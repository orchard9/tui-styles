package mock

import (
	"testing"
	"time"

	"github.com/masquerade/creator-api/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSignupResponse(t *testing.T) {
	tests := []struct {
		name        string
		creatorName string
		email       string
	}{
		{
			name:        "standard signup",
			creatorName: "Test Creator",
			email:       "test@example.com",
		},
		{
			name:        "creator with special characters",
			creatorName: "José García",
			email:       "jose@example.com",
		},
		{
			name:        "long creator name",
			creatorName: "A Very Long Creator Name That Exceeds Normal Length",
			email:       "long@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := GenerateSignupResponse(tt.creatorName, tt.email)

			// Validate structure
			assert.NotEmpty(t, resp.CreatorID, "CreatorID should not be empty")
			assert.Contains(t, resp.Message, tt.creatorName, "Message should contain creator name")
			assert.Equal(t, "trial", resp.Status, "Status should be 'trial'")

			// Validate trial period (14 days from now)
			now := time.Now().UTC()
			expectedTrialEnd := now.AddDate(0, 0, 14)
			assert.WithinDuration(t, expectedTrialEnd, resp.TrialEndsAt, 2*time.Second,
				"Trial should end 14 days from now")

			// Validate next steps
			assert.Len(t, resp.NextSteps, 3, "Should have 3 next steps")
			assert.Contains(t, resp.NextSteps[0], "verification", "First step should mention verification")
			assert.Contains(t, resp.NextSteps[1], "avatar", "Second step should mention avatar")
			assert.Contains(t, resp.NextSteps[2], "marketplace", "Third step should mention marketplace")
		})
	}
}

func TestGenerateUploadResponse(t *testing.T) {
	tests := []struct {
		name string
		req  models.UploadRequest
	}{
		{
			name: "standard upload",
			req: models.UploadRequest{
				CreatorID:  "creator_123",
				AvatarName: "Test Avatar",
				PhotoCount: 5,
			},
		},
		{
			name: "upload with multiple photos",
			req: models.UploadRequest{
				CreatorID:  "creator_456",
				AvatarName: "Multi-Photo Avatar",
				PhotoCount: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := GenerateUploadResponse(tt.req)

			// Validate IDs
			assert.NotEmpty(t, resp.UploadID, "UploadID should not be empty")
			assert.Contains(t, resp.UploadID, "upload_", "UploadID should have upload_ prefix")
			assert.NotEmpty(t, resp.AvatarID, "AvatarID should not be empty")
			assert.Contains(t, resp.AvatarID, "avatar_", "AvatarID should have avatar_ prefix")

			// Validate status and message
			assert.Equal(t, "processing", resp.Status, "Status should be 'processing'")
			assert.NotEmpty(t, resp.Message, "Message should not be empty")

			// Validate estimated completion (10 minutes from now)
			now := time.Now().UTC()
			expectedCompletion := now.Add(10 * time.Minute)
			assert.WithinDuration(t, expectedCompletion, resp.EstimatedCompletion, 2*time.Second,
				"Estimated completion should be 10 minutes from now")

			// Validate processing stages
			require.Len(t, resp.ProcessingStages, 5, "Should have 5 processing stages")
			expectedStages := []string{"face_detection", "3d_reconstruction", "texture_generation", "expression_rigging", "quality_verification"}
			for i, stage := range resp.ProcessingStages {
				assert.Equal(t, expectedStages[i], stage.Stage, "Stage name should match")
				assert.Equal(t, "queued", stage.Status, "All stages should start as 'queued'")
			}
		})
	}
}

func TestGenerateAvatarData(t *testing.T) {
	tests := []struct {
		name      string
		avatarID  string
		expectErr bool
		validate  func(t *testing.T, avatar models.Avatar)
	}{
		{
			name:      "cyber warrior avatar",
			avatarID:  "avatar_a1b2c3",
			expectErr: false,
			validate: func(t *testing.T, avatar models.Avatar) {
				assert.Equal(t, "avatar_a1b2c3", avatar.AvatarID)
				assert.Equal(t, "creator_123", avatar.CreatorID)
				assert.Equal(t, "Cyber Warrior", avatar.Name)
				assert.Equal(t, "published", avatar.Status)
				assert.Equal(t, 8.7, avatar.QualityScore)
				assert.Equal(t, "sci-fi", avatar.Category)
				assert.Contains(t, avatar.Tags, "cyborg")
				assert.Equal(t, 12.99, avatar.Price)
				assert.Greater(t, avatar.SalesCount, 0)
				assert.Greater(t, avatar.Rating, 0.0)
				assert.Len(t, avatar.PreviewImages, 2)
			},
		},
		{
			name:      "mystic elf avatar",
			avatarID:  "avatar_d4e5f6",
			expectErr: false,
			validate: func(t *testing.T, avatar models.Avatar) {
				assert.Equal(t, "avatar_d4e5f6", avatar.AvatarID)
				assert.Equal(t, "Mystic Elf", avatar.Name)
				assert.Equal(t, "fantasy", avatar.Category)
				assert.Contains(t, avatar.Tags, "elf")
				assert.Equal(t, 14.99, avatar.Price)
				assert.Len(t, avatar.PreviewImages, 3)
			},
		},
		{
			name:      "processing avatar",
			avatarID:  "avatar_g7h8i9",
			expectErr: false,
			validate: func(t *testing.T, avatar models.Avatar) {
				assert.Equal(t, "avatar_g7h8i9", avatar.AvatarID)
				assert.Equal(t, "processing", avatar.Status)
				assert.Equal(t, 0.0, avatar.QualityScore, "Processing avatar should have 0 quality score")
				assert.Equal(t, 0, avatar.SalesCount, "Processing avatar should have 0 sales")
				assert.Empty(t, avatar.PreviewImages, "Processing avatar should have no preview images")
			},
		},
		{
			name:      "unknown avatar returns error",
			avatarID:  "unknown_avatar",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avatar, err := GenerateAvatarData(tt.avatarID)

			if tt.expectErr {
				assert.Error(t, err, "Should return error for unknown avatar")
				assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")
			} else {
				require.NoError(t, err, "Should not return error for known avatar")
				if tt.validate != nil {
					tt.validate(t, avatar)
				}
			}
		})
	}
}

