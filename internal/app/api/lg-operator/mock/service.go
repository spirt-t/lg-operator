// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock_lg_operator is a generated GoMock package.
package mock_lg_operator

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCleaner is a mock of Cleaner interface.
type MockCleaner struct {
	ctrl     *gomock.Controller
	recorder *MockCleanerMockRecorder
}

// MockCleanerMockRecorder is the mock recorder for MockCleaner.
type MockCleanerMockRecorder struct {
	mock *MockCleaner
}

// NewMockCleaner creates a new mock instance.
func NewMockCleaner(ctrl *gomock.Controller) *MockCleaner {
	mock := &MockCleaner{ctrl: ctrl}
	mock.recorder = &MockCleanerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCleaner) EXPECT() *MockCleanerMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockCleaner) Run(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockCleanerMockRecorder) Run(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockCleaner)(nil).Run), ctx)
}
