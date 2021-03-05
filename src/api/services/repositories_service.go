package services

import (
	"log"
	"net/http"
	"sync"

	"github.com/jebo87/golang-microservices/src/api/config"
	"github.com/jebo87/golang-microservices/src/api/domain/github"
	"github.com/jebo87/golang-microservices/src/api/domain/github/providers/github_provider"
	"github.com/jebo87/golang-microservices/src/api/domain/repositories"
	"github.com/jebo87/golang-microservices/src/api/utils/errors"
)

type reposService struct{}

type repoServiceInterface interface {
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(request []repositories.CreateRepoRequest) repositories.CreateReposResponse
}

var (
	RepositoryService repoServiceInterface
)

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	log.Println(config.GetGithubAccessToken())
	if err := input.Validate(); err != nil {
		return nil, err
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Private:     false,
		Description: input.Description,
	}

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	result := repositories.CreateRepoResponse{
		ID:    response.ID,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}
	return &result, nil

}

func (s *reposService) CreateRepos(requests []repositories.CreateRepoRequest) repositories.CreateReposResponse {
	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup
	go s.handleRepoResults(&wg, input, output)

	for _, current := range requests {
		wg.Add(1)
		go s.createRepoConcurrent(current, input)
	}
	wg.Wait()
	close(input)
	result := <-output

	successfullCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successfullCreations++
		}
	}

	if successfullCreations == 0 {
		result.StatusCode = result.Results[0].Error.Status()
	} else if successfullCreations == len(requests) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}

	return result

}

func (s *reposService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRepositoriesResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for incomingEvent := range input {
		repoResult := repositories.CreateRepositoriesResult{
			Response: incomingEvent.Response,
			Error:    incomingEvent.Error,
		}
		results.Results = append(results.Results, repoResult)
		wg.Done()
	}
	output <- results
}

func (s *reposService) createRepoConcurrent(input repositories.CreateRepoRequest, output chan repositories.CreateRepositoriesResult) {
	if err := input.Validate(); err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	result, err := s.CreateRepo(input)

	if err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	output <- repositories.CreateRepositoriesResult{Response: result}
}
