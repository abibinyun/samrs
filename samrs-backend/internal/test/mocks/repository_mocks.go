package mocks

import (
	"reflect"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	return &MockUserRepository{ctrl: ctrl, recorder: &MockUserRepositoryMockRecorder{}}
}

func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUsername", username)
	user, _ := ret[0].(*domain.User)
	err, _ := ret[1].(error)
	return user, err
}

func (mr *MockUserRepositoryMockRecorder) GetByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUsername", reflect.TypeOf((*MockUserRepository)(nil).GetByUsername), username)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	user, _ := ret[0].(*domain.User)
	err, _ := ret[1].(error)
	return user, err
}

func (mr *MockUserRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserRepository)(nil).FindByID), id)
}

func (m *MockUserRepository) FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDAndTenant", id, tenantID)
	user, _ := ret[0].(*domain.User)
	err, _ := ret[1].(error)
	return user, err
}

func (mr *MockUserRepositoryMockRecorder) FindByIDAndTenant(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDAndTenant", reflect.TypeOf((*MockUserRepository)(nil).FindByIDAndTenant), id, tenantID)
}

func (m *MockUserRepository) Create(user *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockUserRepositoryMockRecorder) Create(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), user)
}

func (m *MockUserRepository) FindAllByTenant(tenantID uuid.UUID) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID)
	users, _ := ret[0].([]domain.User)
	err, _ := ret[1].(error)
	return users, err
}

func (mr *MockUserRepositoryMockRecorder) FindAllByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockUserRepository)(nil).FindAllByTenant), tenantID)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", user)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockUserRepositoryMockRecorder) Update(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserRepository)(nil).Update), user)
}

func (m *MockUserRepository) Delete(user *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", user)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockUserRepositoryMockRecorder) Delete(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserRepository)(nil).Delete), user)
}

func (m *MockUserRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockUserRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockUserRepository)(nil).CountByTenant), tenantID)
}

func (m *MockUserRepository) ListByTenant(tenantID uuid.UUID, filter repository.UserFilter) ([]domain.User, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByTenant", tenantID, filter)
	users, _ := ret[0].([]domain.User)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return users, total, err
}

func (mr *MockUserRepositoryMockRecorder) ListByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByTenant", reflect.TypeOf((*MockUserRepository)(nil).ListByTenant), tenantID, filter)
}

type MockPermissionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPermissionRepositoryMockRecorder
}

type MockPermissionRepositoryMockRecorder struct {
	mock *MockPermissionRepository
}

func NewMockPermissionRepository(ctrl *gomock.Controller) *MockPermissionRepository {
	return &MockPermissionRepository{ctrl: ctrl, recorder: &MockPermissionRepositoryMockRecorder{}}
}

func (m *MockPermissionRepository) EXPECT() *MockPermissionRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockPermissionRepository) EnsurePermissions(perms []domain.Permission) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnsurePermissions", perms)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockPermissionRepositoryMockRecorder) EnsurePermissions(perms interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsurePermissions", reflect.TypeOf((*MockPermissionRepository)(nil).EnsurePermissions), perms)
}

func (m *MockPermissionRepository) ListAll() ([]domain.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAll")
	perms, _ := ret[0].([]domain.Permission)
	err, _ := ret[1].(error)
	return perms, err
}

func (mr *MockPermissionRepositoryMockRecorder) ListAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAll", reflect.TypeOf((*MockPermissionRepository)(nil).ListAll))
}

func (m *MockPermissionRepository) ListByRoleID(roleID int) ([]domain.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByRoleID", roleID)
	perms, _ := ret[0].([]domain.Permission)
	err, _ := ret[1].(error)
	return perms, err
}

func (mr *MockPermissionRepositoryMockRecorder) ListByRoleID(roleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByRoleID", reflect.TypeOf((*MockPermissionRepository)(nil).ListByRoleID), roleID)
}

func (m *MockPermissionRepository) ListSlugsByRoleID(roleID int) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSlugsByRoleID", roleID)
	slugs, _ := ret[0].([]string)
	err, _ := ret[1].(error)
	return slugs, err
}

func (mr *MockPermissionRepositoryMockRecorder) ListSlugsByRoleID(roleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSlugsByRoleID", reflect.TypeOf((*MockPermissionRepository)(nil).ListSlugsByRoleID), roleID)
}

func (m *MockPermissionRepository) RoleHasPermission(roleID int, slug string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoleHasPermission", roleID, slug)
	ok, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return ok, err
}

func (mr *MockPermissionRepositoryMockRecorder) RoleHasPermission(roleID, slug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoleHasPermission", reflect.TypeOf((*MockPermissionRepository)(nil).RoleHasPermission), roleID, slug)
}

func (m *MockPermissionRepository) AssignPermissions(roleID int, permissionIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignPermissions", roleID, permissionIDs)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockPermissionRepositoryMockRecorder) AssignPermissions(roleID, permissionIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignPermissions", reflect.TypeOf((*MockPermissionRepository)(nil).AssignPermissions), roleID, permissionIDs)
}

func (m *MockPermissionRepository) RevokePermissions(roleID int, permissionIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokePermissions", roleID, permissionIDs)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockPermissionRepositoryMockRecorder) RevokePermissions(roleID, permissionIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokePermissions", reflect.TypeOf((*MockPermissionRepository)(nil).RevokePermissions), roleID, permissionIDs)
}

type MockRoleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRoleRepositoryMockRecorder
}

type MockRoleRepositoryMockRecorder struct {
	mock *MockRoleRepository
}

func NewMockRoleRepository(ctrl *gomock.Controller) *MockRoleRepository {
	return &MockRoleRepository{ctrl: ctrl, recorder: &MockRoleRepositoryMockRecorder{}}
}

