package test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/hash"
	"github.com/hanoys/sigma-music/internal/adapters/repository/mocks"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/service"
	"go.uber.org/zap"
	"testing"
)

var userID = uuid.New()

var musicians = []domain.Musician{
	{
		ID:          uuid.New(),
		Name:        "mus1",
		Email:       "email1",
		Password:    "pas1",
		Country:     "con1",
		Description: "descr1",
	},
	{
		ID:          uuid.New(),
		Name:        "mus2",
		Email:       "email2",
		Password:    "pas2",
		Country:     "con2",
		Description: "descr2",
	},
	{
		ID:          uuid.New(),
		Name:        "mus2",
		Email:       "email2",
		Password:    "pas2",
		Country:     "con2",
		Description: "descr2",
	},
}

var genres = []domain.Genre{
	{
		ID:   uuid.New(),
		Name: "rock",
	},
	{
		ID:   uuid.New(),
		Name: "rap",
	},
}

var users = []domain.User{
	{
		ID:       uuid.New(),
		Name:     "user1",
		Email:    "usermail",
		Phone:    "+793293293293",
		Password: "userpassword",
		Country:  "usercountry",
	},
}

var listenedMusicians = []domain.UserMusiciansStat{
	{
		MusicianID:  musicians[0].ID,
		UserID:      userID,
		ListenCount: 10,
	},
	{
		MusicianID:  musicians[1].ID,
		UserID:      userID,
		ListenCount: 3,
	},
	{
		MusicianID:  musicians[2].ID,
		UserID:      userID,
		ListenCount: 6,
	},
}

var listenedGenres = []domain.UserGenresStat{
	{
		GenreID:     genres[0].ID,
		UserID:      userID,
		ListenCount: 10,
	},
	{
		GenreID:     genres[1].ID,
		UserID:      userID,
		ListenCount: 13,
	},
}

