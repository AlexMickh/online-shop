package models

type User struct {
	ID              string `redis:"id"`
	Login           string
	Email           string
	Password        string
	Role            string `redis:"role"`
	IsEmailVerified bool
}

type Category struct {
	ID   string `redis:"-"`
	Name string `redis:"name"`
}