func (m *MockRoleRepository) EXPECT() *MockRoleRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockRoleRepository) FindByID(id int) (*domain.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	role, _ := ret[0].(*domain.Role)
	err, _ := ret[1].(error)
	return role, err
}

func (mr *MockRoleRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockRoleRepository)(nil).FindByID), id)
}

func (m *MockRoleRepository) FindByIDAndTenant(id int, tenantID uuid.UUID) (*domain.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDAndTenant", id, tenantID)
	role, _ := ret[0].(*domain.Role)
	err, _ := ret[1].(error)
	return role, err
}

func (mr *MockRoleRepositoryMockRecorder) FindByIDAndTenant(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDAndTenant", reflect.TypeOf((*MockRoleRepository)(nil).FindByIDAndTenant), id, tenantID)
}

func (m *MockRoleRepository) FindAllByTenant(tenantID uuid.UUID) ([]domain.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID)
	roles, _ := ret[0].([]domain.Role)
	err, _ := ret[1].(error)
	return roles, err
}

func (mr *MockRoleRepositoryMockRecorder) FindAllByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockRoleRepository)(nil).FindAllByTenant), tenantID)
}

func (m *MockRoleRepository) FindAllByTenantIDs(tenantIDs []uuid.UUID) ([]domain.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenantIDs", tenantIDs)
	roles, _ := ret[0].([]domain.Role)
	err, _ := ret[1].(error)
	return roles, err
}

func (mr *MockRoleRepositoryMockRecorder) FindAllByTenantIDs(tenantIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenantIDs", reflect.TypeOf((*MockRoleRepository)(nil).FindAllByTenantIDs), tenantIDs)
}

func (m *MockRoleRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockRoleRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockRoleRepository)(nil).CountByTenant), tenantID)
}

func (m *MockRoleRepository) ListByTenant(tenantID uuid.UUID, filter repository.RoleFilter) ([]domain.Role, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByTenant", tenantID, filter)
	roles, _ := ret[0].([]domain.Role)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return roles, total, err
}

func (mr *MockRoleRepositoryMockRecorder) ListByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByTenant", reflect.TypeOf((*MockRoleRepository)(nil).ListByTenant), tenantID, filter)
}

func (m *MockRoleRepository) Create(role *domain.Role) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", role)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockRoleRepositoryMockRecorder) Create(role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRoleRepository)(nil).Create), role)
}

func (m *MockRoleRepository) Update(role *domain.Role) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", role)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockRoleRepositoryMockRecorder) Update(role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRoleRepository)(nil).Update), role)
}

func (m *MockRoleRepository) Delete(role *domain.Role) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", role)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockRoleRepositoryMockRecorder) Delete(role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoleRepository)(nil).Delete), role)
}

type MockRoomRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRoomRepositoryMockRecorder
}

type MockRoomRepositoryMockRecorder struct {
	mock *MockRoomRepository
}

func NewMockRoomRepository(ctrl *gomock.Controller) *MockRoomRepository {
	return &MockRoomRepository{ctrl: ctrl, recorder: &MockRoomRepositoryMockRecorder{}}
}

func (m *MockRoomRepository) EXPECT() *MockRoomRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockRoomRepository) Create(room *domain.Room) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", room)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockRoomRepositoryMockRecorder) Create(room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRoomRepository)(nil).Create), room)
}

func (m *MockRoomRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.RoomFilter) ([]domain.Room, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	rooms, _ := ret[0].([]domain.Room)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return rooms, total, err
}

func (mr *MockRoomRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockRoomRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockRoomRepository) FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDAndTenant", id, tenantID)
	room, _ := ret[0].(*domain.Room)
	err, _ := ret[1].(error)
	return room, err
}

func (mr *MockRoomRepositoryMockRecorder) FindByIDAndTenant(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDAndTenant", reflect.TypeOf((*MockRoomRepository)(nil).FindByIDAndTenant), id, tenantID)
}

func (m *MockRoomRepository) Update(room *domain.Room) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", room)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockRoomRepositoryMockRecorder) Update(room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRoomRepository)(nil).Update), room)
}

func (m *MockRoomRepository) Delete(room *domain.Room) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", room)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockRoomRepositoryMockRecorder) Delete(room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoomRepository)(nil).Delete), room)
}

func (m *MockRoomRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByCode", tenantID, code, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockRoomRepositoryMockRecorder) ExistsByCode(tenantID, code, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByCode", reflect.TypeOf((*MockRoomRepository)(nil).ExistsByCode), tenantID, code, excludeID)
}

func (m *MockRoomRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockRoomRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockRoomRepository)(nil).CountByTenant), tenantID)
}

type MockBedRepository struct {
	ctrl     *gomock.Controller
	recorder *MockBedRepositoryMockRecorder
}

type MockBedRepositoryMockRecorder struct {
	mock *MockBedRepository
}

func NewMockBedRepository(ctrl *gomock.Controller) *MockBedRepository {
	return &MockBedRepository{ctrl: ctrl, recorder: &MockBedRepositoryMockRecorder{}}
}

func (m *MockBedRepository) EXPECT() *MockBedRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockBedRepository) Create(bed *domain.Bed) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", bed)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockBedRepositoryMockRecorder) Create(bed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBedRepository)(nil).Create), bed)
}

func (m *MockBedRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.BedFilter) ([]domain.Bed, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	beds, _ := ret[0].([]domain.Bed)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return beds, total, err
}

func (mr *MockBedRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockBedRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockBedRepository) FindByIDAndTenant(id uint, tenantID uuid.UUID) (*domain.Bed, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDAndTenant", id, tenantID)
	bed, _ := ret[0].(*domain.Bed)
	err, _ := ret[1].(error)
	return bed, err
}

