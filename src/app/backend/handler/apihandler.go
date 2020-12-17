package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateHTTPAPIHandler() (http.Handler, error) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome GIN",
			},
		)
	})

	return r, nil
}
