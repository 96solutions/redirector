// Code generated by MockGen. DO NOT EDIT.
// Source: domain/interactor/redirect_interactor.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	net "net"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRedirectInteractor is a mock of RedirectInteractor interface.
type MockRedirectInteractor struct {
	ctrl     *gomock.Controller
	recorder *MockRedirectInteractorMockRecorder
}

// MockRedirectInteractorMockRecorder is the mock recorder for MockRedirectInteractor.
type MockRedirectInteractorMockRecorder struct {
	mock *MockRedirectInteractor
}

// NewMockRedirectInteractor creates a new mock instance.
func NewMockRedirectInteractor(ctrl *gomock.Controller) *MockRedirectInteractor {
	mock := &MockRedirectInteractor{ctrl: ctrl}
	mock.recorder = &MockRedirectInteractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedirectInteractor) EXPECT() *MockRedirectInteractorMockRecorder {
	return m.recorder
}

// Redirect mocks base method.
func (m *MockRedirectInteractor) Redirect(ctx context.Context, slug string, params map[string][]string, headers map[string]string, userAgent string, ip net.IP, protocol string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Redirect", ctx, slug, params, headers, userAgent, ip, protocol)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Redirect indicates an expected call of Redirect.
func (mr *MockRedirectInteractorMockRecorder) Redirect(ctx, slug, params, headers, userAgent, ip, protocol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Redirect", reflect.TypeOf((*MockRedirectInteractor)(nil).Redirect), ctx, slug, params, headers, userAgent, ip, protocol)
}
