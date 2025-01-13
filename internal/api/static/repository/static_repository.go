package repository

import "errors"

var users = []User{
	{ID: "1", Username: "john", Email: "john@example.com"},
	{ID: "2", Username: "jane", Email: "jane@example.com"},
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func FindUserByID(id string) (User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}
