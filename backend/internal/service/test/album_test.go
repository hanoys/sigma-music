package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"github.com/stretchr/testify/mock"
	"testing"
)

var createAlbumReq = ports.CreateAlbumServiceReq{
	Name:        "Test Album Name",
	Description: "Album description",
}

func TestAlbumServiceCreate(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock func(repository *mocks.AlbumRepository)
		req            ports.CreateAlbumServiceReq
		expected       error
	}{
		{
			name: "create album success",
			req:  createAlbumReq,
			repositoryMock: func(repository *mocks.AlbumRepository) {
				repository.
					On("Create", context.Background(), mock.AnythingOfType("domain.Album")).
					Return(domain.Album{}, nil)
			},
			expected: nil,
		},
		{
			name: "create album failure",
			req:  createAlbumReq,
			repositoryMock: func(repository *mocks.AlbumRepository) {
				repository.
					On("Create", context.Background(), mock.AnythingOfType("domain.Album")).
					Return(domain.Album{}, ports.ErrAlbumDuplicate)
			},
			expected: ports.ErrAlbumDuplicate,
		},
	}

	for _, test := range tests {
		t.Logf("Test: %s", test.name)
		albumRepository := mocks.NewAlbumRepository(t)
		albumService := service.NewAlbumService(albumRepository)
		test.repositoryMock(albumRepository)

		_, err := albumService.Create(context.Background(), test.req)
		if !errors.Is(err, test.expected) {
			t.Errorf("got %v, want %v", err, test.expected)
		}
	}
}

func TestAlbumServicePublish(t *testing.T) {
	tests := []struct {
		name           string
		repositoryMock func(repository *mocks.AlbumRepository)
		id             uuid.UUID
		expected       error
	}{
		{
			name: "publish album success",
			id:   uuid.New(),
			repositoryMock: func(repository *mocks.AlbumRepository) {
				repository.
					On("Publish", context.Background(), mock.AnythingOfType("uuid.UUID")).
					Return(nil)
			},
			expected: nil,
		},
		{
			name: "publish album failure",
			id:   uuid.New(),
			repositoryMock: func(repository *mocks.AlbumRepository) {
				repository.
					On("Publish", context.Background(), mock.AnythingOfType("uuid.UUID")).
					Return(ports.ErrAlbumPublish)
			},
			expected: ports.ErrAlbumPublish,
		},
	}

	for _, test := range tests {
		t.Logf("Test: %s", test.name)
		albumRepository := mocks.NewAlbumRepository(t)
		albumService := service.NewAlbumService(albumRepository)
		test.repositoryMock(albumRepository)

		err := albumService.Publish(context.Background(), test.id)
		if !errors.Is(err, test.expected) {
			t.Errorf("got %v, want %v", err, test.expected)
		}
	}
}
