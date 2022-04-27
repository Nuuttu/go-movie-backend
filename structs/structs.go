package mystructs

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

type User struct {
	Id           xid.ID  `json:"Id"`
	Username     string  `json:"Name" validate:"required,min=2,max=69"`
	Movielist    []Movie `json:"Movielist"`
	PasswordHash string  `json:"PasswordHash"`
}

type Credentials struct {
	Username string `json:"username" validate:"required,min=2,max=69"`
	Password string `json:"password" validate:"required,min=2,max=69"`
}

// uppercase json hmmhmm
type UserDetails struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type Session struct {
	Username   string
	Expiration time.Time
}

type Movie struct {
	Id       xid.ID    `json:"Id"`
	Name     string    `json:"Name" validate:"required,min=2,max=169"`
	Year     int       `json:"Year" validate:"numeric,gte=0,lte=2100"`
	Rating   int       `json:"Rating" validate:"numeric,gte=0,lte=10"`
	Review   string    `json:"Review" validate:"min=0,max=16900"`
	Watches  []Watch   `json:"Watches"`
	Comments []Comment `json:"Comments"`
}

type EditMovie struct {
	Name   string `json:"Name" validate:"required,min=2,max=169"`
	Year   int    `json:"Year" validate:"numeric,gte=0,lte=2100"`
	Rating int    `json:"Rating" validate:"numeric,gte=0,lte=10"`
	Review string `json:"Review" validate:"min=0,max=16900"`
}

type Watch struct {
	Id    xid.ID    `json:"Id"`
	Date  time.Time `json:"Date" validate:"datetime"`
	Place string    `json:"Place"`
	Note  string    `json:"Note"`
}

type Comment struct {
	Id            xid.ID    `json:"Id"`
	Owner         string    `json:"Owner"`
	Content       string    `json:"Content"`
	Creation_Time time.Time `json:"Creation_Time"`
}

var MovieList []Movie

// Set id to Pointed Movie struct
func (movie *Movie) SetId() {
	movie.Id = xid.New()
}

func (watch *Watch) SetId() {
	watch.Id = xid.New()
}

func (user *User) SetId() {
	user.Id = xid.New()
}

func (m *Movie) ModifyMovie(newMovie EditMovie) {
	fmt.Printf("new name: %s, new year: %d\n", (&newMovie).Name, &newMovie.Year)
	m.Name = newMovie.Name
	m.Year = newMovie.Year
	m.Rating = newMovie.Rating
	m.Review = newMovie.Review
}

func (s *Session) IsExpired() bool {
	return s.Expiration.Before(time.Now())
}