func (mr *MockBedRepositoryMockRecorder) FindByIDAndTenant(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDAndTenant", reflect.TypeOf((*MockBedRepository)(nil).FindByIDAndTenant), id, tenantID)
}

func (m *MockBedRepository) Update(bed *domain.Bed) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", bed)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockBedRepositoryMockRecorder) Update(bed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockBedRepository)(nil).Update), bed)
}

func (m *MockBedRepository) Delete(bed *domain.Bed) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", bed)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockBedRepositoryMockRecorder) Delete(bed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBedRepository)(nil).Delete), bed)
}

func (m *MockBedRepository) ExistsByCode(tenantID uuid.UUID, roomID uuid.UUID, code string, excludeID *uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByCode", tenantID, roomID, code, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockBedRepositoryMockRecorder) ExistsByCode(tenantID, roomID, code, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByCode", reflect.TypeOf((*MockBedRepository)(nil).ExistsByCode), tenantID, roomID, code, excludeID)
}

func (m *MockBedRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockBedRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockBedRepository)(nil).CountByTenant), tenantID)
}

type MockCategoryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryRepositoryMockRecorder
}

type MockCategoryRepositoryMockRecorder struct {
	mock *MockCategoryRepository
}

func NewMockCategoryRepository(ctrl *gomock.Controller) *MockCategoryRepository {
	return &MockCategoryRepository{ctrl: ctrl, recorder: &MockCategoryRepositoryMockRecorder{}}
}

func (m *MockCategoryRepository) EXPECT() *MockCategoryRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", category)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockCategoryRepositoryMockRecorder) Create(category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCategoryRepository)(nil).Create), category)
}

func (m *MockCategoryRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.CategoryFilter) ([]domain.Category, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	categories, _ := ret[0].([]domain.Category)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return categories, total, err
}

func (mr *MockCategoryRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockCategoryRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockCategoryRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	category, _ := ret[0].(*domain.Category)
	err, _ := ret[1].(error)
	return category, err
}

func (mr *MockCategoryRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockCategoryRepository)(nil).FindByID), tenantID, id)
}

func (m *MockCategoryRepository) Update(category *domain.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", category)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockCategoryRepositoryMockRecorder) Update(category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCategoryRepository)(nil).Update), category)
}

func (m *MockCategoryRepository) Delete(category *domain.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", category)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockCategoryRepositoryMockRecorder) Delete(category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCategoryRepository)(nil).Delete), category)
}

func (m *MockCategoryRepository) ExistsBySlug(tenantID uuid.UUID, slug string, excludeID *uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsBySlug", tenantID, slug, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockCategoryRepositoryMockRecorder) ExistsBySlug(tenantID, slug, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsBySlug", reflect.TypeOf((*MockCategoryRepository)(nil).ExistsBySlug), tenantID, slug, excludeID)
}

func (m *MockCategoryRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockCategoryRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockCategoryRepository)(nil).CountByTenant), tenantID)
}

type MockAssetRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAssetRepositoryMockRecorder
}

type MockAssetRepositoryMockRecorder struct {
	mock *MockAssetRepository
}

func NewMockAssetRepository(ctrl *gomock.Controller) *MockAssetRepository {
	return &MockAssetRepository{ctrl: ctrl, recorder: &MockAssetRepositoryMockRecorder{}}
}

func (m *MockAssetRepository) EXPECT() *MockAssetRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockAssetRepository) Create(asset *domain.Asset) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", asset)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetRepositoryMockRecorder) Create(asset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAssetRepository)(nil).Create), asset)
}

func (m *MockAssetRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.AssetFilter) ([]domain.Asset, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	assets, _ := ret[0].([]domain.Asset)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return assets, total, err
}

func (mr *MockAssetRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockAssetRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockAssetRepository) FindAllByTenantExport(tenantID uuid.UUID, filter repository.AssetFilter) ([]domain.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenantExport", tenantID, filter)
	assets, _ := ret[0].([]domain.Asset)
	err, _ := ret[1].(error)
	return assets, err
}

func (mr *MockAssetRepositoryMockRecorder) FindAllByTenantExport(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenantExport", reflect.TypeOf((*MockAssetRepository)(nil).FindAllByTenantExport), tenantID, filter)
}

func (m *MockAssetRepository) FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDAndTenant", id, tenantID)
	asset, _ := ret[0].(*domain.Asset)
	err, _ := ret[1].(error)
	return asset, err
}

func (mr *MockAssetRepositoryMockRecorder) FindByIDAndTenant(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDAndTenant", reflect.TypeOf((*MockAssetRepository)(nil).FindByIDAndTenant), id, tenantID)
}

func (m *MockAssetRepository) FindByCodeAndTenant(code string, tenantID uuid.UUID) (*domain.Asset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCodeAndTenant", code, tenantID)
	asset, _ := ret[0].(*domain.Asset)
	err, _ := ret[1].(error)
	return asset, err
}

func (mr *MockAssetRepositoryMockRecorder) FindByCodeAndTenant(code, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCodeAndTenant", reflect.TypeOf((*MockAssetRepository)(nil).FindByCodeAndTenant), code, tenantID)
}

func (m *MockAssetRepository) Update(asset *domain.Asset) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", asset)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetRepositoryMockRecorder) Update(asset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAssetRepository)(nil).Update), asset)
}

func (m *MockAssetRepository) Delete(asset *domain.Asset) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", asset)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetRepositoryMockRecorder) Delete(asset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAssetRepository)(nil).Delete), asset)
}

