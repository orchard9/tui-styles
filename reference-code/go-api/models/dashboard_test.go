package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardMetricsJSON(t *testing.T) {
	metrics := DashboardMetrics{
		TotalEarnings:    1234.56,
		EarningsChange:   5.2,
		TotalSales:       42,
		SalesChange:      3.7,
		MarketplaceViews: 1500,
		ViewsChange:      8.1,
		ActiveAvatars:    15,
		PendingReview:    2,
		AverageRating:    4.6,
		ConversionRate:   5.3,
	}

	data, err := json.Marshal(metrics)
	require.NoError(t, err, "Should marshal without error")

	var decoded DashboardMetrics
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")
	assert.Equal(t, metrics, decoded, "Round-trip should preserve data")
}

func TestTopAvatarJSON(t *testing.T) {
	avatar := TopAvatar{
		AvatarID: "avatar_123",
		Name:     "Cyber Warrior",
		Sales:    42,
		Revenue:  499.58,
		Rating:   4.7,
	}

	data, err := json.Marshal(avatar)
	require.NoError(t, err, "Should marshal without error")

	var decoded TopAvatar
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")
	assert.Equal(t, avatar, decoded, "Round-trip should preserve data")
}

func TestRecentActivityJSON(t *testing.T) {
	timestamp := time.Now().UTC()

	tests := []struct {
		name     string
		activity RecentActivity
	}{
		{
			name: "sale activity",
			activity: RecentActivity{
				Type:       "sale",
				AvatarName: "Test Avatar",
				Amount:     12.99,
				Timestamp:  timestamp,
			},
		},
		{
			name: "review activity",
			activity: RecentActivity{
				Type:       "review",
				AvatarName: "Test Avatar",
				Rating:     5,
				Timestamp:  timestamp,
			},
		},
		{
			name: "activity with both amount and rating",
			activity: RecentActivity{
				Type:       "sale",
				AvatarName: "Test Avatar",
				Amount:     9.99,
				Rating:     4,
				Timestamp:  timestamp,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.activity)
			require.NoError(t, err, "Should marshal without error")

			var decoded RecentActivity
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err, "Should unmarshal without error")

			assert.Equal(t, tt.activity.Type, decoded.Type)
			assert.Equal(t, tt.activity.AvatarName, decoded.AvatarName)
			assert.Equal(t, tt.activity.Amount, decoded.Amount)
			assert.Equal(t, tt.activity.Rating, decoded.Rating)
		})
	}
}

func TestRecentActivityOmitempty(t *testing.T) {
	// Test that Amount and Rating are omitted when zero
	activity := RecentActivity{
		Type:       "view",
		AvatarName: "Test Avatar",
		Timestamp:  time.Now().UTC(),
		// Amount and Rating are zero values
	}

	data, err := json.Marshal(activity)
	require.NoError(t, err, "Should marshal without error")

	// Zero values with omitempty should be omitted from JSON
	dataStr := string(data)
	assert.NotContains(t, dataStr, "amount", "Amount should be omitted when zero")
	assert.NotContains(t, dataStr, "rating", "Rating should be omitted when zero")
}

