package test_utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMockedContext(request *http.Request, response http.ResponseWriter) *gin.Context {
	c, _ := gin.CreateTestContext(response)
	c.Request = request

	return c
}
