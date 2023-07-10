// Code generated by mockery v2.30.16. DO NOT EDIT.

package ddb

import (
	context "context"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	mock "github.com/stretchr/testify/mock"
)

// MockDynamodbAPI is an autogenerated mock type for the DynamodbAPI type
type MockDynamodbAPI struct {
	mock.Mock
}

type MockDynamodbAPI_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDynamodbAPI) EXPECT() *MockDynamodbAPI_Expecter {
	return &MockDynamodbAPI_Expecter{mock: &_m.Mock}
}

// DeleteItem provides a mock function with given fields: ctx, params, optFns
func (_m *MockDynamodbAPI) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *dynamodb.DeleteItemOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) *dynamodb.DeleteItemOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.DeleteItemOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDynamodbAPI_DeleteItem_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteItem'
type MockDynamodbAPI_DeleteItem_Call struct {
	*mock.Call
}

// DeleteItem is a helper method to define mock.On call
//   - ctx context.Context
//   - params *dynamodb.DeleteItemInput
//   - optFns ...func(*dynamodb.Options)
func (_e *MockDynamodbAPI_Expecter) DeleteItem(ctx interface{}, params interface{}, optFns ...interface{}) *MockDynamodbAPI_DeleteItem_Call {
	return &MockDynamodbAPI_DeleteItem_Call{Call: _e.mock.On("DeleteItem",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *MockDynamodbAPI_DeleteItem_Call) Run(run func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options))) *MockDynamodbAPI_DeleteItem_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*dynamodb.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*dynamodb.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*dynamodb.DeleteItemInput), variadicArgs...)
	})
	return _c
}

func (_c *MockDynamodbAPI_DeleteItem_Call) Return(_a0 *dynamodb.DeleteItemOutput, _a1 error) *MockDynamodbAPI_DeleteItem_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDynamodbAPI_DeleteItem_Call) RunAndReturn(run func(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)) *MockDynamodbAPI_DeleteItem_Call {
	_c.Call.Return(run)
	return _c
}

// GetItem provides a mock function with given fields: ctx, params, optFns
func (_m *MockDynamodbAPI) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *dynamodb.GetItemOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) *dynamodb.GetItemOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.GetItemOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDynamodbAPI_GetItem_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetItem'
type MockDynamodbAPI_GetItem_Call struct {
	*mock.Call
}

// GetItem is a helper method to define mock.On call
//   - ctx context.Context
//   - params *dynamodb.GetItemInput
//   - optFns ...func(*dynamodb.Options)
func (_e *MockDynamodbAPI_Expecter) GetItem(ctx interface{}, params interface{}, optFns ...interface{}) *MockDynamodbAPI_GetItem_Call {
	return &MockDynamodbAPI_GetItem_Call{Call: _e.mock.On("GetItem",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *MockDynamodbAPI_GetItem_Call) Run(run func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options))) *MockDynamodbAPI_GetItem_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*dynamodb.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*dynamodb.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*dynamodb.GetItemInput), variadicArgs...)
	})
	return _c
}

func (_c *MockDynamodbAPI_GetItem_Call) Return(_a0 *dynamodb.GetItemOutput, _a1 error) *MockDynamodbAPI_GetItem_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDynamodbAPI_GetItem_Call) RunAndReturn(run func(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)) *MockDynamodbAPI_GetItem_Call {
	_c.Call.Return(run)
	return _c
}

// PutItem provides a mock function with given fields: ctx, params, optFns
func (_m *MockDynamodbAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *dynamodb.PutItemOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) *dynamodb.PutItemOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.PutItemOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDynamodbAPI_PutItem_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutItem'
type MockDynamodbAPI_PutItem_Call struct {
	*mock.Call
}

// PutItem is a helper method to define mock.On call
//   - ctx context.Context
//   - params *dynamodb.PutItemInput
//   - optFns ...func(*dynamodb.Options)
func (_e *MockDynamodbAPI_Expecter) PutItem(ctx interface{}, params interface{}, optFns ...interface{}) *MockDynamodbAPI_PutItem_Call {
	return &MockDynamodbAPI_PutItem_Call{Call: _e.mock.On("PutItem",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *MockDynamodbAPI_PutItem_Call) Run(run func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options))) *MockDynamodbAPI_PutItem_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*dynamodb.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*dynamodb.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*dynamodb.PutItemInput), variadicArgs...)
	})
	return _c
}

func (_c *MockDynamodbAPI_PutItem_Call) Return(_a0 *dynamodb.PutItemOutput, _a1 error) *MockDynamodbAPI_PutItem_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDynamodbAPI_PutItem_Call) RunAndReturn(run func(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)) *MockDynamodbAPI_PutItem_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDynamodbAPI creates a new instance of MockDynamodbAPI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDynamodbAPI(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDynamodbAPI {
	mock := &MockDynamodbAPI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