func TestGenerateDashboardData(t *testing.T) {
	tests := []struct {
		name      string
		creatorID string
	}{
		{
			name:      "creator 123",
			creatorID: "creator_123",
		},
		{
			name:      "creator 456",
			creatorID: "creator_456",
		},
		{
			name:      "creator 789",
			creatorID: "creator_789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := GenerateDashboardData(tt.creatorID)

			// Validate basic structure
			assert.Equal(t, tt.creatorID, resp.CreatorID)
			assert.Equal(t, "month", resp.Period)

			// Validate metrics are within reasonable ranges
			assert.InDelta(t, 3000.0, resp.Metrics.TotalEarnings, 2000.0, "Earnings should be within range")
			assert.InDelta(t, 250, resp.Metrics.TotalSales, 100, "Sales should be within range")
			assert.InDelta(t, 5000, resp.Metrics.MarketplaceViews, 2000, "Views should be within range")
			assert.GreaterOrEqual(t, resp.Metrics.EarningsChange, 2.0, "Earnings change should be >= 2%")
			assert.LessOrEqual(t, resp.Metrics.EarningsChange, 10.0, "Earnings change should be <= 10%")
			assert.GreaterOrEqual(t, resp.Metrics.ActiveAvatars, 15, "Active avatars should be >= 15")
			assert.LessOrEqual(t, resp.Metrics.ActiveAvatars, 30, "Active avatars should be <= 30")
			assert.GreaterOrEqual(t, resp.Metrics.AverageRating, 4.3, "Rating should be >= 4.3")
			assert.LessOrEqual(t, resp.Metrics.AverageRating, 4.8, "Rating should be <= 4.8")
			assert.GreaterOrEqual(t, resp.Metrics.ConversionRate, 4.0, "Conversion rate should be >= 4%")
			assert.LessOrEqual(t, resp.Metrics.ConversionRate, 7.0, "Conversion rate should be <= 7%")

			// Validate top avatars
			assert.Len(t, resp.TopAvatars, 3, "Should have 3 top avatars")
			for i, avatar := range resp.TopAvatars {
				assert.NotEmpty(t, avatar.AvatarID, "Avatar ID should not be empty")
				assert.NotEmpty(t, avatar.Name, "Avatar name should not be empty")
				assert.Greater(t, avatar.Sales, 0, "Sales should be positive")
				assert.Greater(t, avatar.Revenue, 0.0, "Revenue should be positive")
				assert.GreaterOrEqual(t, avatar.Rating, 4.4, "Rating should be >= 4.4")
				assert.LessOrEqual(t, avatar.Rating, 4.9, "Rating should be <= 4.9")

				// First avatar should have highest revenue
				if i > 0 {
					assert.LessOrEqual(t, avatar.Revenue, resp.TopAvatars[i-1].Revenue,
						"Avatars should be sorted by revenue descending")
				}
			}

			// Validate recent activity
			assert.Len(t, resp.RecentActivity, 4, "Should have 4 recent activities")
			for i, activity := range resp.RecentActivity {
				assert.NotEmpty(t, activity.Type, "Activity type should not be empty")
				assert.NotEmpty(t, activity.AvatarName, "Avatar name should not be empty")
				assert.True(t, activity.Timestamp.Before(time.Now().UTC()), "Activity should be in the past")

				// Activities should be in reverse chronological order
				if i > 0 {
					assert.True(t, activity.Timestamp.Before(resp.RecentActivity[i-1].Timestamp),
						"Activities should be sorted by timestamp descending")
				}

				// Validate type-specific fields
				switch activity.Type {
				case "sale":
					assert.Greater(t, activity.Amount, 0.0, "Sale amount should be positive")
				case "review":
					assert.GreaterOrEqual(t, activity.Rating, 1, "Review rating should be >= 1")
					assert.LessOrEqual(t, activity.Rating, 5, "Review rating should be <= 5")
				}
			}
		})
	}
}

