package entity

import (
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
)

type MongoMusician struct {
	ID          uuid.UUID `bson:"_id"`
	Name        string    `bson:"name"`
	Email       string    `bson:"email"`
	Password    string    `bson:"password"`
	Salt        string    `bson:"salt"`
	Country     string    `bson:"country"`
	Description string    `bson:"description"`
}

func (m *MongoMusician) ToDomain() domain.Musician {
	return domain.Musician{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		Password:    m.Password,
		Salt:        m.Salt,
		Country:     m.Country,
		Description: m.Description,
	}
}

func NewMongoMusician(musician domain.Musician) MongoMusician {
	return MongoMusician{
		ID:          musician.ID,
		Name:        musician.Name,
		Email:       musician.Email,
		Password:    musician.Password,
		Salt:        musician.Salt,
		Country:     musician.Country,
		Description: musician.Description,
	}
}
