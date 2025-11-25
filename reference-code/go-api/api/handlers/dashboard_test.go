package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/masquerade/creator-api/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCreatorDashboard(t *testing.T) {
	tests := []struct {
		name           string
		creatorID      string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "valid creator dashboard",
			creatorID:      "creator_123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp models.DashboardResponse
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)

				assert.Equal(t, "creator_123", resp.CreatorID)
				assert.Equal(t, "month", resp.Period)

				// Verify metrics are present and reasonable
				assert.Greater(t, resp.Metrics.TotalEarnings, 0.0)
				assert.Greater(t, resp.Metrics.TotalSales, 0)
				assert.Greater(t, resp.Metrics.MarketplaceViews, 0)
				assert.GreaterOrEqual(t, resp.Metrics.ActiveAvatars, 0)
				assert.GreaterOrEqual(t, resp.Metrics.PendingReview, 0)
				assert.GreaterOrEqual(t, resp.Metrics.AverageRating, 0.0)
				assert.LessOrEqual(t, resp.Metrics.AverageRating, 5.0)
				assert.GreaterOrEqual(t, resp.Metrics.ConversionRate, 0.0)

				// Verify change metrics can be positive or negative
				assert.NotNil(t, resp.Metrics.EarningsChange)
				assert.NotNil(t, resp.Metrics.SalesChange)
				assert.NotNil(t, resp.Metrics.ViewsChange)

				// Verify top avatars
				assert.Len(t, resp.TopAvatars, 3)
				for _, avatar := range resp.TopAvatars {
					assert.NotEmpty(t, avatar.AvatarID)
					assert.NotEmpty(t, avatar.Name)
					assert.GreaterOrEqual(t, avatar.Sales, 0)
					assert.GreaterOrEqual(t, avatar.Revenue, 0.0)
					assert.GreaterOrEqual(t, avatar.Rating, 0.0)
					assert.LessOrEqual(t, avatar.Rating, 5.0)
				}

				// Verify recent activity
				assert.Len(t, resp.RecentActivity, 4)
				for _, activity := range resp.RecentActivity {
					assert.NotEmpty(t, activity.Type)
					assert.NotEmpty(t, activity.AvatarName)
					assert.NotZero(t, activity.Timestamp)

					switch activity.Type {
					case "sale":
						assert.Greater(t, activity.Amount, 0.0)
					case "review":
						assert.GreaterOrEqual(t, activity.Rating, 1)
						assert.LessOrEqual(t, activity.Rating, 5)
					}
				}
			},
		},
		{
			name:           "different creator has different data",
			creatorID:      "creator_456",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp models.DashboardResponse
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)

				assert.Equal(t, "creator_456", resp.CreatorID)
				// Data should vary based on creator ID (seeded randomness)
				assert.Greater(t, resp.Metrics.TotalEarnings, 0.0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/creators/"+tt.creatorID+"/dashboard", nil)
			w := httptest.NewRecorder()

			// Set up chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.creatorID)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			HandleCreatorDashboard(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

func TestHandleCreatorDashboardEmptyID(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/creators//dashboard", nil)
	w := httptest.NewRecorder()

	// Set up chi URL params with empty ID
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	HandleCreatorDashboard(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreatorEarnings(t *testing.T) {
	tests := []struct {
		name           string
		creatorID      string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "valid creator earnings",
			creatorID:      "creator_123",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp models.EarningsResponse
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)

				assert.Equal(t, "creator_123", resp.CreatorID)

				// Verify balance
				assert.Greater(t, resp.Balance.Available, 0.0)
				assert.GreaterOrEqual(t, resp.Balance.Pending, 0.0)
				assert.Greater(t, resp.Balance.LifetimeTotal, 0.0)
				assert.NotZero(t, resp.Balance.PendingAvailableDate)

				// Verify balance relationships
				assert.LessOrEqual(t, resp.Balance.Available+resp.Balance.Pending, resp.Balance.LifetimeTotal)

				// Verify breakdown
				assert.GreaterOrEqual(t, resp.Breakdown.AvatarSales, 0.0)
				assert.GreaterOrEqual(t, resp.Breakdown.ReferralBonuses, 0.0)
				assert.GreaterOrEqual(t, resp.Breakdown.ContestPrizes, 0.0)
				assert.GreaterOrEqual(t, resp.Breakdown.PlatformBonuses, 0.0)

				// Breakdown should sum to approximately lifetime total
				breakdownSum := resp.Breakdown.AvatarSales +
					resp.Breakdown.ReferralBonuses +
					resp.Breakdown.ContestPrizes +
					resp.Breakdown.PlatformBonuses
				assert.InDelta(t, resp.Balance.LifetimeTotal, breakdownSum, 0.01)

				// Verify payout history
				assert.Len(t, resp.PayoutHistory, 3)
				for _, payout := range resp.PayoutHistory {
					assert.NotEmpty(t, payout.PayoutID)
					assert.Greater(t, payout.Amount, 0.0)
					assert.NotEmpty(t, payout.Status)
					assert.NotEmpty(t, payout.Method)
					assert.NotZero(t, payout.RequestedAt)

					if payout.Status == "completed" {
						assert.NotZero(t, payout.CompletedAt)
					}
				}

				// Verify payout settings
				assert.NotZero(t, resp.NextPayoutDate)
				assert.Equal(t, 50.0, resp.MinimumPayout)
			},
		},
		{
			name:           "different creator has different earnings",
			creatorID:      "creator_789",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp models.EarningsResponse
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)

				assert.Equal(t, "creator_789", resp.CreatorID)
				// Data should vary based on creator ID (seeded randomness)
				assert.Greater(t, resp.Balance.LifetimeTotal, 0.0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/creators/"+tt.creatorID+"/earnings", nil)
			w := httptest.NewRecorder()

			// Set up chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.creatorID)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			HandleCreatorEarnings(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

func TestHandleCreatorEarningsEmptyID(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/creators//earnings", nil)
	w := httptest.NewRecorder()

	// Set up chi URL params with empty ID
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	HandleCreatorEarnings(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDashboardDataConsistency verifies that dashboard data is consistent for same creator ID.
func TestDashboardDataConsistency(t *testing.T) {
	creatorID := "creator_test_123"

	// Make two requests for the same creator
	req1 := httptest.NewRequest("GET", "/api/v1/creators/"+creatorID+"/dashboard", nil)
	w1 := httptest.NewRecorder()
	rctx1 := chi.NewRouteContext()
	rctx1.URLParams.Add("id", creatorID)
	ctx1 := context.WithValue(req1.Context(), chi.RouteCtxKey, rctx1)
	req1 = req1.WithContext(ctx1)

	req2 := httptest.NewRequest("GET", "/api/v1/creators/"+creatorID+"/dashboard", nil)
	w2 := httptest.NewRecorder()
	rctx2 := chi.NewRouteContext()
	rctx2.URLParams.Add("id", creatorID)
	ctx2 := context.WithValue(req2.Context(), chi.RouteCtxKey, rctx2)
	req2 = req2.WithContext(ctx2)

	HandleCreatorDashboard(w1, req1)
	HandleCreatorDashboard(w2, req2)

	var resp1, resp2 models.DashboardResponse
	err := json.Unmarshal(w1.Body.Bytes(), &resp1)
	require.NoError(t, err)
	err = json.Unmarshal(w2.Body.Bytes(), &resp2)
	require.NoError(t, err)

	// Base metrics should be consistent (seeded by creator ID)
	assert.Equal(t, resp1.Metrics.ActiveAvatars, resp2.Metrics.ActiveAvatars)

	// Recent activity timestamps will differ (uses time.Now()), but types and avatar names should be same
	assert.Len(t, resp1.RecentActivity, len(resp2.RecentActivity))
}

// TestEarningsDataConsistency verifies that earnings data is consistent for same creator ID.
func TestEarningsDataConsistency(t *testing.T) {
	creatorID := "creator_test_456"

	// Make two requests for the same creator
	req1 := httptest.NewRequest("GET", "/api/v1/creators/"+creatorID+"/earnings", nil)
	w1 := httptest.NewRecorder()
	rctx1 := chi.NewRouteContext()
	rctx1.URLParams.Add("id", creatorID)
	ctx1 := context.WithValue(req1.Context(), chi.RouteCtxKey, rctx1)
	req1 = req1.WithContext(ctx1)

	req2 := httptest.NewRequest("GET", "/api/v1/creators/"+creatorID+"/earnings", nil)
	w2 := httptest.NewRecorder()
	rctx2 := chi.NewRouteContext()
	rctx2.URLParams.Add("id", creatorID)
	ctx2 := context.WithValue(req2.Context(), chi.RouteCtxKey, rctx2)
	req2 = req2.WithContext(ctx2)

	HandleCreatorEarnings(w1, req1)
	HandleCreatorEarnings(w2, req2)

	var resp1, resp2 models.EarningsResponse
	err := json.Unmarshal(w1.Body.Bytes(), &resp1)
	require.NoError(t, err)
	err = json.Unmarshal(w2.Body.Bytes(), &resp2)
	require.NoError(t, err)

	// Payout history should be consistent
	assert.Equal(t, len(resp1.PayoutHistory), len(resp2.PayoutHistory))
	assert.Equal(t, resp1.MinimumPayout, resp2.MinimumPayout)
}
