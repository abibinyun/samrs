package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	return recorder
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		body        string
		setup       func(uc *mocks.MockAuthUsecase)
		wantStatus  int
		wantMessage string
		wantToken   bool
	}{
		{
			name:        "missing fields",
			body:        `{}`,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Username dan password wajib diisi",
		},
		{
			name: "invalid credentials",
			body: `{"username":"alice","password":"wrong"}`,
			setup: func(uc *mocks.MockAuthUsecase) {
				uc.EXPECT().Login("alice", "wrong").Return("", util.ErrUnauthorized("username atau password salah"))
			},
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Login gagal",
		},
		{
			name: "success",
			body: `{"username":"alice","password":"secret"}`,
			setup: func(uc *mocks.MockAuthUsecase) {
				uc.EXPECT().Login("alice", "secret").Return("token-123", nil)
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Login berhasil",
			wantToken:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockAuth := mocks.NewMockAuthUsecase(ctrl)
			if tt.setup != nil {
				tt.setup(mockAuth)
			}
			handler := NewAuthHandler(mockAuth)
			router := gin.New()
			router.POST("/api/v1/auth/login", handler.Login)

			recorder := performRequest(router, http.MethodPost, "/api/v1/auth/login", []byte(tt.body))
			if recorder.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, recorder.Code)
			}
			resp := decodeResponse(t, recorder)
			if resp.Message != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, resp.Message)
			}
			if tt.wantToken {
				var payload struct {
					Token string `json:"token"`
				}
				if err := json.Unmarshal(resp.Data, &payload); err != nil {
					t.Fatalf("failed to decode token: %v", err)
				}
				if payload.Token == "" {
					t.Fatalf("expected token, got empty string")
				}
			}
		})
	}
}

func TestAuthHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validUserID := uuid.New()
	user := &domain.User{
		ID:       validUserID,
		TenantID: uuid.New(),
		RoleID:   1,
		Username: "tester",
		IsActive: true,
	}

	tests := []struct {
		name        string
		userID      interface{}
		setup       func(uc *mocks.MockAuthUsecase)
		wantStatus  int
		wantMessage string
		wantPerms   int
	}{
		{
			name:        "missing user id",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "User ID tidak ditemukan",
		},
		{
			name:        "invalid user id",
			userID:      "not-uuid",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Format User ID tidak valid",
		},
		{
			name:   "usecase error",
			userID: validUserID.String(),
			setup: func(uc *mocks.MockAuthUsecase) {
				uc.EXPECT().GetMe(validUserID).Return(nil, nil, errTest("db down"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "Gagal mengambil data user",
		},
		{
			name:   "success",
			userID: validUserID.String(),
			setup: func(uc *mocks.MockAuthUsecase) {
				uc.EXPECT().GetMe(validUserID).Return(user, []string{"room:read", "asset:read"}, nil)
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Data user ditemukan",
			wantPerms:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockAuth := mocks.NewMockAuthUsecase(ctrl)
			if tt.setup != nil {
				tt.setup(mockAuth)
			}
			handler := NewAuthHandler(mockAuth)

			router := gin.New()
			if tt.userID != nil {
				router.Use(func(c *gin.Context) {
					c.Set("user_id", tt.userID)
					c.Next()
				})
			}
			router.GET("/api/v1/auth/me", handler.Me)

			recorder := performRequest(router, http.MethodGet, "/api/v1/auth/me", nil)
			if recorder.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, recorder.Code)
			}
			resp := decodeResponse(t, recorder)
			if resp.Message != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, resp.Message)
			}
			if tt.wantPerms > 0 {
				var payload struct {
					User        domain.User `json:"user"`
					Permissions []string    `json:"permissions"`
				}
				if err := json.Unmarshal(resp.Data, &payload); err != nil {
					t.Fatalf("failed to decode payload: %v", err)
				}
				if payload.User.ID != validUserID {
					t.Fatalf("expected user %v, got %v", validUserID, payload.User.ID)
				}
				if len(payload.Permissions) != tt.wantPerms {
					t.Fatalf("expected %d permissions, got %d", tt.wantPerms, len(payload.Permissions))
				}
			}
		})
	}
}

type errTest string

func (e errTest) Error() string {
	return string(e)
}
