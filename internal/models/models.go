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

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float32
	ImageUrl    string
	Categories  []Category
}

type ProductCard struct {
	ID       string
	Name     string
	Price    float32
	ImageUrl string
}

type Cart struct {
	ID       string
	UserId   string
	Price    float32
	Products []ProductCard
}
