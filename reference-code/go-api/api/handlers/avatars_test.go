package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/masquerade/creator-api/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleAvatarUpload(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name: "successful upload",
			requestBody: `{
				"creator_id": "creator_123",
				"avatar_name": "Test Avatar",
				"photo_count": 25
			}`,
			expectedStatus: http.StatusAccepted,
			checkResponse: func(t *testing.T, body string) {
				var resp models.UploadResponse
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)

				assert.NotEmpty(t, resp.UploadID)
				assert.NotEmpty(t, resp.AvatarID)
				assert.Equal(t, "processing", resp.Status)
				assert.Equal(t, "Avatar upload started successfully", resp.Message)
				assert.NotZero(t, resp.EstimatedCompletion)
				assert.Len(t, resp.ProcessingStages, 5)

				// Verify all stages are queued
				for _, stage := range resp.ProcessingStages {
					assert.Equal(t, "queued", stage.Status)
				}

				// Verify stage names
				stageNames := []string{"face_detection", "3d_reconstruction", "texture_generation", "expression_rigging", "quality_verification"}
				for i, expectedName := range stageNames {
					assert.Equal(t, expectedName, resp.ProcessingStages[i].Stage)
				}
			},
		},
		{
			name:           "invalid json",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "missing creator_id",
			requestBody: `{
				"avatar_name": "Test Avatar",
				"photo_count": 25
			}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "missing avatar_name",
			requestBody: `{
				"creator_id": "creator_123",
				"photo_count": 25
			}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "invalid photo_count zero",
			requestBody: `{
				"creator_id": "creator_123",
				"avatar_name": "Test Avatar",
				"photo_count": 0
			}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "invalid photo_count negative",
			requestBody: `{
				"creator_id": "creator_123",
				"avatar_name": "Test Avatar",
				"photo_count": -5
			}`,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/avatars/upload", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			HandleAvatarUpload(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

func TestHandleGetAvatar(t *testing.T) {
	tests := []struct {
		name           string
		avatarID       string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "existing avatar - cyber warrior",
			avatarID:       "avatar_a1b2c3",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp models.Avatar
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)

				assert.Equal(t, "avatar_a1b2c3", resp.AvatarID)
				assert.Equal(t, "creator_123", resp.CreatorID)
				assert.Equal(t, "Cyber Warrior", resp.Name)
				assert.Equal(t, "published", resp.Status)
				assert.Equal(t, 8.7, resp.QualityScore)
				assert.Equal(t, "sci-fi", resp.Category)
				assert.Equal(t, []string{"cyborg", "futuristic", "masculine"}, resp.Tags)
				assert.Equal(t, 12.99, resp.Price)
				assert.Equal(t, 47, resp.SalesCount)
				assert.Equal(t, 4.6, resp.Rating)
				assert.Equal(t, 12, resp.ReviewCount)
				assert.Len(t, resp.PreviewImages, 2)
				assert.NotZero(t, resp.CreatedAt)
				assert.NotZero(t, resp.PublishedAt)
			},
		},
		{
			name:           "existing avatar - mystic elf",
			avatarID:       "avatar_d4e5f6",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp models.Avatar
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)

				assert.Equal(t, "avatar_d4e5f6", resp.AvatarID)
				assert.Equal(t, "Mystic Elf", resp.Name)
				assert.Equal(t, "published", resp.Status)
				assert.Equal(t, 9.2, resp.QualityScore)
				assert.Equal(t, "fantasy", resp.Category)
				assert.Len(t, resp.PreviewImages, 3)
			},
		},
		{
			name:           "existing avatar - processing status",
			avatarID:       "avatar_g7h8i9",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp models.Avatar
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err)

				assert.Equal(t, "avatar_g7h8i9", resp.AvatarID)
				assert.Equal(t, "Zombie Apocalypse", resp.Name)
				assert.Equal(t, "processing", resp.Status)
				assert.Equal(t, 0.0, resp.QualityScore)
				assert.Equal(t, 0, resp.SalesCount)
				assert.Empty(t, resp.PreviewImages)
			},
		},
		{
			name:           "avatar not found",
			avatarID:       "avatar_nonexistent",
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/avatars/"+tt.avatarID, nil)
			w := httptest.NewRecorder()

			// Set up chi URL params
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.avatarID)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			HandleGetAvatar(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

func TestHandleGetAvatarEmptyID(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/avatars/", nil)
	w := httptest.NewRecorder()

	// Set up chi URL params with empty ID
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	HandleGetAvatar(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
