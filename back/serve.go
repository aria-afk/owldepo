// serve.go is the entry point to the API
package main

import (
	"back/pg"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db := pg.NewPG()
	r := gin.Default()
	r.UseRawPath = true
	r.UnescapePathValues = false
	r.GET("/internal/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "im healthy :)",
		})
	})

	r.GET("/api/v1/items/recent-by-name", func(c *gin.Context) {
		item := c.DefaultQuery("item", "")
		limit := c.DefaultQuery("limit", "20")

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			log.Printf("Error convertin limit to int \n%s", err)
			c.JSON(400, gin.H{"message": "Bad request, limit must be a valid integer"})
		}
		item = formatItemURI(item)

		log.Printf("\n ITEM : %s \n LIMIT: %d", item, limitInt)

		var jsonData []byte
		err = db.QueryRow("get-items-recent-entries", &jsonData, item, limitInt)
		if err != nil {
			log.Printf("Error performing query get-items-recent-entries\n%s", err)
			c.JSON(500, gin.H{"message": "internal server error"})
		}
		c.Data(200, "application/json", jsonData)
	})
	r.Run()
}

// This is really spaghetti and will be fixed eventually
// TLDR; instead of Scroll for Earring for INT 60%
// pass --> Scroll for Earring for INT 60
func formatItemURI(param string) string {
	if !strings.Contains(param, "Scroll") {
		return param
	}
	return param + "%"
}
