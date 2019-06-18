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
		c.String(200, strings.Repeat("s", 4096))
	})
	router.Run(":3001")
}
