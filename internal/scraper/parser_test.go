package scraper

import (
	"os"
	"strings"
	"testing"
)

func TestParseQuota(t *testing.T) {
	data, err := os.ReadFile("../../docs/opencodego-抓取页面获取限额.md")
	if err != nil {
		t.Fatal(err)
	}

	// Extract HTML block from markdown
	content := string(data)
	idx := strings.Index(content, "```html")
	if idx == -1 {
		t.Fatal("no html block found")
	}
	content = content[idx+7:]
	end := strings.LastIndex(content, "```")
	if end == -1 {
		t.Fatal("no closing html block")
	}
	html := content[:end]

	snap, err := ParseQuota(html, "test-account")
	if err != nil {
		t.Fatal("ParseQuota failed:", err)
	}

	if snap.RollingPercent != 22 {
		t.Errorf("rolling percent: want 22, got %d", snap.RollingPercent)
	}
	if snap.WeeklyPercent != 9 {
		t.Errorf("weekly percent: want 9, got %d", snap.WeeklyPercent)
	}
	if snap.MonthlyPercent != 7 {
		t.Errorf("monthly percent: want 7, got %d", snap.MonthlyPercent)
	}
	if snap.RollingStatus != "ok" {
		t.Errorf("rolling status: want ok, got %s", snap.RollingStatus)
	}

	t.Logf("Quota parsed: rolling=%d%% weekly=%d%% monthly=%d%%",
		snap.RollingPercent, snap.WeeklyPercent, snap.MonthlyPercent)
}

func TestParseUsage(t *testing.T) {
	data, err := os.ReadFile("../../docs/opencodego-抓取页面获取使用量情况.md")
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	idx := strings.Index(content, "```html")
	if idx == -1 {
		t.Fatal("no html block found")
	}
	content = content[idx+7:]
	end := strings.LastIndex(content, "```")
	if end == -1 {
		t.Fatal("no closing html block")
	}
	html := content[:end]

	records, err := ParseUsage(html, "test-account")
	if err != nil {
		t.Fatal("ParseUsage failed:", err)
	}

	if len(records) == 0 {
		t.Fatal("no usage records parsed")
	}

	t.Logf("Parsed %d usage records", len(records))

	first := records[0]
	if !strings.HasPrefix(first.ID, "usg_") {
		t.Errorf("first record ID should start with usg_, got %s", first.ID)
	}
	if first.Model == "" {
		t.Error("first record has empty model")
	}

	t.Logf("First record: model=%s provider=%s cost=%d", first.Model, first.Provider, first.Cost)
}

func TestParseQuota_MarkerNotFoundErrorIsEnriched(t *testing.T) {
	html := `<html><body><script>$R[0] = {"unrelated":true}</script></body></html>`

	_, err := ParseQuota(html, "test-account")
	if err == nil {
		t.Fatal("expected error when rollingUsage marker is absent")
	}
	if !strings.Contains(err.Error(), "scanned 1 script(s)") {
		t.Errorf("expected enriched error with script count, got: %v", err)
	}
	if !strings.Contains(err.Error(), "bytes") {
		t.Errorf("expected enriched error with byte count, got: %v", err)
	}
}

func TestParseUsage_MarkerNotFoundErrorIsEnriched(t *testing.T) {
	html := `<html><body><script>$R[0] = {"unrelated":true}</script></body></html>`

	_, err := ParseUsage(html, "test-account")
	if err == nil {
		t.Fatal("expected error when usage.list marker is absent")
	}
	if !strings.Contains(err.Error(), "scanned 1 script(s)") {
		t.Errorf("expected enriched error with script count, got: %v", err)
	}
	if !strings.Contains(err.Error(), "bytes") {
		t.Errorf("expected enriched error with byte count, got: %v", err)
	}
}