func (m *MockAssetRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByCode", tenantID, code, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockAssetRepositoryMockRecorder) ExistsByCode(tenantID, code, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByCode", reflect.TypeOf((*MockAssetRepository)(nil).ExistsByCode), tenantID, code, excludeID)
}

func (m *MockAssetRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockAssetRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockAssetRepository)(nil).CountByTenant), tenantID)
}

type MockVendorRepository struct {
	ctrl     *gomock.Controller
	recorder *MockVendorRepositoryMockRecorder
}

type MockVendorRepositoryMockRecorder struct {
	mock *MockVendorRepository
}

func NewMockVendorRepository(ctrl *gomock.Controller) *MockVendorRepository {
	return &MockVendorRepository{ctrl: ctrl, recorder: &MockVendorRepositoryMockRecorder{}}
}

func (m *MockVendorRepository) EXPECT() *MockVendorRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockVendorRepository) Create(vendor *domain.Vendor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", vendor)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockVendorRepositoryMockRecorder) Create(vendor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVendorRepository)(nil).Create), vendor)
}

func (m *MockVendorRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.VendorFilter) ([]domain.Vendor, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	vendors, _ := ret[0].([]domain.Vendor)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return vendors, total, err
}

func (mr *MockVendorRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockVendorRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockVendorRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.Vendor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	vendor, _ := ret[0].(*domain.Vendor)
	err, _ := ret[1].(error)
	return vendor, err
}

func (mr *MockVendorRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockVendorRepository)(nil).FindByID), tenantID, id)
}

func (m *MockVendorRepository) Update(vendor *domain.Vendor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", vendor)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockVendorRepositoryMockRecorder) Update(vendor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockVendorRepository)(nil).Update), vendor)
}

func (m *MockVendorRepository) Delete(vendor *domain.Vendor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", vendor)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockVendorRepositoryMockRecorder) Delete(vendor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockVendorRepository)(nil).Delete), vendor)
}

func (m *MockVendorRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByCode", tenantID, code, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockVendorRepositoryMockRecorder) ExistsByCode(tenantID, code, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByCode", reflect.TypeOf((*MockVendorRepository)(nil).ExistsByCode), tenantID, code, excludeID)
}

func (m *MockVendorRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockVendorRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockVendorRepository)(nil).CountByTenant), tenantID)
}

type MockAssetBrandRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAssetBrandRepositoryMockRecorder
}

type MockAssetBrandRepositoryMockRecorder struct {
	mock *MockAssetBrandRepository
}

func NewMockAssetBrandRepository(ctrl *gomock.Controller) *MockAssetBrandRepository {
	return &MockAssetBrandRepository{ctrl: ctrl, recorder: &MockAssetBrandRepositoryMockRecorder{}}
}

func (m *MockAssetBrandRepository) EXPECT() *MockAssetBrandRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockAssetBrandRepository) Create(brand *domain.AssetBrand) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", brand)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetBrandRepositoryMockRecorder) Create(brand interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAssetBrandRepository)(nil).Create), brand)
}

func (m *MockAssetBrandRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.AssetBrandFilter) ([]domain.AssetBrand, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	brands, _ := ret[0].([]domain.AssetBrand)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return brands, total, err
}

func (mr *MockAssetBrandRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockAssetBrandRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockAssetBrandRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AssetBrand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	brand, _ := ret[0].(*domain.AssetBrand)
	err, _ := ret[1].(error)
	return brand, err
}

func (mr *MockAssetBrandRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockAssetBrandRepository)(nil).FindByID), tenantID, id)
}

func (m *MockAssetBrandRepository) Update(brand *domain.AssetBrand) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", brand)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetBrandRepositoryMockRecorder) Update(brand interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAssetBrandRepository)(nil).Update), brand)
}

func (m *MockAssetBrandRepository) Delete(brand *domain.AssetBrand) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", brand)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetBrandRepositoryMockRecorder) Delete(brand interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAssetBrandRepository)(nil).Delete), brand)
}

func (m *MockAssetBrandRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByCode", tenantID, code, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockAssetBrandRepositoryMockRecorder) ExistsByCode(tenantID, code, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByCode", reflect.TypeOf((*MockAssetBrandRepository)(nil).ExistsByCode), tenantID, code, excludeID)
}

func (m *MockAssetBrandRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockAssetBrandRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockAssetBrandRepository)(nil).CountByTenant), tenantID)
}

type MockAssetModelRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAssetModelRepositoryMockRecorder
}

type MockAssetModelRepositoryMockRecorder struct {
	mock *MockAssetModelRepository
}

func NewMockAssetModelRepository(ctrl *gomock.Controller) *MockAssetModelRepository {
	return &MockAssetModelRepository{ctrl: ctrl, recorder: &MockAssetModelRepositoryMockRecorder{}}
}

func (m *MockAssetModelRepository) EXPECT() *MockAssetModelRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockAssetModelRepository) Create(model *domain.AssetModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", model)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetModelRepositoryMockRecorder) Create(model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAssetModelRepository)(nil).Create), model)
}

func (m *MockAssetModelRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.AssetModelFilter) ([]domain.AssetModel, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	models, _ := ret[0].([]domain.AssetModel)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return models, total, err
}

func (mr *MockAssetModelRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockAssetModelRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockAssetModelRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AssetModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	model, _ := ret[0].(*domain.AssetModel)
	err, _ := ret[1].(error)
	return model, err
}

func (mr *MockAssetModelRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockAssetModelRepository)(nil).FindByID), tenantID, id)
}

func (m *MockAssetModelRepository) Update(model *domain.AssetModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", model)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetModelRepositoryMockRecorder) Update(model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAssetModelRepository)(nil).Update), model)
}

