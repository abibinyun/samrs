package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestAuditTrailUsecase_Log(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		meta      AuditMeta
		oldData   interface{}
		newData   interface{}
		setup     func(repo *mocks.MockAuditTrailRepository)
		wantErr   string
	}{
		{
			name:    "missing tenant",
			meta:    AuditMeta{UserID: userID},
			wantErr: "metadata audit tidak lengkap",
		},
		{
			name:    "missing user",
			meta:    AuditMeta{TenantID: tenantID},
			wantErr: "metadata audit tidak lengkap",
		},
		{
			name:    "old data marshal error",
			meta:    AuditMeta{TenantID: tenantID, UserID: userID},
			oldData: func() {},
			wantErr: "json: unsupported type: func()",
		},
		{
			name:    "new data marshal error",
			meta:    AuditMeta{TenantID: tenantID, UserID: userID},
			newData: func() {},
			wantErr: "json: unsupported type: func()",
		},
		{
			name: "repo error",
			meta: AuditMeta{TenantID: tenantID, UserID: userID},
			setup: func(repo *mocks.MockAuditTrailRepository) {
				repo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			meta: AuditMeta{TenantID: tenantID, UserID: userID, IP: "127.0.0.1", UserAgent: "tester"},
			oldData: map[string]interface{}{
				"status": "old",
			},
			newData: map[string]interface{}{
				"status": "new",
			},
			setup: func(repo *mocks.MockAuditTrailRepository) {
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(audit *domain.AuditTrail) error {
					if audit.TenantID != tenantID {
						t.Fatalf("expected tenant id %s, got %s", tenantID, audit.TenantID)
					}
					if audit.UserID != userID {
						t.Fatalf("expected user id %s, got %s", userID, audit.UserID)
					}
					if audit.Action != "UPDATE" {
						t.Fatalf("expected action UPDATE, got %s", audit.Action)
					}
					if audit.TableName != "assets" {
						t.Fatalf("expected table assets, got %s", audit.TableName)
					}
					if audit.RecordID != "abc" {
						t.Fatalf("expected record id abc, got %s", audit.RecordID)
					}
					if audit.OldData == nil || audit.NewData == nil {
						t.Fatalf("expected old/new data to be set")
					}
					if audit.IP != "127.0.0.1" {
						t.Fatalf("expected ip 127.0.0.1, got %s", audit.IP)
					}
					if audit.UserAgent != "tester" {
						t.Fatalf("expected user agent tester, got %s", audit.UserAgent)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockAuditTrailRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			uc := NewAuditTrailUsecase(repo)
			err := uc.Log("UPDATE", "assets", "abc", tt.oldData, tt.newData, tt.meta)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAuditTrailUsecase_ListAndGet(t *testing.T) {
	tenantID := uuid.New()
	filter := repository.AuditTrailFilter{Page: 1, PerPage: 10}

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockAuditTrailRepository)
		wantErr string
	}{
		{
			name: "list error",
			setup: func(repo *mocks.MockAuditTrailRepository) {
				repo.EXPECT().ListByTenant(tenantID, filter).Return(nil, int64(0), errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "list success",
			setup: func(repo *mocks.MockAuditTrailRepository) {
				repo.EXPECT().ListByTenant(tenantID, filter).Return([]domain.AuditTrail{{ID: 1}}, int64(1), nil)
			},
		},
		{
			name: "get error",
			setup: func(repo *mocks.MockAuditTrailRepository) {
				repo.EXPECT().FindByID(tenantID, uint(1)).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "get success",
			setup: func(repo *mocks.MockAuditTrailRepository) {
				repo.EXPECT().FindByID(tenantID, uint(1)).Return(&domain.AuditTrail{ID: 1}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockAuditTrailRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			uc := NewAuditTrailUsecase(repo)

			switch tt.name {
			case "list error", "list success":
				_, _, err := uc.ListByTenant(tenantID, filter)
				if tt.wantErr != "" {
					if err == nil {
						t.Fatalf("expected error %q, got nil", tt.wantErr)
					}
					if err.Error() != tt.wantErr {
						t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			case "get error", "get success":
				_, err := uc.GetByID(tenantID, 1)
				if tt.wantErr != "" {
					if err == nil {
						t.Fatalf("expected error %q, got nil", tt.wantErr)
					}
					if err.Error() != tt.wantErr {
						t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
