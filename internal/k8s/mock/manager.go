// Code generated by MockGen. DO NOT EDIT.
// Source: ./manager.go

// Package mock_k8s is a generated GoMock package.
package mock_k8s

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	k8s  "github.com/spirt-t/lg-operator/internal/k8s"
	model  "github.com/spirt-t/lg-operator/internal/model"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockManager) Create(ctx context.Context, cfg k8s.CreationConfig) (*model.LoadGenerator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, cfg)
	ret0, _ := ret[0].(*model.LoadGenerator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockManagerMockRecorder) Create(ctx, cfg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockManager)(nil).Create), ctx, cfg)
}

// Delete mocks base method.
func (m *MockManager) Delete(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockManagerMockRecorder) Delete(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockManager)(nil).Delete), ctx, name)
}

// DeleteAll mocks base method.
func (m *MockManager) DeleteAll(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAll", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll.
func (mr *MockManagerMockRecorder) DeleteAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockManager)(nil).DeleteAll), ctx)
}

// List mocks base method.
func (m *MockManager) List(ctx context.Context) ([]model.LoadGenerator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx)
	ret0, _ := ret[0].([]model.LoadGenerator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockManagerMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockManager)(nil).List), ctx)
}