func (m *MockAssetModelRepository) Delete(model *domain.AssetModel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", model)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetModelRepositoryMockRecorder) Delete(model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAssetModelRepository)(nil).Delete), model)
}

func (m *MockAssetModelRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByCode", tenantID, code, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockAssetModelRepositoryMockRecorder) ExistsByCode(tenantID, code, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByCode", reflect.TypeOf((*MockAssetModelRepository)(nil).ExistsByCode), tenantID, code, excludeID)
}

func (m *MockAssetModelRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockAssetModelRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockAssetModelRepository)(nil).CountByTenant), tenantID)
}

type MockAssetStatusRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAssetStatusRepositoryMockRecorder
}

type MockAssetStatusRepositoryMockRecorder struct {
	mock *MockAssetStatusRepository
}

func NewMockAssetStatusRepository(ctrl *gomock.Controller) *MockAssetStatusRepository {
	return &MockAssetStatusRepository{ctrl: ctrl, recorder: &MockAssetStatusRepositoryMockRecorder{}}
}

func (m *MockAssetStatusRepository) EXPECT() *MockAssetStatusRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockAssetStatusRepository) Create(status *domain.AssetStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", status)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetStatusRepositoryMockRecorder) Create(status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAssetStatusRepository)(nil).Create), status)
}

func (m *MockAssetStatusRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.AssetStatusFilter) ([]domain.AssetStatus, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	statuses, _ := ret[0].([]domain.AssetStatus)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return statuses, total, err
}

func (mr *MockAssetStatusRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockAssetStatusRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockAssetStatusRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AssetStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	status, _ := ret[0].(*domain.AssetStatus)
	err, _ := ret[1].(error)
	return status, err
}

func (mr *MockAssetStatusRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockAssetStatusRepository)(nil).FindByID), tenantID, id)
}

func (m *MockAssetStatusRepository) Update(status *domain.AssetStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", status)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetStatusRepositoryMockRecorder) Update(status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAssetStatusRepository)(nil).Update), status)
}

func (m *MockAssetStatusRepository) Delete(status *domain.AssetStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", status)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetStatusRepositoryMockRecorder) Delete(status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAssetStatusRepository)(nil).Delete), status)
}

func (m *MockAssetStatusRepository) ExistsByCode(tenantID uuid.UUID, code string, excludeID *uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByCode", tenantID, code, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockAssetStatusRepositoryMockRecorder) ExistsByCode(tenantID, code, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsByCode", reflect.TypeOf((*MockAssetStatusRepository)(nil).ExistsByCode), tenantID, code, excludeID)
}

func (m *MockAssetStatusRepository) HasAny(tenantID uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasAny", tenantID)
	hasAny, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return hasAny, err
}

func (mr *MockAssetStatusRepositoryMockRecorder) HasAny(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasAny", reflect.TypeOf((*MockAssetStatusRepository)(nil).HasAny), tenantID)
}

func (m *MockAssetStatusRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockAssetStatusRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockAssetStatusRepository)(nil).CountByTenant), tenantID)
}

type MockAssetEventRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAssetEventRepositoryMockRecorder
}

type MockAssetEventRepositoryMockRecorder struct {
	mock *MockAssetEventRepository
}

func NewMockAssetEventRepository(ctrl *gomock.Controller) *MockAssetEventRepository {
	return &MockAssetEventRepository{ctrl: ctrl, recorder: &MockAssetEventRepositoryMockRecorder{}}
}

func (m *MockAssetEventRepository) EXPECT() *MockAssetEventRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockAssetEventRepository) Create(event *domain.AssetEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", event)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAssetEventRepositoryMockRecorder) Create(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAssetEventRepository)(nil).Create), event)
}

func (m *MockAssetEventRepository) ListByAsset(tenantID uuid.UUID, assetID uuid.UUID, filter repository.AssetEventFilter) ([]domain.AssetEvent, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByAsset", tenantID, assetID, filter)
	events, _ := ret[0].([]domain.AssetEvent)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return events, total, err
}

func (mr *MockAssetEventRepositoryMockRecorder) ListByAsset(tenantID, assetID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByAsset", reflect.TypeOf((*MockAssetEventRepository)(nil).ListByAsset), tenantID, assetID, filter)
}

type MockMaintenanceScheduleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMaintenanceScheduleRepositoryMockRecorder
}

type MockMaintenanceScheduleRepositoryMockRecorder struct {
	mock *MockMaintenanceScheduleRepository
}

func NewMockMaintenanceScheduleRepository(ctrl *gomock.Controller) *MockMaintenanceScheduleRepository {
	return &MockMaintenanceScheduleRepository{ctrl: ctrl, recorder: &MockMaintenanceScheduleRepositoryMockRecorder{}}
}

func (m *MockMaintenanceScheduleRepository) EXPECT() *MockMaintenanceScheduleRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockMaintenanceScheduleRepository) Create(schedule *domain.MaintenanceSchedule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", schedule)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockMaintenanceScheduleRepositoryMockRecorder) Create(schedule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMaintenanceScheduleRepository)(nil).Create), schedule)
}

func (m *MockMaintenanceScheduleRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	schedules, _ := ret[0].([]domain.MaintenanceSchedule)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return schedules, total, err
}

func (mr *MockMaintenanceScheduleRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockMaintenanceScheduleRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockMaintenanceScheduleRepository) FindAllByTenantExport(tenantID uuid.UUID, filter repository.MaintenanceScheduleFilter) ([]domain.MaintenanceSchedule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenantExport", tenantID, filter)
	schedules, _ := ret[0].([]domain.MaintenanceSchedule)
	err, _ := ret[1].(error)
	return schedules, err
}

