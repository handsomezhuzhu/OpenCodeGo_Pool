package model

import "time"

type UsageRecord struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	Model           string    `json:"model"`
	Provider        string    `json:"provider"`
	InputTokens     int       `json:"inputTokens"`
	OutputTokens    int       `json:"outputTokens"`
	ReasoningTokens *int      `json:"reasoningTokens"`
	CacheReadTokens int       `json:"cacheReadTokens"`
	Cost            int64     `json:"cost"`
	TimeCreated     time.Time `json:"timeCreated"`
	Plan            string    `json:"plan"`
	ScrapedAt       time.Time `json:"scraped_at"`
}

type UsageDailyHistory struct {
	AccountID string    `json:"account_id"`
	Date      string    `json:"date"`
	Model     string    `json:"model"`
	KeyID     string    `json:"key_id"`
	Plan      string    `json:"plan"`
	TotalCost int64     `json:"total_cost"`
	ScrapedAt time.Time `json:"scraped_at"`
}

type UsageDailySummary struct {
	Date            string `json:"date"`
	Model           string `json:"model"`
	Provider        string `json:"provider"`
	InputTokens     int    `json:"inputTokens"`
	OutputTokens    int    `json:"outputTokens"`
	ReasoningTokens int    `json:"reasoningTokens"`
	CacheReadTokens int    `json:"cacheReadTokens"`
	Cost            int64  `json:"cost"`
}
