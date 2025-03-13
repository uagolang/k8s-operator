// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/controller-runtime/pkg/client (interfaces: SubResourceWriter)
//
// Generated by this command:
//
//	mockgen -destination ./mocks/mock_k8s_status_client.go -package mocks -mock_names SubResourceWriter=MockK8sStatusClient sigs.k8s.io/controller-runtime/pkg/client SubResourceWriter
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// MockK8sStatusClient is a mock of SubResourceWriter interface.
type MockK8sStatusClient struct {
	ctrl     *gomock.Controller
	recorder *MockK8sStatusClientMockRecorder
	isgomock struct{}
}

// MockK8sStatusClientMockRecorder is the mock recorder for MockK8sStatusClient.
type MockK8sStatusClientMockRecorder struct {
	mock *MockK8sStatusClient
}

// NewMockK8sStatusClient creates a new mock instance.
func NewMockK8sStatusClient(ctrl *gomock.Controller) *MockK8sStatusClient {
	mock := &MockK8sStatusClient{ctrl: ctrl}
	mock.recorder = &MockK8sStatusClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockK8sStatusClient) EXPECT() *MockK8sStatusClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockK8sStatusClient) Create(ctx context.Context, obj, subResource client.Object, opts ...client.SubResourceCreateOption) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, obj, subResource}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockK8sStatusClientMockRecorder) Create(ctx, obj, subResource any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, obj, subResource}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockK8sStatusClient)(nil).Create), varargs...)
}

// Patch mocks base method.
func (m *MockK8sStatusClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.SubResourcePatchOption) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, obj, patch}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Patch indicates an expected call of Patch.
func (mr *MockK8sStatusClientMockRecorder) Patch(ctx, obj, patch any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, obj, patch}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockK8sStatusClient)(nil).Patch), varargs...)
}

// Update mocks base method.
func (m *MockK8sStatusClient) Update(ctx context.Context, obj client.Object, opts ...client.SubResourceUpdateOption) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, obj}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockK8sStatusClientMockRecorder) Update(ctx, obj any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, obj}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockK8sStatusClient)(nil).Update), varargs...)
}
