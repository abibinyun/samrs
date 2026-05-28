package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestNotificationUsecase_Send(t *testing.T) {
	tests := []struct {
		name       string
		req        NotificationRequest
		setup      func(provider *mockNotificationProvider)
		usecase    func(provider *mockNotificationProvider) NotificationUsecase
		wantErr    string
	}{
		{
			name:    "missing channel",
			req:     NotificationRequest{Channel: " "},
			usecase: func(provider *mockNotificationProvider) NotificationUsecase { return NewNotificationUsecase(nil) },
			wantErr: "channel wajib diisi",
		},
		{
			name: "unsupported channel",
			req:  NotificationRequest{Channel: "sms"},
			setup: func(provider *mockNotificationProvider) {
				provider.EXPECT().Channel().Return("email")
			},
			usecase: func(provider *mockNotificationProvider) NotificationUsecase {
				return NewNotificationUsecase([]NotificationProvider{provider})
			},
			wantErr: "channel tidak didukung",
		},
		{
			name: "provider error",
			req: NotificationRequest{
				Channel: " EMAIL ",
				To:      "ops@example.com",
				Message: "ping",
			},
			setup: func(provider *mockNotificationProvider) {
				provider.EXPECT().Channel().Return("Email")
				provider.EXPECT().Send(gomock.Any(), NotificationRequest{
					Channel: " EMAIL ",
					To:      "ops@example.com",
					Message: "ping",
				}).Return(errors.New("send failed"))
			},
			usecase: func(provider *mockNotificationProvider) NotificationUsecase {
				return NewNotificationUsecase([]NotificationProvider{provider})
			},
			wantErr: "send failed",
		},
		{
			name: "success with normalized channel",
			req: NotificationRequest{
				Channel: "EMAIL",
				To:      "ops@example.com",
				Message: "ok",
			},
			setup: func(provider *mockNotificationProvider) {
				provider.EXPECT().Channel().Return("Email")
				provider.EXPECT().Send(gomock.Any(), NotificationRequest{
					Channel: "EMAIL",
					To:      "ops@example.com",
					Message: "ok",
				}).Return(nil)
			},
			usecase: func(provider *mockNotificationProvider) NotificationUsecase {
				return NewNotificationUsecase([]NotificationProvider{provider})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			provider := newMockNotificationProvider(ctrl)
			if tt.setup != nil {
				tt.setup(provider)
			}
			uc := tt.usecase(provider)
			err := uc.Send(context.Background(), tt.req)
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

type mockNotificationProvider struct {
	ctrl     *gomock.Controller
	recorder *mockNotificationProviderRecorder
}

type mockNotificationProviderRecorder struct {
	mock *mockNotificationProvider
}

func newMockNotificationProvider(ctrl *gomock.Controller) *mockNotificationProvider {
	return &mockNotificationProvider{ctrl: ctrl, recorder: &mockNotificationProviderRecorder{}}
}

func (m *mockNotificationProvider) EXPECT() *mockNotificationProviderRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *mockNotificationProvider) Channel() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Channel")
	channel, _ := ret[0].(string)
	return channel
}

func (mr *mockNotificationProviderRecorder) Channel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Channel", reflect.TypeOf((*mockNotificationProvider)(nil).Channel))
}

func (m *mockNotificationProvider) Send(ctx context.Context, req NotificationRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, req)
	err, _ := ret[0].(error)
	return err
}

func (mr *mockNotificationProviderRecorder) Send(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*mockNotificationProvider)(nil).Send), ctx, req)
}
