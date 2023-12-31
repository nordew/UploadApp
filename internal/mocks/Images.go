// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"
	image "image"

	entity "github.com/nordew/UploadApp/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// Images is an autogenerated mock type for the Images type
type Images struct {
	mock.Mock
}

// GetAll provides a mock function with given fields: ctx, id
func (_m *Images) GetAll(ctx context.Context, id string) ([]entity.Image, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []entity.Image
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]entity.Image, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []entity.Image); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Image)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBySize provides a mock function with given fields: ctx, id, size
func (_m *Images) GetBySize(ctx context.Context, id string, size int) (*entity.Image, error) {
	ret := _m.Called(ctx, id, size)

	if len(ret) == 0 {
		panic("no return value specified for GetBySize")
	}

	var r0 *entity.Image
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) (*entity.Image, error)); ok {
		return rf(ctx, id, size)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int) *entity.Image); ok {
		r0 = rf(ctx, id, size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Image)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int) error); ok {
		r1 = rf(ctx, id, size)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Upload provides a mock function with given fields: ctx, _a1
func (_m *Images) Upload(ctx context.Context, _a1 image.Image) (string, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Upload")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, image.Image) (string, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, image.Image) string); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, image.Image) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewImages creates a new instance of Images. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewImages(t interface {
	mock.TestingT
	Cleanup(func())
}) *Images {
	mock := &Images{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