func (mr *MockMaintenanceScheduleRepositoryMockRecorder) FindAllByTenantExport(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenantExport", reflect.TypeOf((*MockMaintenanceScheduleRepository)(nil).FindAllByTenantExport), tenantID, filter)
}

func (m *MockMaintenanceScheduleRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceSchedule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	schedule, _ := ret[0].(*domain.MaintenanceSchedule)
	err, _ := ret[1].(error)
	return schedule, err
}

func (mr *MockMaintenanceScheduleRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockMaintenanceScheduleRepository)(nil).FindByID), tenantID, id)
}

func (m *MockMaintenanceScheduleRepository) Update(schedule *domain.MaintenanceSchedule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", schedule)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockMaintenanceScheduleRepositoryMockRecorder) Update(schedule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMaintenanceScheduleRepository)(nil).Update), schedule)
}

func (m *MockMaintenanceScheduleRepository) Delete(schedule *domain.MaintenanceSchedule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", schedule)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockMaintenanceScheduleRepositoryMockRecorder) Delete(schedule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMaintenanceScheduleRepository)(nil).Delete), schedule)
}

type MockMaintenanceDocumentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMaintenanceDocumentRepositoryMockRecorder
}

type MockMaintenanceDocumentRepositoryMockRecorder struct {
	mock *MockMaintenanceDocumentRepository
}

func NewMockMaintenanceDocumentRepository(ctrl *gomock.Controller) *MockMaintenanceDocumentRepository {
	return &MockMaintenanceDocumentRepository{ctrl: ctrl, recorder: &MockMaintenanceDocumentRepositoryMockRecorder{}}
}

func (m *MockMaintenanceDocumentRepository) EXPECT() *MockMaintenanceDocumentRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockMaintenanceDocumentRepository) Create(doc *domain.MaintenanceDocument) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", doc)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockMaintenanceDocumentRepositoryMockRecorder) Create(doc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMaintenanceDocumentRepository)(nil).Create), doc)
}

func (m *MockMaintenanceDocumentRepository) ListBySchedule(tenantID uuid.UUID, scheduleID uint) ([]domain.MaintenanceDocument, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBySchedule", tenantID, scheduleID)
	docs, _ := ret[0].([]domain.MaintenanceDocument)
	err, _ := ret[1].(error)
	return docs, err
}

func (mr *MockMaintenanceDocumentRepositoryMockRecorder) ListBySchedule(tenantID, scheduleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBySchedule", reflect.TypeOf((*MockMaintenanceDocumentRepository)(nil).ListBySchedule), tenantID, scheduleID)
}

func (m *MockMaintenanceDocumentRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.MaintenanceDocument, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	doc, _ := ret[0].(*domain.MaintenanceDocument)
	err, _ := ret[1].(error)
	return doc, err
}

func (mr *MockMaintenanceDocumentRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockMaintenanceDocumentRepository)(nil).FindByID), tenantID, id)
}

func (m *MockMaintenanceDocumentRepository) Delete(doc *domain.MaintenanceDocument) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", doc)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockMaintenanceDocumentRepositoryMockRecorder) Delete(doc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMaintenanceDocumentRepository)(nil).Delete), doc)
}

type MockDocumentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDocumentRepositoryMockRecorder
}

type MockDocumentRepositoryMockRecorder struct {
	mock *MockDocumentRepository
}

func NewMockDocumentRepository(ctrl *gomock.Controller) *MockDocumentRepository {
	return &MockDocumentRepository{ctrl: ctrl, recorder: &MockDocumentRepositoryMockRecorder{}}
}

func (m *MockDocumentRepository) EXPECT() *MockDocumentRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockDocumentRepository) Create(doc *domain.Document) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", doc)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockDocumentRepositoryMockRecorder) Create(doc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDocumentRepository)(nil).Create), doc)
}

func (m *MockDocumentRepository) Update(doc *domain.Document) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", doc)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockDocumentRepositoryMockRecorder) Update(doc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDocumentRepository)(nil).Update), doc)
}

func (m *MockDocumentRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.Document, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	doc, _ := ret[0].(*domain.Document)
	err, _ := ret[1].(error)
	return doc, err
}

func (mr *MockDocumentRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockDocumentRepository)(nil).FindByID), tenantID, id)
}

func (m *MockDocumentRepository) List(tenantID uuid.UUID, filter repository.DocumentFilter) ([]domain.Document, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", tenantID, filter)
	docs, _ := ret[0].([]domain.Document)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return docs, total, err
}

func (mr *MockDocumentRepositoryMockRecorder) List(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDocumentRepository)(nil).List), tenantID, filter)
}

func (m *MockDocumentRepository) Delete(doc *domain.Document) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", doc)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockDocumentRepositoryMockRecorder) Delete(doc interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDocumentRepository)(nil).Delete), doc)
}

type MockDocumentFileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDocumentFileRepositoryMockRecorder
}

type MockDocumentFileRepositoryMockRecorder struct {
	mock *MockDocumentFileRepository
}

func NewMockDocumentFileRepository(ctrl *gomock.Controller) *MockDocumentFileRepository {
	return &MockDocumentFileRepository{ctrl: ctrl, recorder: &MockDocumentFileRepositoryMockRecorder{}}
}

func (m *MockDocumentFileRepository) EXPECT() *MockDocumentFileRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockDocumentFileRepository) CreateMany(files []domain.DocumentFile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMany", files)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockDocumentFileRepositoryMockRecorder) CreateMany(files interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMany", reflect.TypeOf((*MockDocumentFileRepository)(nil).CreateMany), files)
}

func (m *MockDocumentFileRepository) ListByDocument(tenantID uuid.UUID, documentID uint) ([]domain.DocumentFile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByDocument", tenantID, documentID)
	files, _ := ret[0].([]domain.DocumentFile)
	err, _ := ret[1].(error)
	return files, err
}

