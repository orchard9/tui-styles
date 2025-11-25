// Package mock provides mock data generators for development and testing.
// These responses match the product specification for creator onboarding.
package mock

import (
	"fmt"
	"hash/fnv"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/masquerade/creator-api/pkg/models"
)

// SignupResponse represents the response structure for creator signup.
type SignupResponse struct {
	CreatorID   string    `json:"creator_id"`
	Message     string    `json:"message"`
	TrialEndsAt time.Time `json:"trial_ends_at"`
	Status      string    `json:"status"`
	NextSteps   []string  `json:"next_steps"`
}

// GenerateSignupResponse creates a realistic mock signup response for a new creator.
// This follows the product spec: 14-day trial, status "trial", standard onboarding steps.
func GenerateSignupResponse(creatorName, email string) SignupResponse {
	return SignupResponse{
		CreatorID:   uuid.New().String(),
		Message:     fmt.Sprintf("Welcome to Masquerade Creator Studio, %s!", creatorName),
		TrialEndsAt: time.Now().UTC().AddDate(0, 0, 14), // 14 days from now
		Status:      "trial",
		NextSteps: []string{
			"Complete identity verification",
			"Create your first avatar",
			"Publish to marketplace",
		},
	}
}

// GenerateUploadResponse creates a mock upload response with processing stages.
// The response includes an upload ID, avatar ID, processing status, and estimated completion time.
func GenerateUploadResponse(req models.UploadRequest) models.UploadResponse {
	return models.UploadResponse{
		UploadID:            fmt.Sprintf("upload_%s", uuid.New().String()[:8]),
		AvatarID:            fmt.Sprintf("avatar_%s", uuid.New().String()[:8]),
		Status:              "processing",
		Message:             fmt.Sprintf("Processing %s with %d photos", req.AvatarName, req.PhotoCount),
		EstimatedCompletion: time.Now().UTC().Add(10 * time.Minute),
		ProcessingStages: []models.ProcessingStage{
			{Stage: "face_detection", Status: "queued"},
			{Stage: "3d_reconstruction", Status: "queued"},
			{Stage: "texture_generation", Status: "queued"},
			{Stage: "expression_rigging", Status: "queued"},
			{Stage: "quality_verification", Status: "queued"},
		},
	}
}

// GenerateAvatarData generates mock avatar data based on the avatar ID.
// Returns predefined mock avatars for known IDs or an error for unknown IDs.
// Data is deterministic based on avatar ID to support consistent testing.
func GenerateAvatarData(avatarID string) (models.Avatar, error) {
	// Map of known mock avatars
	mockAvatars := map[string]models.Avatar{
		"avatar_a1b2c3": {
			AvatarID:      "avatar_a1b2c3",
			CreatorID:     "creator_123",
			Name:          "Cyber Warrior",
			Status:        "published",
			QualityScore:  8.7,
			CreatedAt:     time.Now().UTC().AddDate(0, -2, 0),
			PublishedAt:   time.Now().UTC().AddDate(0, -1, -15),
			Category:      "sci-fi",
			Tags:          []string{"cyborg", "futuristic", "masculine"},
			Price:         12.99,
			SalesCount:    47,
			Rating:        4.6,
			ReviewCount:   12,
			PreviewImages: []string{"preview1.jpg", "preview2.jpg"},
		},
		"avatar_d4e5f6": {
			AvatarID:      "avatar_d4e5f6",
			CreatorID:     "creator_456",
			Name:          "Mystic Elf",
			Status:        "published",
			QualityScore:  9.2,
			CreatedAt:     time.Now().UTC().AddDate(0, -3, 0),
			PublishedAt:   time.Now().UTC().AddDate(0, -2, -10),
			Category:      "fantasy",
			Tags:          []string{"elf", "magic", "mystical"},
			Price:         14.99,
			SalesCount:    98,
			Rating:        4.8,
			ReviewCount:   32,
			PreviewImages: []string{"preview1.jpg", "preview2.jpg", "preview3.jpg"},
		},
		"avatar_g7h8i9": {
			AvatarID:      "avatar_g7h8i9",
			CreatorID:     "creator_789",
			Name:          "Zombie Apocalypse",
			Status:        "processing",
			QualityScore:  0.0,
			CreatedAt:     time.Now().UTC().Add(-2 * time.Hour),
			Category:      "sci-fi",
			Tags:          []string{"space", "astronaut"},
			Price:         11.99,
			SalesCount:    0,
			Rating:        0.0,
			ReviewCount:   0,
			PreviewImages: []string{},
		},
	}

	avatar, exists := mockAvatars[avatarID]
	if !exists {
		return models.Avatar{}, fmt.Errorf("avatar %s not found", avatarID)
	}

	return avatar, nil
}

