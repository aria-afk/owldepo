// serve.go is the entry point to the API
package main

import (
	"back/pg"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	_ = pg.NewPG()
	r := gin.Default()
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "im healthy :)",
		})
	})
	r.Run()
}
