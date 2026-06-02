package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	complaintusecase "samrs-backend/internal/usecase/complaint"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

type testComplaintUsecase struct {
	ctrl     *gomock.Controller
	recorder *testComplaintUsecaseRecorder
}

type testComplaintUsecaseRecorder struct {
	mock *testComplaintUsecase
}

func newTestComplaintUsecase(ctrl *gomock.Controller) *testComplaintUsecase {
	return &testComplaintUsecase{ctrl: ctrl, recorder: &testComplaintUsecaseRecorder{}}
}

func (m *testComplaintUsecase) EXPECT() *testComplaintUsecaseRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *testComplaintUsecase) CreateComplaint(input complaintusecase.CreateComplaintInput) (*domain.Complaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComplaint", input)
	complaint, _ := ret[0].(*domain.Complaint)
	err, _ := ret[1].(error)
	return complaint, err
}

func (m *testComplaintUsecase) GetAllComplaints(tenantID uuid.UUID, filter repository.ComplaintFilter) ([]domain.Complaint, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllComplaints", tenantID, filter)
	complaints, _ := ret[0].([]domain.Complaint)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return complaints, total, err
}

func (m *testComplaintUsecase) GetComplaintByID(tenantID uuid.UUID, id uuid.UUID) (*domain.Complaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComplaintByID", tenantID, id)
	complaint, _ := ret[0].(*domain.Complaint)
	err, _ := ret[1].(error)
	return complaint, err
}

func (m *testComplaintUsecase) UpdateComplaint(input complaintusecase.UpdateComplaintInput) (*domain.Complaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComplaint", input)
	complaint, _ := ret[0].(*domain.Complaint)
	err, _ := ret[1].(error)
	return complaint, err
}

func (m *testComplaintUsecase) DeleteComplaint(tenantID uuid.UUID, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteComplaint", tenantID, id)
	err, _ := ret[0].(error)
	return err
}

func (mr *testComplaintUsecaseRecorder) GetAllComplaints(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllComplaints", reflect.TypeOf((*testComplaintUsecase)(nil).GetAllComplaints), tenantID, filter)
}

func (mr *testComplaintUsecaseRecorder) GetComplaintByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComplaintByID", reflect.TypeOf((*testComplaintUsecase)(nil).GetComplaintByID), tenantID, id)
}

func TestComplaintHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tenantID := uuid.New()

	tests := []struct {
		name        string
		setup       func(uc *testComplaintUsecase)
		wantStatus  int
		wantMessage string
		wantCount   int
	}{
		{
			name:        "missing tenant id",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Tenant ID tidak ditemukan",
		},
		{
			name: "usecase error",
			setup: func(uc *testComplaintUsecase) {
				uc.EXPECT().GetAllComplaints(tenantID, gomock.Any()).
					Return(nil, int64(0), fmt.Errorf("db error"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "Gagal mengambil data",
		},
		{
			name: "empty result",
			setup: func(uc *testComplaintUsecase) {
				uc.EXPECT().GetAllComplaints(tenantID, gomock.Any()).
					Return([]domain.Complaint{}, int64(0), nil)
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Data complaint ditemukan",
			wantCount:   0,
		},
		{
			name: "success with results",
			setup: func(uc *testComplaintUsecase) {
				uc.EXPECT().GetAllComplaints(tenantID, gomock.Any()).
					Return([]domain.Complaint{
						{ID: uuid.New(), Title: "Complaint 1", Status: domain.ComplaintStatusOpen},
						{ID: uuid.New(), Title: "Complaint 2", Status: domain.ComplaintStatusInProgress},
					}, int64(2), nil)
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Data complaint ditemukan",
			wantCount:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUc := newTestComplaintUsecase(ctrl)
			if tt.setup != nil {
				tt.setup(mockUc)
			}

			handler := NewComplaintHandler(mockUc, nil)
			router := gin.New()
			if tt.name != "missing tenant id" {
				router.Use(func(c *gin.Context) {
					c.Set("tenant_id_uuid", tenantID)
					c.Next()
				})
			}
			router.GET("/api/v1/complaints", handler.GetAll)

			recorder := performRequest(router, http.MethodGet, "/api/v1/complaints", nil)
			if recorder.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, recorder.Code)
			}
			resp := decodeResponse(t, recorder)
			if resp.Message != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, resp.Message)
			}
			if tt.wantCount > 0 {
				var payload struct {
					Data []domain.Complaint `json:"data"`
				}
				_ = json.Unmarshal(resp.Data, &payload)
			}
		})
	}
}

func TestComplaintHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tenantID := uuid.New()
	complaintID := uuid.New()

	tests := []struct {
		name        string
		complaintID string
		setup       func(uc *testComplaintUsecase)
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "missing tenant id",
			complaintID: complaintID.String(),
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Tenant ID tidak ditemukan",
		},
		{
			name:        "invalid complaint id",
			complaintID: "not-uuid",
			setup:       func(uc *testComplaintUsecase) {},
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Format complaint ID tidak valid",
		},
		{
			name:        "not found",
			complaintID: complaintID.String(),
			setup: func(uc *testComplaintUsecase) {
				uc.EXPECT().GetComplaintByID(tenantID, complaintID).
					Return(nil, fmt.Errorf("not found"))
			},
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "Gagal mengambil complaint",
		},
		{
			name:        "success",
			complaintID: complaintID.String(),
			setup: func(uc *testComplaintUsecase) {
				uc.EXPECT().GetComplaintByID(tenantID, complaintID).
					Return(&domain.Complaint{ID: complaintID, Title: "Test Complaint"}, nil)
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Complaint ditemukan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUc := newTestComplaintUsecase(ctrl)
			if tt.setup != nil {
				tt.setup(mockUc)
			}

			handler := NewComplaintHandler(mockUc, nil)
			router := gin.New()
			if tt.name != "missing tenant id" {
				router.Use(func(c *gin.Context) {
					c.Set("tenant_id_uuid", tenantID)
					c.Next()
				})
			}
			router.GET("/api/v1/complaints/:id", handler.GetByID)

			path := fmt.Sprintf("/api/v1/complaints/%s", tt.complaintID)
			recorder := performRequest(router, http.MethodGet, path, nil)
			if recorder.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, recorder.Code)
			}
			resp := decodeResponse(t, recorder)
			if resp.Message != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, resp.Message)
			}
		})
	}
}