// GenerateDashboardData generates mock dashboard data for a creator.
// Uses deterministic randomization based on creator ID to produce consistent data.
func GenerateDashboardData(creatorID string) models.DashboardResponse {
	seed := hashString(creatorID)
	now := time.Now().UTC()

	// Generate deterministic metrics based on seed
	totalEarnings := 1500.0 + float64(seed%3000)
	totalSales := 150 + (seed % 200)
	marketplaceViews := 3500 + (seed % 3000)
	earningsChange := 2.0 + float64(seed%9)
	activeAvatars := 15 + (seed % 16)
	averageRating := 4.3 + (float64(seed%6) / 10.0)
	conversionRate := 4.0 + (float64(seed%4) / 1.0)

	// Generate top 3 avatars
	topAvatars := []models.TopAvatar{
		{
			AvatarID: fmt.Sprintf("avatar_%d", seed%1000),
			Name:     "Cyber Warrior",
			Sales:    45 + (seed % 20),
			Revenue:  650.0 + float64(seed%300),
			Rating:   4.6 + (float64(seed%4) / 10.0),
		},
		{
			AvatarID: fmt.Sprintf("avatar_%d", (seed+1)%1000),
			Name:     "Mystic Elf",
			Sales:    38 + (seed % 15),
			Revenue:  520.0 + float64(seed%200),
			Rating:   4.7 + (float64(seed%3) / 10.0),
		},
		{
			AvatarID: fmt.Sprintf("avatar_%d", (seed+2)%1000),
			Name:     "Space Explorer",
			Sales:    32 + (seed % 10),
			Revenue:  410.0 + float64(seed%150),
			Rating:   4.5 + (float64(seed%5) / 10.0),
		},
	}

	// Generate 4 recent activities in reverse chronological order
	recentActivity := []models.RecentActivity{
		{
			Type:       "sale",
			AvatarName: "Cyber Warrior",
			Amount:     12.99,
			Timestamp:  now.Add(-2 * time.Hour),
		},
		{
			Type:       "review",
			AvatarName: "Mystic Elf",
			Rating:     5,
			Timestamp:  now.Add(-6 * time.Hour),
		},
		{
			Type:       "sale",
			AvatarName: "Space Explorer",
			Amount:     11.99,
			Timestamp:  now.Add(-12 * time.Hour),
		},
		{
			Type:       "review",
			AvatarName: "Cyber Warrior",
			Rating:     4,
			Timestamp:  now.Add(-24 * time.Hour),
		},
	}

	return models.DashboardResponse{
		CreatorID: creatorID,
		Period:    "month",
		Metrics: models.DashboardMetrics{
			TotalEarnings:    totalEarnings,
			EarningsChange:   earningsChange,
			TotalSales:       totalSales,
			SalesChange:      5.5,
			MarketplaceViews: marketplaceViews,
			ViewsChange:      12.3,
			ActiveAvatars:    activeAvatars,
			PendingReview:    2,
			AverageRating:    averageRating,
			ConversionRate:   conversionRate,
		},
		TopAvatars:     topAvatars,
		RecentActivity: recentActivity,
	}
}

// GenerateEarningsData generates mock earnings data for a creator.
// Uses deterministic randomization based on creator ID to produce consistent data.
func GenerateEarningsData(creatorID string) models.EarningsResponse {
	seed := hashString(creatorID)
	now := time.Now().UTC()

	// Generate deterministic balances based on seed
	availableBalance := 300.0 + float64(seed%500)
	pendingBalance := 50.0 + float64(seed%150)
	lifetimeTotal := 6000.0 + float64(seed%5000)

	// Calculate breakdown (95% avatar sales, 3% referrals, 2% contests, 0% bonuses for simplicity)
	avatarSales := math.Round(lifetimeTotal*0.95*100) / 100
	referralBonuses := math.Round(lifetimeTotal*0.03*100) / 100
	contestPrizes := math.Round(lifetimeTotal*0.02*100) / 100
	platformBonuses := lifetimeTotal - avatarSales - referralBonuses - contestPrizes

	// Generate 3 payout history entries in reverse chronological order
	payoutHistory := []models.PayoutHistory{
		{
			PayoutID:    fmt.Sprintf("payout_%d", seed%10000),
			Amount:      450.0 + float64(seed%200),
			Status:      "completed",
			Method:      "stripe",
			RequestedAt: now.AddDate(0, 0, -7),
			CompletedAt: now.AddDate(0, 0, -5),
		},
		{
			PayoutID:    fmt.Sprintf("payout_%d", (seed+1)%10000),
			Amount:      380.0 + float64(seed%150),
			Status:      "completed",
			Method:      "paypal",
			RequestedAt: now.AddDate(0, 0, -21),
			CompletedAt: now.AddDate(0, 0, -19),
		},
		{
			PayoutID:    fmt.Sprintf("payout_%d", (seed+2)%10000),
			Amount:      520.0 + float64(seed%250),
			Status:      "completed",
			Method:      "stripe",
			RequestedAt: now.AddDate(0, 0, -35),
			CompletedAt: now.AddDate(0, 0, -33),
		},
	}

	return models.EarningsResponse{
		CreatorID: creatorID,
		Balance: models.Balance{
			Available:            availableBalance,
			Pending:              pendingBalance,
			PendingAvailableDate: now.AddDate(0, 0, 14),
			LifetimeTotal:        lifetimeTotal,
		},
		Breakdown: models.EarningsBreakdown{
			AvatarSales:     avatarSales,
			ReferralBonuses: referralBonuses,
			ContestPrizes:   contestPrizes,
			PlatformBonuses: platformBonuses,
		},
		PayoutHistory:  payoutHistory,
		NextPayoutDate: now.AddDate(0, 0, 7),
		MinimumPayout:  50.0,
	}
}

// hashString creates a deterministic hash seed from a string.
// Used to generate consistent mock data for the same input.
func hashString(s string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s)) // hash.Hash.Write never returns an error
	return int(h.Sum32())
}
