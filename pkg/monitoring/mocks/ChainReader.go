// Code generated by mockery v2.12.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	query "github.com/cosmos/cosmos-sdk/types/query"

	testing "testing"

	tx "github.com/cosmos/cosmos-sdk/types/tx"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ChainReader is an autogenerated mock type for the ChainReader type
type ChainReader struct {
	mock.Mock
}

// ContractState provides a mock function with given fields: ctx, contractAddress, queryMsg
func (_m *ChainReader) ContractState(ctx context.Context, contractAddress types.AccAddress, queryMsg []byte) ([]byte, error) {
	ret := _m.Called(ctx, contractAddress, queryMsg)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress, []byte) []byte); ok {
		r0 = rf(ctx, contractAddress, queryMsg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.AccAddress, []byte) error); ok {
		r1 = rf(ctx, contractAddress, queryMsg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TxsEvents provides a mock function with given fields: ctx, events, paginationParams
func (_m *ChainReader) TxsEvents(ctx context.Context, events []string, paginationParams *query.PageRequest) (*tx.GetTxsEventResponse, error) {
	ret := _m.Called(ctx, events, paginationParams)

	var r0 *tx.GetTxsEventResponse
	if rf, ok := ret.Get(0).(func(context.Context, []string, *query.PageRequest) *tx.GetTxsEventResponse); ok {
		r0 = rf(ctx, events, paginationParams)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tx.GetTxsEventResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string, *query.PageRequest) error); ok {
		r1 = rf(ctx, events, paginationParams)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewChainReader creates a new instance of ChainReader. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewChainReader(t testing.TB) *ChainReader {
	mock := &ChainReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