func TestGenerateDashboardDataConsistency(t *testing.T) {
	// Test that same creator ID produces same data (deterministic based on seed)
	creatorID := "creator_test_consistency"

	resp1 := GenerateDashboardData(creatorID)
	resp2 := GenerateDashboardData(creatorID)

	assert.Equal(t, resp1.Metrics.TotalEarnings, resp2.Metrics.TotalEarnings,
		"Same creator ID should produce same earnings")
	assert.Equal(t, resp1.Metrics.TotalSales, resp2.Metrics.TotalSales,
		"Same creator ID should produce same sales")
	assert.Equal(t, resp1.Metrics.ActiveAvatars, resp2.Metrics.ActiveAvatars,
		"Same creator ID should produce same active avatars count")
}

func TestGenerateEarningsData(t *testing.T) {
	tests := []struct {
		name      string
		creatorID string
	}{
		{
			name:      "creator 123",
			creatorID: "creator_123",
		},
		{
			name:      "creator 456",
			creatorID: "creator_456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := GenerateEarningsData(tt.creatorID)

			// Validate basic structure
			assert.Equal(t, tt.creatorID, resp.CreatorID)

			// Validate balance ranges
			assert.InDelta(t, 550.0, resp.Balance.Available, 300.0, "Available balance should be within range")
			assert.InDelta(t, 100.0, resp.Balance.Pending, 100.0, "Pending balance should be within range")
			assert.InDelta(t, 8500.0, resp.Balance.LifetimeTotal, 3000.0, "Lifetime total should be within range")

			// Pending available date should be 14 days from now
			expectedPendingDate := time.Now().UTC().AddDate(0, 0, 14)
			assert.WithinDuration(t, expectedPendingDate, resp.Balance.PendingAvailableDate, 24*time.Hour,
				"Pending available date should be ~14 days from now")

			// Validate breakdown adds up to lifetime total
			totalBreakdown := resp.Breakdown.AvatarSales + resp.Breakdown.ReferralBonuses +
				resp.Breakdown.ContestPrizes + resp.Breakdown.PlatformBonuses
			assert.InDelta(t, resp.Balance.LifetimeTotal, totalBreakdown, 0.01,
				"Breakdown should sum to lifetime total")

			// Validate breakdown percentages
			assert.InDelta(t, resp.Balance.LifetimeTotal*0.95, resp.Breakdown.AvatarSales, 0.01,
				"Avatar sales should be 95% of lifetime total")
			assert.InDelta(t, resp.Balance.LifetimeTotal*0.03, resp.Breakdown.ReferralBonuses, 0.01,
				"Referral bonuses should be 3% of lifetime total")
			assert.InDelta(t, resp.Balance.LifetimeTotal*0.02, resp.Breakdown.ContestPrizes, 0.01,
				"Contest prizes should be 2% of lifetime total")

			// Validate payout history
			assert.Len(t, resp.PayoutHistory, 3, "Should have 3 payout history entries")
			for i, payout := range resp.PayoutHistory {
				assert.NotEmpty(t, payout.PayoutID, "Payout ID should not be empty")
				assert.Greater(t, payout.Amount, 0.0, "Payout amount should be positive")
				assert.Equal(t, "completed", payout.Status, "All payouts should be completed")
				assert.NotEmpty(t, payout.Method, "Payment method should not be empty")
				assert.True(t, payout.CompletedAt.After(payout.RequestedAt),
					"Completed date should be after requested date")

				// Payouts should be in reverse chronological order
				if i > 0 {
					assert.True(t, payout.RequestedAt.Before(resp.PayoutHistory[i-1].RequestedAt),
						"Payouts should be sorted by requested date descending")
				}
			}

			// Validate minimum payout and next payout date
			assert.Equal(t, 50.0, resp.MinimumPayout, "Minimum payout should be $50")
			assert.False(t, resp.NextPayoutDate.IsZero(), "Next payout date should be set")
		})
	}
}

func TestGenerateEarningsDataConsistency(t *testing.T) {
	// Test that same creator ID produces same data (deterministic based on seed)
	creatorID := "creator_test_consistency"

	resp1 := GenerateEarningsData(creatorID)
	resp2 := GenerateEarningsData(creatorID)

	assert.Equal(t, resp1.Balance.Available, resp2.Balance.Available,
		"Same creator ID should produce same available balance")
	assert.Equal(t, resp1.Balance.Pending, resp2.Balance.Pending,
		"Same creator ID should produce same pending balance")
	assert.Equal(t, resp1.Balance.LifetimeTotal, resp2.Balance.LifetimeTotal,
		"Same creator ID should produce same lifetime total")
}
