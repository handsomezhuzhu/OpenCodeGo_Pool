package scraper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"OpenCodeGoPool/internal/model"

	"github.com/google/uuid"
)

func (c *Client) FetchQuota(cookie, workspaceID string) (*model.QuotaSnapshot, error) {
	html, err := c.FetchQuotaPage(cookie, workspaceID)
	if err != nil {
		return nil, err
	}
	return ParseQuota(html, "")
}

func ParseQuota(html, accountID string) (*model.QuotaSnapshot, error) {
	scripts := extractInlineScripts(html)
	if len(scripts) == 0 {
		return nil, fmt.Errorf("no inline scripts found")
	}

	for _, script := range scripts {
		idx := strings.Index(script, "rollingUsage")
		if idx == -1 {
			continue
		}

		// Walk backwards to find the opening { of the parent object
		objStart := idx
		for objStart > 0 && script[objStart] != '{' {
			objStart--
		}
		if objStart <= 0 {
			continue
		}

		block, err := findBalanced(script, objStart, '{', '}')
		if err != nil {
			continue
		}

		normalized := normalizeJS(block)

		var data model.QuotaData
		if err := json.Unmarshal([]byte(normalized), &data); err != nil {
			return nil, fmt.Errorf("unmarshal quota data: %w (normalized: %s)", err, truncate(normalized, 500))
		}

		return &model.QuotaSnapshot{
			ID:              uuid.New().String(),
			AccountID:       accountID,
			RollingPercent:  data.RollingUsage.UsagePercent,
			RollingStatus:   data.RollingUsage.Status,
			RollingResetSec: data.RollingUsage.ResetInSec,
			WeeklyPercent:   data.WeeklyUsage.UsagePercent,
			WeeklyStatus:    data.WeeklyUsage.Status,
			WeeklyResetSec:  data.WeeklyUsage.ResetInSec,
			MonthlyPercent:  data.MonthlyUsage.UsagePercent,
			MonthlyStatus:   data.MonthlyUsage.Status,
			MonthlyResetSec: data.MonthlyUsage.ResetInSec,
			ScrapedAt:       time.Now().UTC(),
		}, nil
	}

	return nil, fmt.Errorf("quota data not found in page (scanned %d script(s), %d bytes)", len(scripts), len(html))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
