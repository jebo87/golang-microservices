package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jebo87/golang-microservices/mvc/services"
	"github.com/jebo87/golang-microservices/mvc/utils"
)

func GetUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)

	if err != nil {
		appError := &utils.ApplicationError{
			Message: "Error parsing the request",
			Status:  http.StatusBadRequest,
			Code:    "bad_request",
		}
		utils.Respond(c, appError.Status, appError)
		return
	}

	user, appErr := services.GetUser(userID)

	if appErr != nil {
		utils.Respond(c, appErr.Status, appErr)
		return
	}
	utils.Respond(c, http.StatusOK, user)
}
