// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Viverov/guideliner/internal/domains/user/token_provider (interfaces: TokenProvider)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	tokenprovider "github.com/Viverov/guideliner/internal/domains/user/token_provider"
	gomock "github.com/golang/mock/gomock"
)

// MockTokenProvider is a mock of TokenProvider interface.
type MockTokenProvider struct {
	ctrl     *gomock.Controller
	recorder *MockTokenProviderMockRecorder
}

// MockTokenProviderMockRecorder is the mock recorder for MockTokenProvider.
type MockTokenProviderMockRecorder struct {
	mock *MockTokenProvider
}

// NewMockTokenProvider creates a new mock instance.
func NewMockTokenProvider(ctrl *gomock.Controller) *MockTokenProvider {
	mock := &MockTokenProvider{ctrl: ctrl}
	mock.recorder = &MockTokenProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenProvider) EXPECT() *MockTokenProviderMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockTokenProvider) GenerateToken(arg0 uint, arg1 time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockTokenProviderMockRecorder) GenerateToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockTokenProvider)(nil).GenerateToken), arg0, arg1)
}

// ValidateToken mocks base method.
func (m *MockTokenProvider) ValidateToken(arg0 string) (*tokenprovider.AuthClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", arg0)
	ret0, _ := ret[0].(*tokenprovider.AuthClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockTokenProviderMockRecorder) ValidateToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockTokenProvider)(nil).ValidateToken), arg0)
}
