package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"OpenCodeGoPool/internal/model"
)

var (
	reScriptSrc      = regexp.MustCompile(`<script[^>]+src="(/_build/assets/[^"]+\.js)"`)
	reChunkInBundle  = regexp.MustCompile(`"(_build/assets/[^"]+\.js)"`)
	reCreateSrvRef   = regexp.MustCompile(`createServerReference\("([a-f0-9]{64})"\)`)
)

// ExtractServerFnID fetches the usage page, resolves all lazy-loaded JS chunks
// via __vite__mapDeps in the entry bundle, and returns the 64-char hash of the
// server function that returns daily cost history (identified by "totalCost" keyword).
func (c *Client) ExtractServerFnID(cookie, workspaceID string) (string, error) {
	html, err := c.fetchPage(cookie, "/workspace/"+workspaceID+"/usage")
	if err != nil {
		return "", fmt.Errorf("fetch usage page: %w", err)
	}

	// Collect all chunk paths from __vite__mapDeps inside the entry bundle(s).
	var chunkPaths []string
	seen := map[string]bool{}
	for _, m := range reScriptSrc.FindAllStringSubmatch(html, -1) {
		entryBody, err := c.fetchBundle(m[1])
		if err != nil {
			continue
		}
		for _, cm := range reChunkInBundle.FindAllStringSubmatch(entryBody, -1) {
			p := "/" + cm[1]
			if !seen[p] {
				seen[p] = true
				chunkPaths = append(chunkPaths, p)
			}
		}
	}

	// Find the bundle that handles daily cost data (contains "totalCost") and
	// extract its createServerReference hash.
	for _, path := range chunkPaths {
		base := path[strings.LastIndex(path, "/")+1:]
		if strings.Contains(base, "server-runtime-") || strings.Contains(base, "i18n-") {
			continue
		}

		body, err := c.fetchBundle(path)
		if err != nil {
			continue
		}

		if !strings.Contains(body, "totalCost") {
			continue
		}

		if sub := reCreateSrvRef.FindStringSubmatch(body); sub != nil {
			return sub[1], nil
		}
	}

	return "", fmt.Errorf("usage-cost server-fn hash not found in %d chunk(s)", len(chunkPaths))
}

// FetchUsageHistoryMonth fetches daily usage aggregates for one calendar month.
// month is 0-indexed (0=January, 11=December) as required by the API.
func (c *Client) FetchUsageHistoryMonth(cookie, workspaceID, serverFnID string, year, month int) ([]model.UsageDailyHistory, error) {
	reqBody, err := json.Marshal(map[string]any{
		"t": map[string]any{
			"t": 9,
			"i": 0,
			"l": 4,
			"a": []map[string]any{
				{"t": 1, "s": workspaceID},
				{"t": 0, "s": year},
				{"t": 0, "s": month},
				{"t": 1, "s": "+00:00"},
			},
			"o": 0,
		},
		"f": 31,
		"m": []any{},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/_server", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Server-Id", serverFnID)
	req.Header.Set("X-Server-Instance", "server-fn:0")
	req.Header.Set("Cookie", cookie)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server-fn HTTP %d", resp.StatusCode)
	}

	js, err := readFirstFrame(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response frame: %w", err)
	}

	return parseHistoryFrame(js)
}

// fetchBundle downloads a JS bundle from /_build/assets/.
func (c *Client) fetchBundle(path string) (string, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("bundle HTTP %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// readFirstFrame reads one frame from the streaming /_server response.
// Frame header format: ;0xXXXXXXXX; (12 bytes total).
func readFirstFrame(r io.Reader) (string, error) {
	header := make([]byte, 12)
	if _, err := io.ReadFull(r, header); err != nil {
		return "", fmt.Errorf("read frame header: %w", err)
	}

	if string(header[:3]) != ";0x" || header[11] != ';' {
		return "", fmt.Errorf("invalid frame header: %q", header)
	}

	var frameLen uint32
	if _, err := fmt.Sscanf(string(header[3:11]), "%x", &frameLen); err != nil {
		return "", fmt.Errorf("parse frame length from %q: %w", header[3:11], err)
	}

	frame := make([]byte, frameLen)
	if _, err := io.ReadFull(r, frame); err != nil {
		return "", fmt.Errorf("read frame body (%d bytes): %w", frameLen, err)
	}

	return string(frame), nil
}

// parseHistoryFrame extracts the usage array from a normalised seroval JS frame.
func parseHistoryFrame(js string) ([]model.UsageDailyHistory, error) {
	normalized := normalizeJS(js)

	const usageKey = `"usage":`
	idx := strings.Index(normalized, usageKey)
	if idx == -1 {
		return nil, fmt.Errorf("usage key not found in server-fn response")
	}

	rest := normalized[idx+len(usageKey):]
	arrStart := strings.Index(rest, "[")
	if arrStart == -1 {
		return nil, fmt.Errorf("usage array start not found")
	}

	block, err := findBalanced(normalized, idx+len(usageKey)+arrStart, '[', ']')
	if err != nil {
		return nil, fmt.Errorf("extract usage array: %w", err)
	}

	var raw []struct {
		Date      string `json:"date"`
		Model     string `json:"model"`
		TotalCost int64  `json:"totalCost"`
		KeyID     string `json:"keyId"`
		Plan      string `json:"plan"`
	}
	if err := json.Unmarshal([]byte(block), &raw); err != nil {
		return nil, fmt.Errorf("unmarshal history data: %w", err)
	}

	now := time.Now().UTC()
	out := make([]model.UsageDailyHistory, 0, len(raw))
	for _, r := range raw {
		out = append(out, model.UsageDailyHistory{
			Date:      r.Date,
			Model:     r.Model,
			KeyID:     r.KeyID,
			Plan:      r.Plan,
			TotalCost: r.TotalCost,
			ScrapedAt: now,
		})
	}
	return out, nil
}
