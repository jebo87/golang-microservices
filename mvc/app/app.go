package app

import (
	"net/http"

	"github.com/jebo87/golang-microservices/mvc/controllers"
)

 

func StartApp(){
	http.HandleFunc("/users", controllers.GetUser)
	if err:= http.ListenAndServe("localhost:8080", nil); err != nil { 
		panic(err)
	} 
}