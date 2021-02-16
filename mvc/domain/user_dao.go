package domain

import (
	"fmt"

	"github.com/jebo87/golang-microservices/mvc/utils"
)

var(
	users = []*User {
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
)

func GetUser(userId int64) (*User, *utils.ApplicationError) {
	
	for _,user := range users {
		if user.ID == userId {
			return user, nil
		}
	}
	
	return nil, &utils.ApplicationError{
		Message: fmt.Sprintf("User %v not found", userId),
		Status:  404,
		Code:    "not_found",
	}
}