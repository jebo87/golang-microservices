package services

import (
	"log"
	"net/http"
	"testing"

	"github.com/jebo87/golang-microservices/mvc/domain"
	"github.com/jebo87/golang-microservices/mvc/utils"
	"github.com/stretchr/testify/assert"
)

var (
	UserDaoMock userDaoMock
	//Define a function that will be overwritten at each test case
	getUserFunction func(userID int64) (*domain.User, *utils.ApplicationError)
)

type userDaoMock struct{}

func init() {
	//This is the most important line since here we replace the UserDao for the Mock
	domain.UserDao = &userDaoMock{}
}

// This will bridge the actual GetUser implementation from the interface to the
// function we will be defining in our test cases.
func (m *userDaoMock) GetUser(userID int64) (*domain.User, *utils.ApplicationError) {
	return getUserFunction(userID)
}

func TestGetUserNotFoundInDB(t *testing.T) {
	//execute the actual function that will handle the mock data
	getUserFunction = func(userID int64) (*domain.User, *utils.ApplicationError) {
		log.Println("mocking the DB")
		return nil, &utils.ApplicationError{
			Status:  http.StatusNotFound,
			Message: "user 0 does not exist",
		}
	}
	user, err := GetUser(0)
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
}

func TestGetUserFoundInDB(t *testing.T) {
	log.Println("mocking the DB")
	getUserFunction = func(userID int64) (*domain.User, *utils.ApplicationError) {
		return &domain.User{
			ID:        999,
			FirstName: "Test Name",
			LastName:  "Test Lastname",
			Email:     "Test Email",
		}, nil
	}
	user, err := GetUser(1)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, "Test Name", user.FirstName)
	assert.EqualValues(t, "Test Lastname", user.LastName)
	assert.EqualValues(t, "Test Email", user.Email)
}
