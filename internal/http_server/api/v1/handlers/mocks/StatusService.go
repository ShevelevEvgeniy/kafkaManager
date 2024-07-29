// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// StatusService is an autogenerated mock type for the StatusService type
type StatusService struct {
	mock.Mock
}

// GetStatus provides a mock function with given fields: ctx, requestId
func (_m *StatusService) GetStatus(ctx context.Context, requestId string) (string, error) {
	ret := _m.Called(ctx, requestId)

	if len(ret) == 0 {
		panic("no return value specified for GetStatus")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, requestId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, requestId)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, requestId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStatusService creates a new instance of StatusService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStatusService(t interface {
	mock.TestingT
	Cleanup(func())
}) *StatusService {
	mock := &StatusService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
