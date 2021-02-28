package github_provider

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jebo87/golang-microservices/src/api/clients/restclient"
	"github.com/jebo87/golang-microservices/src/api/domain/github"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthorizationHeader(t *testing.T) {
	header := getAuthorizationHeader("abc123")
	assert.EqualValues(t, "token abc123", header)
}

func TestMain(m *testing.M) {
	//This is executed prior to any other test.
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestConstants(t *testing.T) {
	assert.EqualValues(t, "Authorization", headerAuthorization)
	assert.EqualValues(t, "token %s", headerAuthorizationFormat)
	assert.EqualValues(t, "https://api.github.com/user/repos", urlCreateRepo)
}

func TestCreateRepoErrorRestClient(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Err:        errors.New("Invalid restclient response"),
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Invalid restclient response", err.Message)

}

func TestCreateRepoInvalidResponseBody(t *testing.T) {
	restclient.FlushMockups()
	badBody, _ := os.Open("asdas")

	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: http.Response{
			StatusCode: http.StatusCreated,
			Body:       badBody,
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid response body", err.Message)

}

func TestCreateRepoInvalidJSONResponse(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id":"123}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid json response body", err.Message)

}

func TestCreateRepoServerError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message":"test"}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)

}

func TestCreateRepoCreatedButJSONParseError(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"asdasd`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Error when trying to unmarshal body succesful response", err.Message)

}

func TestCreateRepoCreatedSuccessful(t *testing.T) {
	restclient.FlushMockups()

	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: http.Response{
			StatusCode: http.StatusCreated,
			Body: ioutil.NopCloser(strings.NewReader(`{
				"id": 1296269,
				"name": "Hello-World",
				"full_name": "octocat/Hello-World",
				"owner": {
				  "login": "octocat",
				  "id": 1,				  
				  "url": "https://api.github.com/users/octocat",
				  "html_url": "https://github.com/octocat"				  
				},				
				"has_issues": true,
				"has_projects": true,
				"has_wiki": true,				
				"permissions": {
				  "admin": false,
				  "push": false,
				  "pull": true
				}				
			  }`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, "Hello-World", response.Name)
	assert.EqualValues(t, "octocat/Hello-World", response.FullName)

}
