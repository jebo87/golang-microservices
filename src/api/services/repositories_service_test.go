package services

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jebo87/golang-microservices/src/api/clients/restclient"
	"github.com/jebo87/golang-microservices/src/api/domain/repositories"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidInputName(t *testing.T) {
	request := repositories.CreateRepoRequest{}

	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repository name", err.Message())
}

func TestCreateRepoErrorFromGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication","documentation_url": "https://docs.github.com/rest/reference/repos#create-a-repository-for-the-authenticated-user"}`)),
		},
	})
	request := repositories.CreateRepoRequest{
		Name: "golang-example",
	}

	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Requires authentication", err.Message())
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())

}

func TestCreateRepoNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123,"name": "golang-example","description":"This is the description","owner":{"login":"jebo87"}}`)),
		},
	})

	request := repositories.CreateRepoRequest{
		Name:        "golang-example",
		Description: "This is the description",
	}

	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, "golang-example", result.Name)
	assert.EqualValues(t, 123, result.ID)
	assert.EqualValues(t, "jebo87", result.Owner)
}
