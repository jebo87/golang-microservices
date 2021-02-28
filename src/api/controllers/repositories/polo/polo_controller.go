package polo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Polo(c *gin.Context) {
	c.String(http.StatusOK, "polo")
}
