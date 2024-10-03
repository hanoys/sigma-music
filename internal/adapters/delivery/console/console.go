package console

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

type Option int

type Console struct {
	Handler  *Handler
	Routes   map[Option]func(*Console)
	UserID   uuid.UUID
	UserRole int
}

func NewConsole(h *Handler) *Console {
	c := &Console{Handler: h}
	c.InitRoutes()
	return c
}

const (
	exit Option = iota
	menu

	logIn
	logOut

	signUpUser
	getAllUsers
	getByIDUser
	getByNameUser

	signUpMusician
	getAllMusicians
	getByIDMusician
	getByNameMusician

	createTrack
	getAllTrack
	getByIDTrack
	deleteTrack
	getUserFavoritesTrack
	addToUserFavoritesTrack
	getByAlbumIDTrack
	getByMusicianIDTrack
	getOwnTracks

	postComment
	getOnTrackComment
	getUserComments

	getAllGenre

	createAlbum
	getAllAlbum
	getByMusicianIDAlbum
	getOwnAlbums
	getByIDAlbum
	publishAlbum

	getStatistic
	listen
)

func (c *Console) InitRoutes() {
	c.Routes = map[Option]func(*Console){
		exit: func(c *Console) { os.Exit(0) },
		menu: func(c *Console) { c.PrintMenu() },

		logIn:                   c.Handler.LogIn,
		logOut:                  c.Handler.LogOut,
		signUpUser:              c.Handler.SignUpUser,
		signUpMusician:          c.Handler.SignUpMusician,
		getAllUsers:             c.Handler.GetAllUsers,
		getByIDUser:             c.Handler.GetByIdUsers,
		getByNameUser:           c.Handler.GetByNameUser,
		getAllMusicians:         c.Handler.GetAllMusicians,
		getByIDMusician:         c.Handler.GetByIdMusician,
		getByNameMusician:       c.Handler.GetByNameMusician,
		createAlbum:             c.Handler.CreateAlbum,
		getAllAlbum:             c.Handler.GetAllAlbums,
		getByMusicianIDAlbum:    c.Handler.GetByMusicianIDAlbum,
		getOwnAlbums:            c.Handler.GetOwn,
		getByIDAlbum:            c.Handler.GetByIDAlbum,
		publishAlbum:            c.Handler.PublishAlbum,
		createTrack:             c.Handler.CreateTrack,
		getAllTrack:             c.Handler.GetAllTrack,
		getByIDTrack:            c.Handler.GetByIDTrack,
		deleteTrack:             c.Handler.DeleteTrack,
		getUserFavoritesTrack:   c.Handler.GetUserFavoritesTrack,
		addToUserFavoritesTrack: c.Handler.AddToUserFavoritesTrack,
		getByAlbumIDTrack:       c.Handler.GetByAlbumIDTrack,
		getByMusicianIDTrack:    c.Handler.GetByMusicianIDTrack,
		getAllGenre:             c.Handler.GetAllGenre,
		postComment:             c.Handler.PostComment,
		getOnTrackComment:       c.Handler.GetCommentsOnTrack,
		getUserComments:         c.Handler.GetUserComments,
		getStatistic:            c.Handler.GetStat,
		listen:                  c.Handler.Listen,
		getOwnTracks:            c.Handler.GetOwnTrack,
	}
}

func (c *Console) Start() error {
	time.Sleep(1 * time.Second)
	for {
		c.PrintMenu()

		var option Option
		fmt.Print("Choose menu option: ")
		_, err := fmt.Scanf("%d", &option)
		if err != nil {
			fmt.Println("Invalid menu option")
			continue
		}
		fmt.Println()

		handleFunc, ok := c.Routes[option]
		if !ok {
			fmt.Println("Invalid menu option")
			continue
		}
		handleFunc(c)
	}
}

func (c *Console) PrintMenu() {
	fmt.Println("-----------------------")
	fmt.Println("0. Exit")
	fmt.Println("1. Print menu")
	fmt.Println("2. Log In")
	fmt.Println("3. Log Out")
	fmt.Println("-----------------------")
	fmt.Println("4. Sign up user")
	fmt.Println("5. Get all users")
	fmt.Println("6. Get user by id")
	fmt.Println("7. Get user by name")
	fmt.Println("-----------------------")
	fmt.Println("8. Sign up musician")
	fmt.Println("9. Get all musicians")
	fmt.Println("10. Get musician by id")
	fmt.Println("11. Get musician by name")
	fmt.Println("-----------------------")
	fmt.Println("12. Create track (M)")
	fmt.Println("13. Get all tracks")
	fmt.Println("14. Get track by id")
	fmt.Println("15. Delete track (M)")
	fmt.Println("16. Get user favorites (U)")
	fmt.Println("17. Add track to user favorites (U)")
	fmt.Println("18. Get tracks by album id")
	fmt.Println("19. Get tracks by musician id")
	fmt.Println("20. Get own tracks")
	fmt.Println("-----------------------")
	fmt.Println("21. Post comment (U)")
	fmt.Println("22. Get comments on track (U)")
	fmt.Println("23. Get user comments (U)")
	fmt.Println("-----------------------")
	fmt.Println("24. Get all genres (U)")
	fmt.Println("-----------------------")
	fmt.Println("25. Create album (M)")
	fmt.Println("26. Get all albums")
	fmt.Println("27. Get albums by musician id")
	fmt.Println("28. Get own albums (M)")
	fmt.Println("29. Get album by id")
	fmt.Println("30. Publish album (M)")
	fmt.Println("-----------------------")
	fmt.Println("31. Get statistics (U)")
	fmt.Println("32. Listen Track (U)")
	fmt.Println("-----------------------")
}
