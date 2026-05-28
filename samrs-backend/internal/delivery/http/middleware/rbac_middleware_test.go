package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"samrs-backend/internal/test/mocks"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

type apiResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Errors  json.RawMessage `json:"errors"`
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) apiResponse {
	t.Helper()
	var resp apiResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

func TestRBACMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		roleID      interface{}
		setup       func(uc *mocks.MockRBACUsecase)
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "missing role id",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Role ID tidak ditemukan",
		},
		{
			name:        "invalid role id",
			roleID:      "not-int",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Format Role ID tidak valid",
		},
		{
			name:   "role not found",
			roleID: 1,
			setup: func(uc *mocks.MockRBACUsecase) {
				uc.EXPECT().Authorize(1, "room:read").Return(false, errTest("not found"))
			},
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Role tidak ditemukan",
		},
		{
			name:   "access denied",
			roleID: 2,
			setup: func(uc *mocks.MockRBACUsecase) {
				uc.EXPECT().Authorize(2, "room:read").Return(false, nil)
			},
			wantStatus:  http.StatusForbidden,
			wantMessage: "Akses ditolak",
		},
		{
			name:   "access allowed",
			roleID: 3,
			setup: func(uc *mocks.MockRBACUsecase) {
				uc.EXPECT().Authorize(3, "room:read").Return(true, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRBAC := mocks.NewMockRBACUsecase(ctrl)
			if tt.setup != nil {
				tt.setup(mockRBAC)
			}

			router := gin.New()
			if tt.roleID != nil {
				router.Use(func(c *gin.Context) {
					c.Set("role_id", tt.roleID)
					c.Next()
				})
			}
			router.GET("/protected", RBACMiddleware(mockRBAC, "room:read"), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if tt.wantMessage != "" {
				resp := decodeResponse(t, rec)
				if resp.Message != tt.wantMessage {
					t.Fatalf("expected message %q, got %q", tt.wantMessage, resp.Message)
				}
			}
		})
	}
}

type errTest string

func (e errTest) Error() string {
	return string(e)
}
