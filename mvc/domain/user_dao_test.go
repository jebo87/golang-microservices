package domain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserNoUserFound(t *testing.T){	
	// Initialization

	// Execution
	user, err := GetUser(0)

	// Validation
	assert.Nil(t, user, "We were not expecting to find a user with id 0")
	assert.NotNil(t, err, "We were expecting to receive an error when querying id 0")
	assert.EqualValues(t, http.StatusNotFound, err.Status, "We were expecting a 404 for id 0")
	assert.EqualValues(t, "not_found", err.Code, "We were expecting a 404 for id 0")
	assert.EqualValues(t, "User 0 not found", err.Message, "We were expecting a 404 for id 0")
}

func TestGetUserFound(t *testing.T) {
	user, err := GetUser(1)
	//Validate that the user is not nil, if it is nil, then display the error message
	assert.NotNil(t,user,"User should not be nil if the id was valid" )
	assert.Nil(t,err,"We should not receive and error if the id was valid")
	assert.EqualValues(t, 1, user.ID)
	assert.EqualValues(t,"Lionel", user.FirstName)
	assert.EqualValues(t,"Messi", user.LastName)
	assert.EqualValues(t,"lio@gmail.com", user.Email)


}


