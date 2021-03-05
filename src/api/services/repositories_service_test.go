package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/jebo87/golang-microservices/src/api/clients/restclient"
	"github.com/jebo87/golang-microservices/src/api/domain/repositories"
	"github.com/jebo87/golang-microservices/src/api/utils/errors"
	"github.com/jebo87/golang-microservices/src/api/utils/mocks"
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
		Response: &http.Response{
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
		Response: &http.Response{
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

func TestRepoConcurrent(t *testing.T) {
	request := repositories.CreateRepoRequest{}
	output := make(chan repositories.CreateRepositoriesResult)
	defer close(output)

	service := reposService{}

	//we have to do it in a go rutine, otherwise will block
	go service.createRepoConcurrent(request, output)

	//blocks until we get an output.
	result := <-output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Error.Message())

}

func TestRepoConcurrentErrorGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication","documentation_url": "https://docs.github.com/rest/reference/repos#create-a-repository-for-the-authenticated-user"}`)),
		},
	})
	request := repositories.CreateRepoRequest{Name: "testing"}
	output := make(chan repositories.CreateRepositoriesResult)
	defer close(output)

	service := reposService{}

	//we have to do it in a go rutine, otherwise will block
	go service.createRepoConcurrent(request, output)

	//blocks until we get an output.
	result := <-output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusUnauthorized, result.Error.Status())
	assert.EqualValues(t, "Requires authentication", result.Error.Message())

}

func TestRepoConcurrentNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123,"name": "golang-example","description":"This is the description","owner":{"login":"jebo87"}}`)),
		},
	})

	request := repositories.CreateRepoRequest{
		Name:        "golang-example",
		Description: "This is the description",
	}

	output := make(chan repositories.CreateRepositoriesResult)
	defer close(output)
	service := reposService{}

	//we have to do it in a go rutine, otherwise will block
	go service.createRepoConcurrent(request, output)

	//blocks until we get an output.
	result := <-output
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Response)
	assert.EqualValues(t, "golang-example", result.Response.Name)
	assert.EqualValues(t, 123, result.Response.ID)
	assert.EqualValues(t, "jebo87", result.Response.Owner)

}

func TestHandleRepoResults(t *testing.T) {
	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)
	var wg sync.WaitGroup

	service := reposService{}

	//run in another gorutine sinec handleRepoResults
	//is blocking waiting on the input channel.
	go service.handleRepoResults(&wg, input, output)
	wg.Add(1)
	go func() {

		input <- repositories.CreateRepositoriesResult{
			Error: errors.NewBadRequestError("invalid repository name"),
		}
	}()
	wg.Wait()
	close(input)
	result := <-output
	assert.NotNil(t, result)
	assert.EqualValues(t, 0, result.StatusCode)
	assert.EqualValues(t, 1, len(result.Results))
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())

}

func TestCreateReposBadRequestsAllFail(t *testing.T) {

	//no need to add mockups since the requests wil fail
	//when validating the names

	requests := []repositories.CreateRepoRequest{
		{},
		{
			Name: "   ",
		},
	}
	result := RepositoryService.CreateRepos(requests)
	assert.NotNil(t, result)

	assert.EqualValues(t, http.StatusBadRequest, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())

	assert.EqualValues(t, http.StatusBadRequest, result.Results[1].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[1].Error.Message())

}

func TestCreateReposOneFail(t *testing.T) {

	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123,"name": "golang-example","description":"This is the description","owner":{"login":"jebo87"}}`)),
		},
	})
	requests := []repositories.CreateRepoRequest{
		{},
		{
			Name: "golang-example",
		},
	}
	result := RepositoryService.CreateRepos(requests)
	assert.NotNil(t, result)

	assert.EqualValues(t, http.StatusPartialContent, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	for _, resp := range result.Results {
		if resp.Error != nil {
			assert.EqualValues(t, http.StatusBadRequest, resp.Error.Status())
			assert.EqualValues(t, "invalid repository name", resp.Error.Message())
		} else {

			assert.EqualValues(t, 123, resp.Response.ID)
			assert.EqualValues(t, "golang-example", resp.Response.Name)
			assert.EqualValues(t, "jebo87", resp.Response.Owner)

		}
	}

}

func TestCreateReposAllSuccess(t *testing.T) {

	restclient.FlushMockups()

	//Use the mock client to support multiple requests
	restclient.StopMockups()
	restclient.Client = &mocks.MockClient{}

	mocks.DoFunc = func(*http.Request) (*http.Response, error) {
		r := ioutil.NopCloser(strings.NewReader(`{"id": 123,"name": "testing","description":"This is the description","owner":{"login":"jebo87"}}`))

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	requests := []repositories.CreateRepoRequest{
		{Name: "testing"},
		{Name: "testing"},
		{Name: "testing"},
	}
	result := RepositoryService.CreateRepos(requests)
	assert.NotNil(t, result)

	// assert.EqualValues(t, http.StatusCreated, result.StatusCode)
	// assert.EqualValues(t, 2, len(result.Results))
	a, _ := json.Marshal(result)
	log.Println(string(a))
	assert.EqualValues(t, "testing", result.Results[0].Response.Name)
	assert.EqualValues(t, 123, result.Results[0].Response.ID)
	assert.EqualValues(t, "jebo87", result.Results[0].Response.Owner)

	assert.EqualValues(t, "testing", result.Results[1].Response.Name)
	assert.EqualValues(t, 123, result.Results[1].Response.ID)
	assert.EqualValues(t, "jebo87", result.Results[1].Response.Owner)

}
