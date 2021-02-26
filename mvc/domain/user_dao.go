package domain

import (
	"fmt"

	"github.com/jebo87/golang-microservices/mvc/utils"
)

var (
	users = []*User{
		{
			ID:        1,
			FirstName: "Lionel",
			LastName:  "Messi",
			Email:     "lio@gmail.com",
		},
		{
			ID:        2,
			FirstName: "Cristiano",
			LastName:  "Ronaldo",
			Email:     "cr@gmail.com",
		},
	}

	UserDao userDaoInterface
)

func init() {
	UserDao = &userDao{}
}

type userDao struct{}

type userDaoInterface interface {
	GetUser(userID int64) (*User, *utils.ApplicationError)
}

func (u *userDao) GetUser(userID int64) (*User, *utils.ApplicationError) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}

	return nil, &utils.ApplicationError{
		Message: fmt.Sprintf("User %v not found", userID),
		Status:  404,
		Code:    "not_found",
	}
}
