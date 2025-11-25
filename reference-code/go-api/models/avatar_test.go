package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected UploadRequest
		wantErr  bool
	}{
		{
			name:     "valid upload request",
			jsonData: `{"creator_id":"creator_123","avatar_name":"Test Avatar","photo_count":5}`,
			expected: UploadRequest{
				CreatorID:  "creator_123",
				AvatarName: "Test Avatar",
				PhotoCount: 5,
			},
			wantErr: false,
		},
		{
			name:     "minimal request",
			jsonData: `{"creator_id":"","avatar_name":"","photo_count":0}`,
			expected: UploadRequest{
				CreatorID:  "",
				AvatarName: "",
				PhotoCount: 0,
			},
			wantErr: false,
		},
		{
			name:     "invalid JSON",
			jsonData: `{invalid}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req UploadRequest
			err := json.Unmarshal([]byte(tt.jsonData), &req)

			if tt.wantErr {
				assert.Error(t, err, "Should return error for invalid JSON")
			} else {
				require.NoError(t, err, "Should unmarshal valid JSON")
				assert.Equal(t, tt.expected, req, "Unmarshaled request should match expected")
			}
		})
	}
}

func TestUploadRequestMarshal(t *testing.T) {
	req := UploadRequest{
		CreatorID:  "creator_123",
		AvatarName: "Test Avatar",
		PhotoCount: 5,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err, "Should marshal without error")

	// Unmarshal back to verify round-trip
	var decoded UploadRequest
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal back without error")
	assert.Equal(t, req, decoded, "Round-trip should preserve data")
}

func TestProcessingStageJSON(t *testing.T) {
	stage := ProcessingStage{
		Stage:  "face_detection",
		Status: "processing",
	}

	data, err := json.Marshal(stage)
	require.NoError(t, err, "Should marshal without error")

	var decoded ProcessingStage
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")
	assert.Equal(t, stage, decoded, "Round-trip should preserve data")
}

func TestUploadResponseJSON(t *testing.T) {
	now := time.Now().UTC()
	resp := UploadResponse{
		UploadID:            "upload_123",
		AvatarID:            "avatar_abc",
		Status:              "processing",
		Message:             "Upload started",
		EstimatedCompletion: now,
		ProcessingStages: []ProcessingStage{
			{Stage: "face_detection", Status: "queued"},
			{Stage: "3d_reconstruction", Status: "queued"},
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err, "Should marshal without error")

	var decoded UploadResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")
	assert.Equal(t, resp.UploadID, decoded.UploadID)
	assert.Equal(t, resp.AvatarID, decoded.AvatarID)
	assert.Equal(t, resp.Status, decoded.Status)
	assert.Len(t, decoded.ProcessingStages, 2)
}

func TestAvatarJSON(t *testing.T) {
	createdAt := time.Date(2025, 11, 15, 8, 30, 0, 0, time.UTC)
	publishedAt := time.Date(2025, 11, 17, 14, 20, 0, 0, time.UTC)

	avatar := Avatar{
		AvatarID:      "avatar_123",
		CreatorID:     "creator_456",
		Name:          "Test Avatar",
		Status:        "published",
		QualityScore:  8.5,
		CreatedAt:     createdAt,
		PublishedAt:   publishedAt,
		Category:      "sci-fi",
		Tags:          []string{"robot", "futuristic"},
		Price:         12.99,
		SalesCount:    42,
		Rating:        4.7,
		ReviewCount:   15,
		PreviewImages: []string{"https://example.com/img1.jpg", "https://example.com/img2.jpg"},
	}

	data, err := json.Marshal(avatar)
	require.NoError(t, err, "Should marshal without error")

	var decoded Avatar
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")

	assert.Equal(t, avatar.AvatarID, decoded.AvatarID)
	assert.Equal(t, avatar.CreatorID, decoded.CreatorID)
	assert.Equal(t, avatar.Name, decoded.Name)
	assert.Equal(t, avatar.Status, decoded.Status)
	assert.Equal(t, avatar.QualityScore, decoded.QualityScore)
	assert.Equal(t, avatar.Category, decoded.Category)
	assert.Equal(t, avatar.Tags, decoded.Tags)
	assert.Equal(t, avatar.Price, decoded.Price)
	assert.Equal(t, avatar.SalesCount, decoded.SalesCount)
	assert.Equal(t, avatar.Rating, decoded.Rating)
	assert.Equal(t, avatar.ReviewCount, decoded.ReviewCount)
	assert.Equal(t, avatar.PreviewImages, decoded.PreviewImages)
}

func TestAvatarWithoutPublishedAt(t *testing.T) {
	// Test that PublishedAt serializes as zero time when not published
	// Note: Go's json.Encoder doesn't omit zero time.Time values even with omitempty
	// This is expected behavior - zero time serializes as "0001-01-01T00:00:00Z"
	avatar := Avatar{
		AvatarID:   "avatar_123",
		CreatorID:  "creator_456",
		Name:       "Processing Avatar",
		Status:     "processing",
		CreatedAt:  time.Now().UTC(),
		Category:   "fantasy",
		Tags:       []string{},
		Price:      9.99,
		SalesCount: 0,
	}

	data, err := json.Marshal(avatar)
	require.NoError(t, err, "Should marshal without error")

	var decoded Avatar
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")
	assert.True(t, decoded.PublishedAt.IsZero(), "PublishedAt should be zero value when not set")
}

func TestAvatarEmptySlices(t *testing.T) {
	// Test that empty slices serialize correctly (as [] not null)
	avatar := Avatar{
		AvatarID:      "avatar_123",
		CreatorID:     "creator_456",
		Name:          "Test",
		Status:        "processing",
		CreatedAt:     time.Now().UTC(),
		Category:      "fantasy",
		Tags:          []string{},
		PreviewImages: []string{},
	}

	data, err := json.Marshal(avatar)
	require.NoError(t, err, "Should marshal without error")

	var decoded Avatar
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err, "Should unmarshal without error")

	// Empty slices should unmarshal as empty slices, not nil
	assert.NotNil(t, decoded.Tags, "Tags should not be nil")
	assert.Empty(t, decoded.Tags, "Tags should be empty")
	assert.NotNil(t, decoded.PreviewImages, "PreviewImages should not be nil")
	assert.Empty(t, decoded.PreviewImages, "PreviewImages should be empty")
}

func TestAvatarNilSlices(t *testing.T) {
	// Test that nil slices are handled correctly
	avatar := Avatar{
		AvatarID:      "avatar_123",
		CreatorID:     "creator_456",
		Name:          "Test",
		Status:        "processing",
		CreatedAt:     time.Now().UTC(),
		Category:      "fantasy",
		Tags:          nil,
		PreviewImages: nil,
	}

	data, err := json.Marshal(avatar)
	require.NoError(t, err, "Should marshal without error")

	// In Go's json package, nil slices marshal to null
	assert.Contains(t, string(data), `"tags":null`)
	assert.Contains(t, string(data), `"preview_images":null`)
}
