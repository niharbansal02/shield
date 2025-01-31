// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	rule "github.com/odpf/shield/core/rule"
	mock "github.com/stretchr/testify/mock"
)

// RuleService is an autogenerated mock type for the RuleService type
type RuleService struct {
	mock.Mock
}

type RuleService_Expecter struct {
	mock *mock.Mock
}

func (_m *RuleService) EXPECT() *RuleService_Expecter {
	return &RuleService_Expecter{mock: &_m.Mock}
}

// GetAllConfigs provides a mock function with given fields: ctx
func (_m *RuleService) GetAllConfigs(ctx context.Context) ([]rule.Ruleset, error) {
	ret := _m.Called(ctx)

	var r0 []rule.Ruleset
	if rf, ok := ret.Get(0).(func(context.Context) []rule.Ruleset); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]rule.Ruleset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RuleService_GetAllConfigs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllConfigs'
type RuleService_GetAllConfigs_Call struct {
	*mock.Call
}

// GetAllConfigs is a helper method to define mock.On call
//  - ctx context.Context
func (_e *RuleService_Expecter) GetAllConfigs(ctx interface{}) *RuleService_GetAllConfigs_Call {
	return &RuleService_GetAllConfigs_Call{Call: _e.mock.On("GetAllConfigs", ctx)}
}

func (_c *RuleService_GetAllConfigs_Call) Run(run func(ctx context.Context)) *RuleService_GetAllConfigs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RuleService_GetAllConfigs_Call) Return(_a0 []rule.Ruleset, _a1 error) *RuleService_GetAllConfigs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewRuleService interface {
	mock.TestingT
	Cleanup(func())
}

// NewRuleService creates a new instance of RuleService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRuleService(t mockConstructorTestingTNewRuleService) *RuleService {
	mock := &RuleService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
