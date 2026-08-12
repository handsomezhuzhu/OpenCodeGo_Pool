package scraper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	http    *http.Client
	baseURL string
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		http: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		baseURL: "https://opencode.ai",
	}
}

func (c *Client) fetchPage(cookie, path string) (string, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept", "text/html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return "", fmt.Errorf("cookie expired (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode == 302 || resp.StatusCode == 301 {
		return "", fmt.Errorf("cookie expired (redirect to login)")
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	if !bytes.Contains(body, []byte("$R")) {
		return "", fmt.Errorf("cookie expired or invalid, or page format has changed (expected $R marker not found)")
	}

	return string(body), nil
}

func (c *Client) FetchQuotaPage(cookie, workspaceID string) (string, error) {
	return c.fetchPage(cookie, "/workspace/"+workspaceID+"/go")
}

func (c *Client) FetchUsagePage(cookie, workspaceID string) (string, error) {
	return c.fetchPage(cookie, "/workspace/"+workspaceID+"/usage")
}