func TestStatServiceFromReport(t *testing.T) {
	tests := []struct {
		name             string
		statRepoMock     func(repository *mocks.StatRepository)
		musicianRepoMock func(repository *mocks.MusicianRepository)
		genreRepoMock    func(repository *mocks.GenreRepository)
		expected         struct {
			err               error
			listenCnt         int64
			genresPercentages []int64
		}
	}{
		{
			name: "test 1",
			statRepoMock: func(repository *mocks.StatRepository) {
				repository.
					On("GetMostListenedMusicians", context.Background(), userID, 3).
					Return(listenedMusicians, nil).
					On("GetListenedGenres", context.Background(), userID).
					Return(listenedGenres, nil)
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {
				repository.
					On("GetByID", context.Background(), listenedMusicians[0].MusicianID).
					Return(musicians[0], nil).
					On("GetByID", context.Background(), listenedMusicians[1].MusicianID).
					Return(musicians[1], nil).
					On("GetByID", context.Background(), listenedMusicians[2].MusicianID).
					Return(musicians[2], nil)
			},
			genreRepoMock: func(repository *mocks.GenreRepository) {
				repository.
					On("GetByID", context.Background(), listenedGenres[0].GenreID).
					Return(genres[0], nil).
					On("GetByID", context.Background(), listenedGenres[1].GenreID).
					Return(genres[1], nil)
			},
			expected: struct {
				err               error
				listenCnt         int64
				genresPercentages []int64
			}{err: nil, listenCnt: 19, genresPercentages: []int64{43, 57}},
		},
		{
			name: "test 1",
			statRepoMock: func(repository *mocks.StatRepository) {
				repository.
					On("GetMostListenedMusicians", context.Background(), userID, 3).
					Return(listenedMusicians, ports.ErrInternalStatRepo)
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {
			},
			genreRepoMock: func(repository *mocks.GenreRepository) {
			},
			expected: struct {
				err               error
				listenCnt         int64
				genresPercentages []int64
			}{err: ports.ErrInternalStatRepo, listenCnt: 0, genresPercentages: []int64{}},
		},
		{
			name: "test 1",
			statRepoMock: func(repository *mocks.StatRepository) {
				repository.
					On("GetMostListenedMusicians", context.Background(), userID, 3).
					Return(listenedMusicians, nil).
					On("GetListenedGenres", context.Background(), userID).
					Return(listenedGenres, ports.ErrInternalStatRepo)
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {
			},
			genreRepoMock: func(repository *mocks.GenreRepository) {
			},
			expected: struct {
				err               error
				listenCnt         int64
				genresPercentages []int64
			}{err: ports.ErrInternalStatRepo, listenCnt: 0, genresPercentages: []int64{}},
		},
		{
			name: "test 1",
			statRepoMock: func(repository *mocks.StatRepository) {
				repository.
					On("GetMostListenedMusicians", context.Background(), userID, 3).
					Return(listenedMusicians, nil).
					On("GetListenedGenres", context.Background(), userID).
					Return(listenedGenres, nil)
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {
				repository.
					On("GetByID", context.Background(), listenedMusicians[0].MusicianID).
					Return(musicians[0], ports.ErrMusicianIDNotFound)
			},
			genreRepoMock: func(repository *mocks.GenreRepository) {
			},
			expected: struct {
				err               error
				listenCnt         int64
				genresPercentages []int64
			}{err: ports.ErrMusicianIDNotFound, listenCnt: 0, genresPercentages: []int64{}},
		},
		{
			name: "test 1",
			statRepoMock: func(repository *mocks.StatRepository) {
				repository.
					On("GetMostListenedMusicians", context.Background(), userID, 3).
					Return(listenedMusicians, nil).
					On("GetListenedGenres", context.Background(), userID).
					Return(listenedGenres, nil)
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {
				repository.
					On("GetByID", context.Background(), listenedMusicians[0].MusicianID).
					Return(musicians[0], nil).
					On("GetByID", context.Background(), listenedMusicians[1].MusicianID).
					Return(musicians[1], nil).
					On("GetByID", context.Background(), listenedMusicians[2].MusicianID).
					Return(musicians[2], nil)
			},
			genreRepoMock: func(repository *mocks.GenreRepository) {
				repository.
					On("GetByID", context.Background(), listenedGenres[0].GenreID).
					Return(genres[0], ports.ErrGenreIDNotFound)
			},
			expected: struct {
				err               error
				listenCnt         int64
				genresPercentages []int64
			}{err: ports.ErrGenreIDNotFound, listenCnt: 0, genresPercentages: []int64{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			musicianRepo := mocks.NewMusicianRepository(t)
			musicianService := service.NewMusicianService(musicianRepo, hash.NewHashPasswordProvider(), logger)
			test.musicianRepoMock(musicianRepo)

			genreRepo := mocks.NewGenreRepository(t)
			genreService := service.NewGenreService(genreRepo, logger)
			test.genreRepoMock(genreRepo)

			statRepo := mocks.NewStatRepository(t)
			statService := service.NewStatService(statRepo, genreService, musicianService, logger)
			test.statRepoMock(statRepo)

			res, err := statService.FormReport(context.Background(), userID)
			if !errors.Is(err, test.expected.err) {
				t.Errorf("got: %v, expected: %v", err, test.expected.err)
			}

			if err != nil {
				for i, percentage := range test.expected.genresPercentages {
					if res.ListenedGenres[i].ListenPercentage != percentage {
						t.Errorf("got: %v, expected: %v", res.ListenedGenres[i].ListenPercentage, percentage)
					}
				}
			}
		})
	}
}

var trackID = uuid.New()

func TestStatServiceAdd(t *testing.T) {
	tests := []struct {
		name             string
		statRepoMock     func(repository *mocks.StatRepository)
		musicianRepoMock func(repository *mocks.MusicianRepository)
		genreRepoMock    func(repository *mocks.GenreRepository)
		expected         error
	}{
		{
			name: "test1",
			statRepoMock: func(repository *mocks.StatRepository) {
				repository.
					On("Add", context.Background(), userID, trackID).
					Return(nil)
			},
			musicianRepoMock: func(repository *mocks.MusicianRepository) {
			},
			genreRepoMock: func(repository *mocks.GenreRepository) {
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			musicianRepo := mocks.NewMusicianRepository(t)
			musicianService := service.NewMusicianService(musicianRepo, hash.NewHashPasswordProvider(), logger)
			test.musicianRepoMock(musicianRepo)

			genreRepo := mocks.NewGenreRepository(t)
			genreService := service.NewGenreService(genreRepo, logger)
			test.genreRepoMock(genreRepo)

			statRepo := mocks.NewStatRepository(t)
			statService := service.NewStatService(statRepo, genreService, musicianService, logger)
			test.statRepoMock(statRepo)

			err := statService.Add(context.Background(), userID, trackID)
			if !errors.Is(err, test.expected) {
				t.Errorf("got: %v, expected: %v", err, test.expected)
			}
		})
	}
}
