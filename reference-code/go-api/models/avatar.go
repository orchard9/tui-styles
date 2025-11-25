// Package models defines data structures for the Creator API.
package models

import "time"

// UploadRequest represents the request body for avatar upload.
type UploadRequest struct {
	CreatorID  string `json:"creator_id"`
	AvatarName string `json:"avatar_name"`
	PhotoCount int    `json:"photo_count"`
}

// ProcessingStage represents a single stage in the avatar processing pipeline.
type ProcessingStage struct {
	Stage  string `json:"stage"`
	Status string `json:"status"`
}

// UploadResponse represents the response for avatar upload initiation.
type UploadResponse struct {
	UploadID            string            `json:"upload_id"`
	AvatarID            string            `json:"avatar_id"`
	Status              string            `json:"status"`
	Message             string            `json:"message"`
	EstimatedCompletion time.Time         `json:"estimated_completion"`
	ProcessingStages    []ProcessingStage `json:"processing_stages"`
}

// Avatar represents complete avatar metadata and status.
type Avatar struct {
	AvatarID      string    `json:"avatar_id"`
	CreatorID     string    `json:"creator_id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	QualityScore  float64   `json:"quality_score"`
	CreatedAt     time.Time `json:"created_at"`
	PublishedAt   time.Time `json:"published_at,omitempty"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags"`
	Price         float64   `json:"price"`
	SalesCount    int       `json:"sales_count"`
	Rating        float64   `json:"rating"`
	ReviewCount   int       `json:"review_count"`
	PreviewImages []string  `json:"preview_images"`
}
