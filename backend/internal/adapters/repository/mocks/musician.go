// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/hanoys/sigma-music/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// MusicianRepository is an autogenerated mock type for the IMusicianRepository type
type MusicianRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, musician
func (_m *MusicianRepository) Create(ctx context.Context, musician domain.Musician) (domain.Musician, error) {
	ret := _m.Called(ctx, musician)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 domain.Musician
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Musician) (domain.Musician, error)); ok {
		return rf(ctx, musician)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Musician) domain.Musician); ok {
		r0 = rf(ctx, musician)
	} else {
		r0 = ret.Get(0).(domain.Musician)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Musician) error); ok {
		r1 = rf(ctx, musician)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *MusicianRepository) GetByEmail(ctx context.Context, email string) (domain.Musician, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetByEmail")
	}

	var r0 domain.Musician
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.Musician, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Musician); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(domain.Musician)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: ctx, name
func (_m *MusicianRepository) GetByName(ctx context.Context, name string) (domain.Musician, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetByName")
	}

	var r0 domain.Musician
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.Musician, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Musician); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(domain.Musician)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMusicianRepository creates a new instance of MusicianRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMusicianRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MusicianRepository {
	mock := &MusicianRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}