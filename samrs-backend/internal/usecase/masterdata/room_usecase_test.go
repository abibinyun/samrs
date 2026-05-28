package usecase

import (
	"errors"
	"testing"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/test/mocks"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestRoomUsecase_CreateRoom(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name    string
		room    *domain.Room
		setup   func(repo *mocks.MockRoomRepository)
		wantErr string
	}{
		{
			name: "missing name",
			room: &domain.Room{TenantID: tenantID, Name: " ", Code: "R-1"},
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().ExistsByCode(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: "nama dan kode ruangan wajib diisi",
		},
		{
			name: "missing code",
			room: &domain.Room{TenantID: tenantID, Name: "Room A", Code: " "},
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().ExistsByCode(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: "nama dan kode ruangan wajib diisi",
		},
		{
			name: "code exists",
			room: &domain.Room{TenantID: tenantID, Name: "Room A", Code: "R-1"},
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().ExistsByCode(tenantID, "R-1", nil).Return(true, nil)
			},
			wantErr: "kode ruangan sudah digunakan",
		},
		{
			name: "exists check error",
			room: &domain.Room{TenantID: tenantID, Name: "Room A", Code: "R-1"},
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().ExistsByCode(tenantID, "R-1", nil).Return(false, errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			room: &domain.Room{TenantID: tenantID, Name: "Room A", Code: "R-1"},
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().ExistsByCode(tenantID, "R-1", nil).Return(false, nil)
				repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(room *domain.Room) error {
					if room.Name != "Room A" || room.Code != "R-1" {
						t.Fatalf("unexpected room data: %+v", room)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRoomRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewRoomUsecase(repo)
			err := usecase.CreateRoom(tt.room)
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

func TestRoomUsecase_UpdateRoom(t *testing.T) {
	tenantID := uuid.New()
	roomID := uuid.New()

	tests := []struct {
		name     string
		input    struct{ name, code, location string }
		setup    func(repo *mocks.MockRoomRepository)
		wantErr  string
	}{
		{
			name: "missing name",
			input: struct{ name, code, location string }{
				name: " ",
				code: "R-1",
			},
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().FindByIDAndTenant(gomock.Any(), gomock.Any()).Times(0)
			},
			wantErr: "nama dan kode ruangan wajib diisi",
		},
		{
			name: "room not found",
			input: struct{ name, code, location string }{
				name: "Room A",
				code: "R-1",
			},
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "code exists",
			input: struct{ name, code, location string }{
				name: "Room A",
				code: "R-1",
			},
			setup: func(repo *mocks.MockRoomRepository) {
				room := &domain.Room{ID: roomID, TenantID: tenantID, Name: "Old", Code: "R-Old"}
				repo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(room, nil)
				repo.EXPECT().ExistsByCode(tenantID, "R-1", &roomID).Return(true, nil)
			},
			wantErr: "kode ruangan sudah digunakan",
		},
		{
			name: "update error",
			input: struct{ name, code, location string }{
				name: "Room A",
				code: "R-1",
			},
			setup: func(repo *mocks.MockRoomRepository) {
				room := &domain.Room{ID: roomID, TenantID: tenantID, Name: "Old", Code: "R-Old"}
				repo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(room, nil)
				repo.EXPECT().ExistsByCode(tenantID, "R-1", &roomID).Return(false, nil)
				repo.EXPECT().Update(room).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			input: struct{ name, code, location string }{
				name:     "Room A",
				code:     "R-1",
				location: "L1",
			},
			setup: func(repo *mocks.MockRoomRepository) {
				room := &domain.Room{ID: roomID, TenantID: tenantID, Name: "Old", Code: "R-Old"}
				repo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(room, nil)
				repo.EXPECT().ExistsByCode(tenantID, "R-1", &roomID).Return(false, nil)
				repo.EXPECT().Update(room).DoAndReturn(func(updated *domain.Room) error {
					if updated.Name != "Room A" || updated.Code != "R-1" || updated.Location != "L1" {
						t.Fatalf("unexpected update: %+v", updated)
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRoomRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewRoomUsecase(repo)
			_, err := usecase.UpdateRoom(tenantID, roomID, tt.input.name, tt.input.code, tt.input.location)
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

func TestRoomUsecase_DeleteRoom(t *testing.T) {
	tenantID := uuid.New()
	roomID := uuid.New()

	tests := []struct {
		name    string
		setup   func(repo *mocks.MockRoomRepository)
		wantErr string
	}{
		{
			name: "room not found",
			setup: func(repo *mocks.MockRoomRepository) {
				repo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(nil, errors.New("not found"))
			},
			wantErr: "not found",
		},
		{
			name: "delete error",
			setup: func(repo *mocks.MockRoomRepository) {
				room := &domain.Room{ID: roomID, TenantID: tenantID}
				repo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(room, nil)
				repo.EXPECT().Delete(room).Return(errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "success",
			setup: func(repo *mocks.MockRoomRepository) {
				room := &domain.Room{ID: roomID, TenantID: tenantID}
				repo.EXPECT().FindByIDAndTenant(roomID, tenantID).Return(room, nil)
				repo.EXPECT().Delete(room).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mocks.NewMockRoomRepository(ctrl)
			if tt.setup != nil {
				tt.setup(repo)
			}
			usecase := NewRoomUsecase(repo)
			err := usecase.DeleteRoom(tenantID, roomID)
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
