// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	fcdclient "github.com/smartcontractkit/chainlink-terra/pkg/monitoring/fcdclient"

	testing "testing"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// GetBlockAtHeight provides a mock function with given fields: _a0, _a1
func (_m *Client) GetBlockAtHeight(_a0 context.Context, _a1 uint64) (fcdclient.Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 fcdclient.Response
	if rf, ok := ret.Get(0).(func(context.Context, uint64) fcdclient.Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(fcdclient.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTxList provides a mock function with given fields: _a0, _a1
func (_m *Client) GetTxList(_a0 context.Context, _a1 fcdclient.GetTxListParams) (fcdclient.Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 fcdclient.Response
	if rf, ok := ret.Get(0).(func(context.Context, fcdclient.GetTxListParams) fcdclient.Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(fcdclient.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, fcdclient.GetTxListParams) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewClient creates a new instance of Client. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewClient(t testing.TB) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