func (mr *MockDocumentFileRepositoryMockRecorder) ListByDocument(tenantID, documentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByDocument", reflect.TypeOf((*MockDocumentFileRepository)(nil).ListByDocument), tenantID, documentID)
}

func (m *MockDocumentFileRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.DocumentFile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	file, _ := ret[0].(*domain.DocumentFile)
	err, _ := ret[1].(error)
	return file, err
}

func (mr *MockDocumentFileRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockDocumentFileRepository)(nil).FindByID), tenantID, id)
}

func (m *MockDocumentFileRepository) Delete(file *domain.DocumentFile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", file)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockDocumentFileRepositoryMockRecorder) Delete(file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDocumentFileRepository)(nil).Delete), file)
}

func (m *MockDocumentFileRepository) DeleteByDocument(tenantID uuid.UUID, documentID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByDocument", tenantID, documentID)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockDocumentFileRepositoryMockRecorder) DeleteByDocument(tenantID, documentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByDocument", reflect.TypeOf((*MockDocumentFileRepository)(nil).DeleteByDocument), tenantID, documentID)
}

type MockComplaintRepository struct {
	ctrl     *gomock.Controller
	recorder *MockComplaintRepositoryMockRecorder
}

type MockComplaintRepositoryMockRecorder struct {
	mock *MockComplaintRepository
}

func NewMockComplaintRepository(ctrl *gomock.Controller) *MockComplaintRepository {
	return &MockComplaintRepository{ctrl: ctrl, recorder: &MockComplaintRepositoryMockRecorder{}}
}

func (m *MockComplaintRepository) EXPECT() *MockComplaintRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockComplaintRepository) Create(complaint *domain.Complaint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", complaint)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockComplaintRepositoryMockRecorder) Create(complaint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockComplaintRepository)(nil).Create), complaint)
}

func (m *MockComplaintRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.ComplaintFilter) ([]domain.Complaint, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	complaints, _ := ret[0].([]domain.Complaint)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return complaints, total, err
}

func (mr *MockComplaintRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockComplaintRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockComplaintRepository) FindAllByTenantExport(tenantID uuid.UUID, filter repository.ComplaintFilter) ([]domain.Complaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenantExport", tenantID, filter)
	complaints, _ := ret[0].([]domain.Complaint)
	err, _ := ret[1].(error)
	return complaints, err
}

func (mr *MockComplaintRepositoryMockRecorder) FindAllByTenantExport(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenantExport", reflect.TypeOf((*MockComplaintRepository)(nil).FindAllByTenantExport), tenantID, filter)
}

func (m *MockComplaintRepository) FindByIDAndTenant(id uuid.UUID, tenantID uuid.UUID) (*domain.Complaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDAndTenant", id, tenantID)
	complaint, _ := ret[0].(*domain.Complaint)
	err, _ := ret[1].(error)
	return complaint, err
}

func (mr *MockComplaintRepositoryMockRecorder) FindByIDAndTenant(id, tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDAndTenant", reflect.TypeOf((*MockComplaintRepository)(nil).FindByIDAndTenant), id, tenantID)
}

func (m *MockComplaintRepository) Update(complaint *domain.Complaint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", complaint)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockComplaintRepositoryMockRecorder) Update(complaint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockComplaintRepository)(nil).Update), complaint)
}

func (m *MockComplaintRepository) Delete(complaint *domain.Complaint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", complaint)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockComplaintRepositoryMockRecorder) Delete(complaint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockComplaintRepository)(nil).Delete), complaint)
}

type MockStockOpnameRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStockOpnameRepositoryMockRecorder
}

type MockStockOpnameRepositoryMockRecorder struct {
	mock *MockStockOpnameRepository
}

func NewMockStockOpnameRepository(ctrl *gomock.Controller) *MockStockOpnameRepository {
	return &MockStockOpnameRepository{ctrl: ctrl, recorder: &MockStockOpnameRepositoryMockRecorder{}}
}

func (m *MockStockOpnameRepository) EXPECT() *MockStockOpnameRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockStockOpnameRepository) Create(session *domain.StockOpnameSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", session)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockStockOpnameRepositoryMockRecorder) Create(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStockOpnameRepository)(nil).Create), session)
}

func (m *MockStockOpnameRepository) FindAllByTenant(tenantID uuid.UUID, filter repository.StockOpnameFilter) ([]domain.StockOpnameSession, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByTenant", tenantID, filter)
	sessions, _ := ret[0].([]domain.StockOpnameSession)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return sessions, total, err
}

func (mr *MockStockOpnameRepositoryMockRecorder) FindAllByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByTenant", reflect.TypeOf((*MockStockOpnameRepository)(nil).FindAllByTenant), tenantID, filter)
}

func (m *MockStockOpnameRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.StockOpnameSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	session, _ := ret[0].(*domain.StockOpnameSession)
	err, _ := ret[1].(error)
	return session, err
}

func (mr *MockStockOpnameRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockStockOpnameRepository)(nil).FindByID), tenantID, id)
}

func (m *MockStockOpnameRepository) Update(session *domain.StockOpnameSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", session)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockStockOpnameRepositoryMockRecorder) Update(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStockOpnameRepository)(nil).Update), session)
}

func (m *MockStockOpnameRepository) Delete(session *domain.StockOpnameSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", session)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockStockOpnameRepositoryMockRecorder) Delete(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStockOpnameRepository)(nil).Delete), session)
}

type MockStockOpnameItemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStockOpnameItemRepositoryMockRecorder
}

type MockStockOpnameItemRepositoryMockRecorder struct {
	mock *MockStockOpnameItemRepository
}

