package user

import "errors"

type User struct {
	ID    int
	Name  string
	Email string
}

var ErrNotFound = errors.New("user not found")
var ErrMailTaken = errors.New("email already taken")

type Repository interface {
	Save(u User) (User, error)
	FindByEmail(email string) (User, error)
	List() ([]User, error)
}
