// Package models defines data structures for the Creator API.
package models

import "time"

// DashboardMetrics represents key performance metrics for a creator.
type DashboardMetrics struct {
	TotalEarnings    float64 `json:"total_earnings"`
	EarningsChange   float64 `json:"earnings_change"`
	TotalSales       int     `json:"total_sales"`
	SalesChange      float64 `json:"sales_change"`
	MarketplaceViews int     `json:"marketplace_views"`
	ViewsChange      float64 `json:"views_change"`
	ActiveAvatars    int     `json:"active_avatars"`
	PendingReview    int     `json:"pending_review"`
	AverageRating    float64 `json:"average_rating"`
	ConversionRate   float64 `json:"conversion_rate"`
}

// TopAvatar represents a high-performing avatar in the creator's catalog.
type TopAvatar struct {
	AvatarID string  `json:"avatar_id"`
	Name     string  `json:"name"`
	Sales    int     `json:"sales"`
	Revenue  float64 `json:"revenue"`
	Rating   float64 `json:"rating"`
}

// RecentActivity represents a recent event in the creator's account.
type RecentActivity struct {
	Type       string    `json:"type"`
	AvatarName string    `json:"avatar_name"`
	Amount     float64   `json:"amount,omitempty"`
	Rating     int       `json:"rating,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// DashboardResponse represents the complete dashboard data for a creator.
type DashboardResponse struct {
	CreatorID      string           `json:"creator_id"`
	Period         string           `json:"period"`
	Metrics        DashboardMetrics `json:"metrics"`
	TopAvatars     []TopAvatar      `json:"top_avatars"`
	RecentActivity []RecentActivity `json:"recent_activity"`
}

// Balance represents the creator's current balance information.
type Balance struct {
	Available            float64   `json:"available"`
	Pending              float64   `json:"pending"`
	PendingAvailableDate time.Time `json:"pending_available_date"`
	LifetimeTotal        float64   `json:"lifetime_total"`
}

// EarningsBreakdown shows the breakdown of earnings by source.
type EarningsBreakdown struct {
	AvatarSales     float64 `json:"avatar_sales"`
	ReferralBonuses float64 `json:"referral_bonuses"`
	ContestPrizes   float64 `json:"contest_prizes"`
	PlatformBonuses float64 `json:"platform_bonuses"`
}

// PayoutHistory represents a single payout transaction.
type PayoutHistory struct {
	PayoutID    string    `json:"payout_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Method      string    `json:"method"`
	RequestedAt time.Time `json:"requested_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// EarningsResponse represents earnings breakdown and payout history.
type EarningsResponse struct {
	CreatorID      string            `json:"creator_id"`
	Balance        Balance           `json:"balance"`
	Breakdown      EarningsBreakdown `json:"breakdown"`
	PayoutHistory  []PayoutHistory   `json:"payout_history"`
	NextPayoutDate time.Time         `json:"next_payout_date"`
	MinimumPayout  float64           `json:"minimum_payout"`
}
