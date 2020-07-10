// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/charmhub (interfaces: Client)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	transport "github.com/juju/juju/charmhub/transport"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockClient) Find(arg0 context.Context, arg1 string) ([]transport.FindResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].([]transport.FindResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockClientMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockClient)(nil).Find), arg0, arg1)
}

// Info mocks base method
func (m *MockClient) Info(arg0 context.Context, arg1 string) (transport.InfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info", arg0, arg1)
	ret0, _ := ret[0].(transport.InfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Info indicates an expected call of Info
func (mr *MockClientMockRecorder) Info(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockClient)(nil).Info), arg0, arg1)
}

// URL mocks base method
func (m *MockClient) URL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "URL")
	ret0, _ := ret[0].(string)
	return ret0
}

// URL indicates an expected call of URL
func (mr *MockClientMockRecorder) URL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*MockClient)(nil).URL))
}
