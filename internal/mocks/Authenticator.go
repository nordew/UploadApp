// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	auth "github.com/nordew/UploadApp/pkg/auth"
	mock "github.com/stretchr/testify/mock"
)

// Authenticator is an autogenerated mock type for the Authenticator type
type Authenticator struct {
	mock.Mock
}

// GenerateToken provides a mock function with given fields: options
func (_m *Authenticator) GenerateToken(options *auth.GenerateTokenClaimsOptions) (string, error) {
	ret := _m.Called(options)

	if len(ret) == 0 {
		panic("no return value specified for GenerateToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*auth.GenerateTokenClaimsOptions) (string, error)); ok {
		return rf(options)
	}
	if rf, ok := ret.Get(0).(func(*auth.GenerateTokenClaimsOptions) string); ok {
		r0 = rf(options)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*auth.GenerateTokenClaimsOptions) error); ok {
		r1 = rf(options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseToken provides a mock function with given fields: accessToken
func (_m *Authenticator) ParseToken(accessToken string) (*auth.ParseTokenClaimsOutput, error) {
	ret := _m.Called(accessToken)

	if len(ret) == 0 {
		panic("no return value specified for ParseToken")
	}

	var r0 *auth.ParseTokenClaimsOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*auth.ParseTokenClaimsOutput, error)); ok {
		return rf(accessToken)
	}
	if rf, ok := ret.Get(0).(func(string) *auth.ParseTokenClaimsOutput); ok {
		r0 = rf(accessToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.ParseTokenClaimsOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accessToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuthenticator creates a new instance of Authenticator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthenticator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Authenticator {
	mock := &Authenticator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
