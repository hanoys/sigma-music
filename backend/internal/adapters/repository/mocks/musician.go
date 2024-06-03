// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/hanoys/sigma-music/internal/domain"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
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

// GetAll provides a mock function with given fields: ctx
func (_m *MusicianRepository) GetAll(ctx context.Context) ([]domain.Musician, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.Musician
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.Musician, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Musician); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Musician)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByAlbumID provides a mock function with given fields: ctx, albumID
func (_m *MusicianRepository) GetByAlbumID(ctx context.Context, albumID uuid.UUID) (domain.Musician, error) {
	ret := _m.Called(ctx, albumID)

	if len(ret) == 0 {
		panic("no return value specified for GetByAlbumID")
	}

	var r0 domain.Musician
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (domain.Musician, error)); ok {
		return rf(ctx, albumID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) domain.Musician); ok {
		r0 = rf(ctx, albumID)
	} else {
		r0 = ret.Get(0).(domain.Musician)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, albumID)
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

// GetByID provides a mock function with given fields: ctx, musicianID
func (_m *MusicianRepository) GetByID(ctx context.Context, musicianID uuid.UUID) (domain.Musician, error) {
	ret := _m.Called(ctx, musicianID)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 domain.Musician
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (domain.Musician, error)); ok {
		return rf(ctx, musicianID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) domain.Musician); ok {
		r0 = rf(ctx, musicianID)
	} else {
		r0 = ret.Get(0).(domain.Musician)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, musicianID)
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

// GetByTrackID provides a mock function with given fields: ctx, trackID
func (_m *MusicianRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID) (domain.Musician, error) {
	ret := _m.Called(ctx, trackID)

	if len(ret) == 0 {
		panic("no return value specified for GetByTrackID")
	}

	var r0 domain.Musician
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (domain.Musician, error)); ok {
		return rf(ctx, trackID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) domain.Musician); ok {
		r0 = rf(ctx, trackID)
	} else {
		r0 = ret.Get(0).(domain.Musician)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, trackID)
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
