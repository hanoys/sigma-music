package entity

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/hanoys/sigma-music/internal/domain"
)

type PgMusician struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	Email       string      `db:"email"`
	Password    string      `db:"password"`
	Salt        string      `db:"salt"`
	Country     string      `db:"country"`
	Description string      `db:"description"`
	ImageURL    null.String `db:"image_url"`
}

func (m *PgMusician) ToDomain() domain.Musician {
	return domain.Musician{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		Password:    m.Password,
		Salt:        m.Salt,
		Country:     m.Country,
		Description: m.Description,
		ImageURL:    m.ImageURL,
	}
}

func NewPgMusician(musician domain.Musician) PgMusician {
	return PgMusician{
		ID:          musician.ID,
		Name:        musician.Name,
		Email:       musician.Email,
		Password:    musician.Password,
		Salt:        musician.Salt,
		Country:     musician.Country,
		Description: musician.Description,
		ImageURL:    musician.ImageURL,
	}
}
