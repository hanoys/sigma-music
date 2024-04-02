package domain

type Comment struct {
	ID      int
	UserID  int
	TrackID int
	Stars   int
	Text    string
}
