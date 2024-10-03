package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/adapters/repository/postgres/entity"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"github.com/hanoys/sigma-music/internal/util"
	"github.com/jmoiron/sqlx"
)

const (
	StatAddQuery             = "INSERT INTO users_history(id, user_id, track_id) VALUES ($1, $2, $3)"
	statGetMostListenedQuery = "select musician_id, $1 user_id, cnt " +
		"from (select a.id musician_id, count(*) cnt from (select m.id, uh.user_id from users_history uh " +
		"join tracks t on uh.track_id = t.id " +
		"join albums a on a.id = t.album_id " +
		"join album_musician am on am.album_id = a.id " +
		"join musicians m on m.id = am.musician_id " +
		"where uh.user_id=$1) as a " +
		"group by a.id " +
		"order by cnt DESC limit $2) t join musicians m on t.musician_id = m.id"
	statGetListenedGenresQuery = "select user_id, g.id genre_id, count(*) cnt from users_history uh " +
		"join tracks t on uh.track_id = t.id " +
		"join track_genre tg on t.id = tg.track_id " +
		"join genres g on g.id = tg.genre_id " +
		"where uh.user_id = $1 " +
		"group by user_id, g.id " +
		"order by cnt DESC"
)

type PostgresStatRepository struct {
	connection *sqlx.DB
}

func NewPostgresStatRepository(connection *sqlx.DB) *PostgresStatRepository {
	return &PostgresStatRepository{connection: connection}
}

func (sr *PostgresStatRepository) Add(ctx context.Context, recordID uuid.UUID, userID uuid.UUID, trackID uuid.UUID) error {
	_, err := sr.connection.ExecContext(ctx, StatAddQuery, recordID, userID, trackID)
	if err != nil {
		return util.WrapError(ports.ErrInternalStatRepo, err)
	}

	return nil
}

func (sr *PostgresStatRepository) GetMostListenedMusicians(ctx context.Context, userID uuid.UUID, maxCnt int) ([]domain.UserMusiciansStat, error) {
	var musiciansStat []entity.PgUserMusiciansStat
	err := sr.connection.SelectContext(ctx, &musiciansStat, statGetMostListenedQuery, userID, maxCnt)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalStatRepo, err)
	}

	domainMusiciansStat := make([]domain.UserMusiciansStat, len(musiciansStat))
	for i, stat := range musiciansStat {
		domainMusiciansStat[i] = stat.ToDomain()
	}

	return domainMusiciansStat, nil
}

func (sr *PostgresStatRepository) GetListenedGenres(ctx context.Context, userID uuid.UUID) ([]domain.UserGenresStat, error) {
	var genresStat []entity.PgUserGenresStat
	err := sr.connection.SelectContext(ctx, &genresStat, statGetListenedGenresQuery, userID)
	if err != nil {
		return nil, util.WrapError(ports.ErrInternalStatRepo, err)
	}

	domainGenresStat := make([]domain.UserGenresStat, len(genresStat))
	for i, stat := range genresStat {
		domainGenresStat[i] = stat.ToDomain()
	}

	return domainGenresStat, nil
}
