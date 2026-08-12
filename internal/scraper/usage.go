package scraper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"OpenCodeGoPool/internal/model"
)

type rawUsageRecord struct {
	ID              string  `json:"id"`
	Model           string  `json:"model"`
	Provider        string  `json:"provider"`
	InputTokens     int     `json:"inputTokens"`
	OutputTokens    int     `json:"outputTokens"`
	ReasoningTokens *int    `json:"reasoningTokens"`
	CacheReadTokens int     `json:"cacheReadTokens"`
	Cost            int64   `json:"cost"`
	TimeCreated     string  `json:"timeCreated"`
	Enrichment      *struct {
		Plan string `json:"plan"`
	} `json:"enrichment"`
}

func (c *Client) FetchUsage(cookie, workspaceID string) ([]model.UsageRecord, error) {
	html, err := c.FetchUsagePage(cookie, workspaceID)
	if err != nil {
		return nil, err
	}
	return ParseUsage(html, "")
}

func ParseUsage(html, accountID string) ([]model.UsageRecord, error) {
	scripts := extractInlineScripts(html)
	if len(scripts) == 0 {
		return nil, fmt.Errorf("no inline scripts found")
	}

	var lastErr error
	for _, script := range scripts {
		for _, block := range extractUsageBlocks(script) {
			normalized := normalizeJS(block)

			var rawRecords []rawUsageRecord
			if err := json.Unmarshal([]byte(normalized), &rawRecords); err != nil {
				lastErr = fmt.Errorf("unmarshal usage data: %w (normalized: %s)", err, truncate(normalized, 500))
				continue
			}

			if len(rawRecords) == 0 {
				return []model.UsageRecord{}, nil
			}

			now := time.Now().UTC()
			records := make([]model.UsageRecord, 0, len(rawRecords))
			for _, r := range rawRecords {
				t, _ := time.Parse(time.RFC3339Nano, r.TimeCreated)
				if t.IsZero() {
					t, _ = time.Parse("2006-01-02T15:04:05.000Z", r.TimeCreated)
				}

				plan := ""
				if r.Enrichment != nil {
					plan = r.Enrichment.Plan
				}

				records = append(records, model.UsageRecord{
					ID:              r.ID,
					AccountID:       accountID,
					Model:           r.Model,
					Provider:        r.Provider,
					InputTokens:     r.InputTokens,
					OutputTokens:    r.OutputTokens,
					ReasoningTokens: r.ReasoningTokens,
					CacheReadTokens: r.CacheReadTokens,
					Cost:            r.Cost,
					TimeCreated:     t,
					Plan:            plan,
					ScrapedAt:       now,
				})
			}
			return records, nil
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf("usage data not found in page (scanned %d script(s), %d bytes)", len(scripts), len(html))
}

func extractUsageBlocks(script string) []string {
	patterns := []string{
		`usg_`,
		`inputTokens:`,
		`outputTokens:`,
		`timeCreated:`,
		`provider:`,
		`model:`,
		`cost:`,
	}

	seen := make(map[int]struct{})
	blocks := make([]string, 0, 4)
	for _, pattern := range patterns {
		searchFrom := 0
		for searchFrom < len(script) {
			idx := strings.Index(script[searchFrom:], pattern)
			if idx == -1 {
				break
			}
			idx += searchFrom

			arrStart, ok := findEnclosingArrayStart(script, idx)
			if ok {
				if _, exists := seen[arrStart]; !exists {
					block, err := findBalanced(script, arrStart, '[', ']')
					if err == nil {
						blocks = append(blocks, block)
						seen[arrStart] = struct{}{}
					}
				}
			}

			searchFrom = idx + len(pattern)
		}
	}

	if len(blocks) > 0 {
		return blocks
	}

	if block, ok := extractUsageBlockAfterMarker(script, `usage.list[`, false); ok {
		return []string{block}
	}

	return nil
}

func extractUsageBlockAfterMarker(script, marker string, requireUsageRecord bool) (string, bool) {
	searchFrom := 0
	for searchFrom < len(script) {
		idx := strings.Index(script[searchFrom:], marker)
		if idx == -1 {
			return "", false
		}
		idx += searchFrom

		arrIdx := strings.Index(script[idx:], `= [`)
		if arrIdx == -1 {
			searchFrom = idx + len(marker)
			continue
		}
		arrStart := idx + arrIdx + 2

		block, err := findBalanced(script, arrStart, '[', ']')
		if err != nil {
			searchFrom = arrStart + 1
			continue
		}
		if requireUsageRecord && !isUsageRecordBlock(block) {
			searchFrom = arrStart + 1
			continue
		}
		return block, true
	}

	return "", false
}

func isUsageRecordBlock(block string) bool {
	trimmed := strings.TrimSpace(block)
	return trimmed == "[]" || strings.Contains(trimmed, `id: "usg_`)
}

func findEnclosingArrayStart(script string, pos int) (int, bool) {
	if pos < 0 || pos >= len(script) {
		return 0, false
	}

	stack := make([]int, 0, 8)
	inString := false
	escaped := false
	for i := 0; i <= pos; i++ {
		ch := script[i]

		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '[':
			stack = append(stack, i)
		case ']':
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}

	if len(stack) == 0 {
		return 0, false
	}

	return stack[len(stack)-1], true
}
