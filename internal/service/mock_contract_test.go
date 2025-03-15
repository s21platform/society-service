// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	society_proto "github.com/s21platform/society-proto/society-proto"
	model "github.com/s21platform/society-service/internal/model"
)

// MockDbRepo is a mock of DbRepo interface.
type MockDbRepo struct {
	ctrl     *gomock.Controller
	recorder *MockDbRepoMockRecorder
}

// MockDbRepoMockRecorder is the mock recorder for MockDbRepo.
type MockDbRepoMockRecorder struct {
	mock *MockDbRepo
}

// NewMockDbRepo creates a new mock instance.
func NewMockDbRepo(ctrl *gomock.Controller) *MockDbRepo {
	mock := &MockDbRepo{ctrl: ctrl}
	mock.recorder = &MockDbRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDbRepo) EXPECT() *MockDbRepoMockRecorder {
	return m.recorder
}

// CountSubscribe mocks base method.
func (m *MockDbRepo) CountSubscribe(societyUUID string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountSubscribe", societyUUID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountSubscribe indicates an expected call of CountSubscribe.
func (mr *MockDbRepoMockRecorder) CountSubscribe(societyUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountSubscribe", reflect.TypeOf((*MockDbRepo)(nil).CountSubscribe), societyUUID)
}

// CreateSociety mocks base method.
func (m *MockDbRepo) CreateSociety(socData *model.SocietyData) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSociety", socData)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSociety indicates an expected call of CreateSociety.
func (mr *MockDbRepoMockRecorder) CreateSociety(socData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSociety", reflect.TypeOf((*MockDbRepo)(nil).CreateSociety), socData)
}

// GetSocietyInfo mocks base method.
func (m *MockDbRepo) GetSocietyInfo(societyUUID string) (*model.SocietyInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSocietyInfo", societyUUID)
	ret0, _ := ret[0].(*model.SocietyInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSocietyInfo indicates an expected call of GetSocietyInfo.
func (mr *MockDbRepoMockRecorder) GetSocietyInfo(societyUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSocietyInfo", reflect.TypeOf((*MockDbRepo)(nil).GetSocietyInfo), societyUUID)
}

// GetTags mocks base method.
func (m *MockDbRepo) GetTags(societyUUID string) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTags", societyUUID)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTags indicates an expected call of GetTags.
func (mr *MockDbRepoMockRecorder) GetTags(societyUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTags", reflect.TypeOf((*MockDbRepo)(nil).GetTags), societyUUID)
}

// IsOwnerAdminModerator mocks base method.
func (m *MockDbRepo) IsOwnerAdminModerator(peerUUID, societyUUID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOwnerAdminModerator", peerUUID, societyUUID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsOwnerAdminModerator indicates an expected call of IsOwnerAdminModerator.
func (mr *MockDbRepoMockRecorder) IsOwnerAdminModerator(peerUUID, societyUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOwnerAdminModerator", reflect.TypeOf((*MockDbRepo)(nil).IsOwnerAdminModerator), peerUUID, societyUUID)
}

// UpdateSociety mocks base method.
func (m *MockDbRepo) UpdateSociety(societyData *society_proto.UpdateSocietyIn) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSociety", societyData)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSociety indicates an expected call of UpdateSociety.
func (mr *MockDbRepoMockRecorder) UpdateSociety(societyData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSociety", reflect.TypeOf((*MockDbRepo)(nil).UpdateSociety), societyData)
}
