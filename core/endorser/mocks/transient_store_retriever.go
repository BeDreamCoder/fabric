// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	transientstore "github.com/hyperledger/fabric/core/transientstore"
	mock "github.com/stretchr/testify/mock"
)

// TransientStoreRetriever is an autogenerated mock type for the TransientStoreRetriever type
type TransientStoreRetriever struct {
	mock.Mock
}

// StoreForChannel provides a mock function with given fields: channel
func (_m *TransientStoreRetriever) StoreForChannel(channel string) *transientstore.Store {
	ret := _m.Called(channel)

	var r0 *transientstore.Store
	if rf, ok := ret.Get(0).(func(string) *transientstore.Store); ok {
		r0 = rf(channel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*transientstore.Store)
		}
	}

	return r0
}
