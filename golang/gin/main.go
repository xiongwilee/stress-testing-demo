package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(time.Millisecond * 50)
		c.String(200, strings.Repeat("haha", 1024))
	})
	router.Run(":3002")
}
