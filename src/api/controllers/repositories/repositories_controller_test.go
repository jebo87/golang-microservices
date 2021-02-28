package repositories

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jebo87/golang-microservices/src/api/clients/restclient"
	"github.com/jebo87/golang-microservices/src/api/domain/repositories"
	"github.com/jebo87/golang-microservices/src/api/utils/errors"
	"github.com/jebo87/golang-microservices/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidJsonRequest(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/repositories", strings.NewReader(``))
	c := test_utils.GetMockedContext(request, response)

	CreateRepo(c)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)

	//We need to create an error from the response body
	apiErr, err := errors.NewApiErrFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "invalid json body", apiErr.Message())
}

func TestCreateRepoErrorGithub(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication","documentation_url": "https://docs.github.com/rest/reference/repos#create-a-repository-for-the-authenticated-user"}`)),
		},
	})
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/repositories", strings.NewReader(`{"name":"test"}`))
	c := test_utils.GetMockedContext(request, response)

	CreateRepo(c)

	assert.EqualValues(t, http.StatusUnauthorized, response.Code) //this comes from the mockup rest client

	//We need to create an error from the response body
	apiErr, err := errors.NewApiErrFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())    //this comes from the mockup rest client
	assert.EqualValues(t, "Requires authentication", apiErr.Message()) //this comes from the mockup rest client
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
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/repositories", strings.NewReader(`{"name":"test"}`))
	c := test_utils.GetMockedContext(request, response)

	CreateRepo(c)

	assert.EqualValues(t, http.StatusCreated, response.Code) //this comes from the mockup rest client

	//We need to create an error from the response body
	result := &repositories.CreateRepoResponse{}
	err := json.Unmarshal(response.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.EqualValues(t, 123, result.ID)                //this comes from the mockup rest client
	assert.EqualValues(t, "golang-example", result.Name) //this comes from the mockup rest client
	assert.EqualValues(t, "jebo87", result.Owner)        //this comes from the mockup rest client
}