func TestDashboardResponseJSON(t *testing.T) {
	response := DashboardResponse{
		CreatorID: "creator_123",
		Period:    "month",
		Metrics: DashboardMetrics{
			TotalEarnings:    2000.0,
			EarningsChange:   5.0,
			TotalSales:       150,
			SalesChange:      3.0,
			MarketplaceViews: 3000,
			ViewsChange:      7.0,
			ActiveAvatars:    20,
			PendingReview:    1,
			AverageRating:    4.5,
			ConversionRate:   5.0,
		},
		TopAvatars: []TopAvatar{
			{AvatarID: "avatar_1", Name: "Avatar 1", Sales: 50, Revenue: 500.0, Rating: 4.8},
			{AvatarID: "avatar_2", Name: "Avatar 2", Sales: 40, Revenue: 400.0, Rating: 4.6},
		},
		RecentActivity: []RecentActivity{
			{Type: "sale", AvatarName: "Avatar 1", Amount: 10.0, Timestamp: time.Now().UTC()},
			{Type: "review", AvatarName: "Avatar 2", Rating: 5, Timestamp: time.Now().UTC()},
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err, "Should marshal without error")

	var decoded DashboardResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")

	assert.Equal(t, response.CreatorID, decoded.CreatorID)
	assert.Equal(t, response.Period, decoded.Period)
	assert.Equal(t, response.Metrics, decoded.Metrics)
	assert.Len(t, decoded.TopAvatars, 2)
	assert.Len(t, decoded.RecentActivity, 2)
}

func TestBalanceJSON(t *testing.T) {
	balance := Balance{
		Available:            500.00,
		Pending:              100.00,
		PendingAvailableDate: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		LifetimeTotal:        5000.00,
	}

	data, err := json.Marshal(balance)
	require.NoError(t, err, "Should marshal without error")

	var decoded Balance
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")
	assert.Equal(t, balance, decoded, "Round-trip should preserve data")
}

func TestEarningsBreakdownJSON(t *testing.T) {
	breakdown := EarningsBreakdown{
		AvatarSales:     4500.00,
		ReferralBonuses: 300.00,
		ContestPrizes:   150.00,
		PlatformBonuses: 50.00,
	}

	data, err := json.Marshal(breakdown)
	require.NoError(t, err, "Should marshal without error")

	var decoded EarningsBreakdown
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")
	assert.Equal(t, breakdown, decoded, "Round-trip should preserve data")
}

func TestPayoutHistoryJSON(t *testing.T) {
	tests := []struct {
		name   string
		payout PayoutHistory
	}{
		{
			name: "completed payout",
			payout: PayoutHistory{
				PayoutID:    "payout_123",
				Amount:      1000.00,
				Status:      "completed",
				Method:      "paypal",
				RequestedAt: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
				CompletedAt: time.Date(2025, 11, 5, 12, 30, 0, 0, time.UTC),
			},
		},
		{
			name: "pending payout without completed date",
			payout: PayoutHistory{
				PayoutID:    "payout_456",
				Amount:      500.00,
				Status:      "pending",
				Method:      "bank_transfer",
				RequestedAt: time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC),
				// CompletedAt is zero value
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.payout)
			require.NoError(t, err, "Should marshal without error")

			var decoded PayoutHistory
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err, "Should unmarshal without error")

			assert.Equal(t, tt.payout.PayoutID, decoded.PayoutID)
			assert.Equal(t, tt.payout.Amount, decoded.Amount)
			assert.Equal(t, tt.payout.Status, decoded.Status)
			assert.Equal(t, tt.payout.Method, decoded.Method)

			// CompletedAt might be omitted in JSON for pending payouts
			if !tt.payout.CompletedAt.IsZero() {
				assert.False(t, decoded.CompletedAt.IsZero(),
					"CompletedAt should be preserved when set")
			}
		})
	}
}

func TestPayoutHistoryOmitempty(t *testing.T) {
	// Test that CompletedAt serializes as zero time when not set
	// Note: Go's json.Encoder doesn't omit zero time.Time values even with omitempty
	// This is expected behavior - zero time serializes as "0001-01-01T00:00:00Z"
	payout := PayoutHistory{
		PayoutID:    "payout_123",
		Amount:      500.00,
		Status:      "pending",
		Method:      "paypal",
		RequestedAt: time.Now().UTC(),
		// CompletedAt is zero value
	}

	data, err := json.Marshal(payout)
	require.NoError(t, err, "Should marshal without error")

	var decoded PayoutHistory
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")

	// CompletedAt should be zero value after round-trip
	assert.True(t, decoded.CompletedAt.IsZero(),
		"CompletedAt should be zero when not set")
}

func TestEarningsResponseJSON(t *testing.T) {
	response := EarningsResponse{
		CreatorID: "creator_123",
		Balance: Balance{
			Available:            500.00,
			Pending:              100.00,
			PendingAvailableDate: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			LifetimeTotal:        5000.00,
		},
		Breakdown: EarningsBreakdown{
			AvatarSales:     4750.00,
			ReferralBonuses: 150.00,
			ContestPrizes:   100.00,
			PlatformBonuses: 0.00,
		},
		PayoutHistory: []PayoutHistory{
			{
				PayoutID:    "payout_001",
				Amount:      1200.00,
				Status:      "completed",
				Method:      "paypal",
				RequestedAt: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
				CompletedAt: time.Date(2025, 11, 5, 12, 30, 0, 0, time.UTC),
			},
		},
		NextPayoutDate: time.Date(2025, 12, 7, 0, 0, 0, 0, time.UTC),
		MinimumPayout:  50.00,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err, "Should marshal without error")

	var decoded EarningsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")

	assert.Equal(t, response.CreatorID, decoded.CreatorID)
	assert.Equal(t, response.Balance, decoded.Balance)
	assert.Equal(t, response.Breakdown, decoded.Breakdown)
	assert.Len(t, decoded.PayoutHistory, 1)
	assert.Equal(t, response.MinimumPayout, decoded.MinimumPayout)
}

func TestEarningsResponseEmptyHistory(t *testing.T) {
	// Test with empty payout history
	response := EarningsResponse{
		CreatorID: "creator_new",
		Balance: Balance{
			Available:            0.00,
			Pending:              0.00,
			PendingAvailableDate: time.Now().UTC(),
			LifetimeTotal:        0.00,
		},
		Breakdown: EarningsBreakdown{
			AvatarSales:     0.00,
			ReferralBonuses: 0.00,
			ContestPrizes:   0.00,
			PlatformBonuses: 0.00,
		},
		PayoutHistory:  []PayoutHistory{},
		NextPayoutDate: time.Now().UTC().AddDate(0, 1, 0),
		MinimumPayout:  50.00,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err, "Should marshal without error")

	var decoded EarningsResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")

	assert.Empty(t, decoded.PayoutHistory, "Payout history should be empty")
}