func NewMockStockOpnameItemRepository(ctrl *gomock.Controller) *MockStockOpnameItemRepository {
	return &MockStockOpnameItemRepository{ctrl: ctrl, recorder: &MockStockOpnameItemRepositoryMockRecorder{}}
}

func (m *MockStockOpnameItemRepository) EXPECT() *MockStockOpnameItemRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockStockOpnameItemRepository) Create(item *domain.StockOpnameItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", item)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockStockOpnameItemRepositoryMockRecorder) Create(item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStockOpnameItemRepository)(nil).Create), item)
}

func (m *MockStockOpnameItemRepository) ListBySession(tenantID uuid.UUID, sessionID uint) ([]domain.StockOpnameItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBySession", tenantID, sessionID)
	items, _ := ret[0].([]domain.StockOpnameItem)
	err, _ := ret[1].(error)
	return items, err
}

func (mr *MockStockOpnameItemRepositoryMockRecorder) ListBySession(tenantID, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBySession", reflect.TypeOf((*MockStockOpnameItemRepository)(nil).ListBySession), tenantID, sessionID)
}

type MockTenantRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTenantRepositoryMockRecorder
}

type MockTenantRepositoryMockRecorder struct {
	mock *MockTenantRepository
}

func NewMockTenantRepository(ctrl *gomock.Controller) *MockTenantRepository {
	return &MockTenantRepository{ctrl: ctrl, recorder: &MockTenantRepositoryMockRecorder{}}
}

func (m *MockTenantRepository) EXPECT() *MockTenantRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockTenantRepository) ListAll(filter repository.TenantFilter) ([]domain.Tenant, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAll", filter)
	tenants, _ := ret[0].([]domain.Tenant)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return tenants, total, err
}

func (mr *MockTenantRepositoryMockRecorder) ListAll(filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAll", reflect.TypeOf((*MockTenantRepository)(nil).ListAll), filter)
}

func (m *MockTenantRepository) FindByID(id uuid.UUID) (*domain.Tenant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", id)
	tenant, _ := ret[0].(*domain.Tenant)
	err, _ := ret[1].(error)
	return tenant, err
}

func (mr *MockTenantRepositoryMockRecorder) FindByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockTenantRepository)(nil).FindByID), id)
}

func (m *MockTenantRepository) FindBySlug(slug string) (*domain.Tenant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBySlug", slug)
	tenant, _ := ret[0].(*domain.Tenant)
	err, _ := ret[1].(error)
	return tenant, err
}

func (mr *MockTenantRepositoryMockRecorder) FindBySlug(slug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBySlug", reflect.TypeOf((*MockTenantRepository)(nil).FindBySlug), slug)
}

func (m *MockTenantRepository) Create(tenant *domain.Tenant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", tenant)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockTenantRepositoryMockRecorder) Create(tenant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTenantRepository)(nil).Create), tenant)
}

func (m *MockTenantRepository) Update(tenant *domain.Tenant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", tenant)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockTenantRepositoryMockRecorder) Update(tenant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTenantRepository)(nil).Update), tenant)
}

func (m *MockTenantRepository) ExistsBySlug(slug string, excludeID *uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsBySlug", slug, excludeID)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (mr *MockTenantRepositoryMockRecorder) ExistsBySlug(slug, excludeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsBySlug", reflect.TypeOf((*MockTenantRepository)(nil).ExistsBySlug), slug, excludeID)
}

type MockAuditTrailRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuditTrailRepositoryMockRecorder
}

type MockAuditTrailRepositoryMockRecorder struct {
	mock *MockAuditTrailRepository
}

func NewMockAuditTrailRepository(ctrl *gomock.Controller) *MockAuditTrailRepository {
	return &MockAuditTrailRepository{ctrl: ctrl, recorder: &MockAuditTrailRepositoryMockRecorder{}}
}

func (m *MockAuditTrailRepository) EXPECT() *MockAuditTrailRepositoryMockRecorder {
	m.recorder.mock = m
	return m.recorder
}

func (m *MockAuditTrailRepository) Create(audit *domain.AuditTrail) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", audit)
	err, _ := ret[0].(error)
	return err
}

func (mr *MockAuditTrailRepositoryMockRecorder) Create(audit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAuditTrailRepository)(nil).Create), audit)
}

func (m *MockAuditTrailRepository) ListByTenant(tenantID uuid.UUID, filter repository.AuditTrailFilter) ([]domain.AuditTrail, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByTenant", tenantID, filter)
	items, _ := ret[0].([]domain.AuditTrail)
	total, _ := ret[1].(int64)
	err, _ := ret[2].(error)
	return items, total, err
}

func (mr *MockAuditTrailRepositoryMockRecorder) ListByTenant(tenantID, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByTenant", reflect.TypeOf((*MockAuditTrailRepository)(nil).ListByTenant), tenantID, filter)
}

func (m *MockAuditTrailRepository) FindByID(tenantID uuid.UUID, id uint) (*domain.AuditTrail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", tenantID, id)
	item, _ := ret[0].(*domain.AuditTrail)
	err, _ := ret[1].(error)
	return item, err
}

func (mr *MockAuditTrailRepositoryMockRecorder) FindByID(tenantID, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockAuditTrailRepository)(nil).FindByID), tenantID, id)
}

func (m *MockAuditTrailRepository) CountByTenant(tenantID uuid.UUID) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByTenant", tenantID)
	count, _ := ret[0].(int64)
	err, _ := ret[1].(error)
	return count, err
}

func (mr *MockAuditTrailRepositoryMockRecorder) CountByTenant(tenantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByTenant", reflect.TypeOf((*MockAuditTrailRepository)(nil).CountByTenant), tenantID)
}
