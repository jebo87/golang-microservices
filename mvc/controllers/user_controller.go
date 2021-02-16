package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jebo87/golang-microservices/mvc/services"
	"github.com/jebo87/golang-microservices/mvc/utils"
)

func GetUser( res http.ResponseWriter, req *http.Request,){
  userID, err := strconv.ParseInt(req.URL.Query().Get("user_id"),10, 64)

  if err != nil {
	  appError:=&utils.ApplicationError{
	  	Message: "Error parsing the request",
	  	Status:  400,
	  	Code:    "bad_request",
	  }
	res.WriteHeader(appError.Status)
	jsonValue, _ := json.MarshalIndent(appError,"","\t")
	res.Write(jsonValue)
	  return 
  }

  user, appErr := services.GetUser(userID)
  if appErr != nil {
	res.WriteHeader(appErr.Status)
	jsonValue,_ := json.MarshalIndent(appErr,"", "\t")
	res.Write(jsonValue)
	  return 
  }

  jsonValue,_ := json.MarshalIndent(user,"", "\t")
  res.Header().Add("Content-Type", "application/json")
  res.WriteHeader(http.StatusOK)  
  res.Write(jsonValue)

}