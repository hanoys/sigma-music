package domain

import "time"

type Album struct {
	ID          int
	MusicianID  int
	Name        string
	Description string
	ReleaseDate time.Time
}
