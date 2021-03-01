package app

import (
	"github.com/jebo87/golang-microservices/src/api/controllers/polo"
	"github.com/jebo87/golang-microservices/src/api/controllers/repositories"
)

func mapURLs() {
	router.POST("/repository", repositories.CreateRepo)
	router.POST("/repositories", repositories.CreateRepos)
	router.GET("/marco", polo.Marco)
}
