// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/services/valkey (interfaces: Service)
//
// Generated by this command:
//
//	mockgen -destination ./mocks/mock_valkey_service.go -package mocks -mock_names Service=MockValkeyService ./internal/services/valkey Service
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	v1alpha1 "github.com/uagolang/k8s-operator/api/v1alpha1"
	valkey "github.com/uagolang/k8s-operator/internal/services/valkey"
	gomock "go.uber.org/mock/gomock"
	types "k8s.io/apimachinery/pkg/types"
)

// MockValkeyService is a mock of Service interface.
type MockValkeyService struct {
	ctrl     *gomock.Controller
	recorder *MockValkeyServiceMockRecorder
	isgomock struct{}
}

// MockValkeyServiceMockRecorder is the mock recorder for MockValkeyService.
type MockValkeyServiceMockRecorder struct {
	mock *MockValkeyService
}

// NewMockValkeyService creates a new mock instance.
func NewMockValkeyService(ctrl *gomock.Controller) *MockValkeyService {
	mock := &MockValkeyService{ctrl: ctrl}
	mock.recorder = &MockValkeyServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValkeyService) EXPECT() *MockValkeyServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockValkeyService) Create(ctx context.Context, i *valkey.CreateRequest) (*v1alpha1.Valkey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, i)
	ret0, _ := ret[0].(*v1alpha1.Valkey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockValkeyServiceMockRecorder) Create(ctx, i any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockValkeyService)(nil).Create), ctx, i)
}

// Delete mocks base method.
func (m *MockValkeyService) Delete(ctx context.Context, i types.NamespacedName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, i)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockValkeyServiceMockRecorder) Delete(ctx, i any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockValkeyService)(nil).Delete), ctx, i)
}

// IsReady mocks base method.
func (m *MockValkeyService) IsReady(ctx context.Context, item *v1alpha1.Valkey) (bool, int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsReady", ctx, item)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(int32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IsReady indicates an expected call of IsReady.
func (mr *MockValkeyServiceMockRecorder) IsReady(ctx, item any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsReady", reflect.TypeOf((*MockValkeyService)(nil).IsReady), ctx, item)
}

// Update mocks base method.
func (m *MockValkeyService) Update(ctx context.Context, i *valkey.UpdateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, i)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockValkeyServiceMockRecorder) Update(ctx, i any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockValkeyService)(nil).Update), ctx, i)
}
