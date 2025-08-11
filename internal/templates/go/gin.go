package golang

const (
	GinGoMod = `module {{.ProjectName}}

go 1.21

require github.com/gin-gonic/gin v1.9.1`

	GinMainGo = `package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World from {{.ProjectName}}!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Run(":3000")
}`
)
