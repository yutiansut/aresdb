// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	"unsafe"
)

// ListVectorParty is an autogenerated mock type for the ListVectorParty type
type ListVectorParty struct {
	mock.Mock
}

// GetListValue provides a mock function with given fields: row
func (_m *ListVectorParty) GetListValue(row int) (unsafe.Pointer, bool) {
	ret := _m.Called(row)

	var r0 unsafe.Pointer
	if rf, ok := ret.Get(0).(func(int) unsafe.Pointer); ok {
		r0 = rf(row)
	} else {
		r0 = ret.Get(0).(unsafe.Pointer)
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(int) bool); ok {
		r1 = rf(row)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// SetListValue provides a mock function with given fields: row, val, valid
func (_m *ListVectorParty) SetListValue(row int, val unsafe.Pointer, valid bool) {
	_m.Called(row, val, valid)
}
