package app

import (
	"github.com/jebo87/golang-microservices/src/api/controllers/repositories"
	"github.com/jebo87/golang-microservices/src/api/controllers/repositories/polo"
)

func mapURLs() {
	router.POST("/repositories", repositories.CreateRepo)
	router.GET("/marco", polo.Polo)
}
