package app

import (
	"github.com/jebo87/golang-microservices/mvc/controllers"
)

func mapURLs() {
	router.GET("/users/:user_id", controllers.GetUser)
}
