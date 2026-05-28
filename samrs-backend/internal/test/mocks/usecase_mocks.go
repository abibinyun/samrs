package mocks

import (
	"reflect"

	"samrs-backend/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

type MockAuthUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecaseMockRecorder
}

type MockAuthUsecaseMockRecorder struct {
	mock *MockAuthUsecase
}

func NewMockAuthUsecase(ctrl *gomock.Controller) *MockAuthUsecase {
	return &MockAuthUsecase{ctrl: ctrl, recorder: &MockAuthUsecaseMockRecorder{}}
}

func (m *MockAuthUsecase) EXPECT() *MockAuthUsecaseMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockAuthUsecase) Login(username string, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", username, password)
	token, _ := ret[0].(string)
	err, _ := ret[1].(error)
	return token, err
}

func (mr *MockAuthUsecaseMockRecorder) Login(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthUsecase)(nil).Login), username, password)
}

func (m *MockAuthUsecase) GetMe(userID uuid.UUID) (*domain.User, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMe", userID)
	user, _ := ret[0].(*domain.User)
	perms, _ := ret[1].([]string)
	err, _ := ret[2].(error)
	return user, perms, err
}

func (mr *MockAuthUsecaseMockRecorder) GetMe(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMe", reflect.TypeOf((*MockAuthUsecase)(nil).GetMe), userID)
}

type MockRBACUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockRBACUsecaseMockRecorder
}

type MockRBACUsecaseMockRecorder struct {
	mock *MockRBACUsecase
}

func NewMockRBACUsecase(ctrl *gomock.Controller) *MockRBACUsecase {
	return &MockRBACUsecase{ctrl: ctrl, recorder: &MockRBACUsecaseMockRecorder{}}
}

func (m *MockRBACUsecase) EXPECT() *MockRBACUsecaseMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockRBACUsecase) Authorize(roleID int, permissionSlug string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", roleID, permissionSlug)
	allowed, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return allowed, err
}

func (mr *MockRBACUsecaseMockRecorder) Authorize(roleID, permissionSlug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockRBACUsecase)(nil).Authorize), roleID, permissionSlug)
}
