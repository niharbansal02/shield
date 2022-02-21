// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/odpf/shield/model"
	mock "github.com/stretchr/testify/mock"
)

// PermissionCheckService is an autogenerated mock type for the PermissionCheckService type
type PermissionCheckService struct {
	mock.Mock
}

// CheckAuthz provides a mock function with given fields: ctx, resource, action
func (_m *PermissionCheckService) CheckAuthz(ctx context.Context, resource model.Resource, action model.Action) (bool, error) {
	ret := _m.Called(ctx, resource, action)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, model.Resource, model.Action) bool); ok {
		r0 = rf(ctx, resource, action)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.Resource, model.Action) error); ok {
		r1 = rf(ctx, resource, action)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
