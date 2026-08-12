package model

import "time"

type UsageLimit struct {
	Status       string `json:"status"`
	ResetInSec   int    `json:"resetInSec"`
	UsagePercent int    `json:"usagePercent"`
}

type QuotaData struct {
	Mine         bool       `json:"mine"`
	UseBalance   bool       `json:"useBalance"`
	RollingUsage UsageLimit `json:"rollingUsage"`
	WeeklyUsage  UsageLimit `json:"weeklyUsage"`
	MonthlyUsage UsageLimit `json:"monthlyUsage"`
}

type QuotaSnapshot struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	RollingPercent  int       `json:"rolling_percent"`
	RollingStatus   string    `json:"rolling_status"`
	RollingResetSec int       `json:"rolling_reset_sec"`
	WeeklyPercent   int       `json:"weekly_percent"`
	WeeklyStatus    string    `json:"weekly_status"`
	WeeklyResetSec  int       `json:"weekly_reset_sec"`
	MonthlyPercent  int       `json:"monthly_percent"`
	MonthlyStatus   string    `json:"monthly_status"`
	MonthlyResetSec int       `json:"monthly_reset_sec"`
	ScrapedAt       time.Time `json:"scraped_at"`
}
