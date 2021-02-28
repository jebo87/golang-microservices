package github

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepoRequest(t *testing.T) {
	request := CreateRepoRequest{
		Name:        "golang introduction",
		Description: "A golang repository",
		Homepage:    "http://github.com",
		Private:     true,
		HasIssues:   false,
		HasProjects: true,
		HasWiki:     false,
	}
	bytes, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var target CreateRepoRequest
	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	assert.EqualValues(t, request.Name, target.Name)
	assert.EqualValues(t, request.Description, target.Description)
	assert.EqualValues(t, request.Homepage, target.Homepage)
	assert.EqualValues(t, request.Private, target.Private)
	assert.EqualValues(t, request.HasIssues, target.HasIssues)
	assert.EqualValues(t, request.HasProjects, target.HasProjects)
	assert.EqualValues(t, request.HasWiki, target.HasWiki)
}

func TestCreateRepoResponse(t *testing.T) {
	response := CreateRepoResponse{
		ID:       123,
		Name:     "Test response",
		FullName: "This is the full name",
		Owner: RepoOwner{
			ID:      22522,
			Login:   "jebo87",
			Url:     "http://test.com",
			HtmlUrl: "http://another.com",
		},
		Permissions: RepoPermissions{
			IsAdmin: true,
			HasPull: true,
			HasPush: true,
		},
	}

	bytes, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	fmt.Println(string(bytes))

	var target CreateRepoResponse
	err = json.Unmarshal(bytes, &target)
	assert.Nil(t, err)
	fmt.Println(err)

	assert.EqualValues(t, response.FullName, target.FullName)
	assert.EqualValues(t, response.Owner.ID, target.Owner.ID)
	assert.EqualValues(t, response.Permissions.IsAdmin, target.Permissions.IsAdmin)

}
