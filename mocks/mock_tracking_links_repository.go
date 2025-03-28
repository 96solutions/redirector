// Code generated by MockGen. DO NOT EDIT.
// Source: domain/repository/tracking_links_repository.go
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=mocks/mock_tracking_links_repository.go -source=domain/repository/tracking_links_repository.go TrackingLinksRepositoryInterface
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/lroman242/redirector/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockTrackingLinksRepositoryInterface is a mock of TrackingLinksRepositoryInterface interface.
type MockTrackingLinksRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTrackingLinksRepositoryInterfaceMockRecorder
	isgomock struct{}
}

// MockTrackingLinksRepositoryInterfaceMockRecorder is the mock recorder for MockTrackingLinksRepositoryInterface.
type MockTrackingLinksRepositoryInterfaceMockRecorder struct {
	mock *MockTrackingLinksRepositoryInterface
}

// NewMockTrackingLinksRepositoryInterface creates a new mock instance.
func NewMockTrackingLinksRepositoryInterface(ctrl *gomock.Controller) *MockTrackingLinksRepositoryInterface {
	mock := &MockTrackingLinksRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockTrackingLinksRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrackingLinksRepositoryInterface) EXPECT() *MockTrackingLinksRepositoryInterfaceMockRecorder {
	return m.recorder
}

// FindTrackingLink mocks base method.
func (m *MockTrackingLinksRepositoryInterface) FindTrackingLink(ctx context.Context, slug string) *entity.TrackingLink {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTrackingLink", ctx, slug)
	ret0, _ := ret[0].(*entity.TrackingLink)
	return ret0
}

// FindTrackingLink indicates an expected call of FindTrackingLink.
func (mr *MockTrackingLinksRepositoryInterfaceMockRecorder) FindTrackingLink(ctx, slug any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTrackingLink", reflect.TypeOf((*MockTrackingLinksRepositoryInterface)(nil).FindTrackingLink), ctx, slug)
}
