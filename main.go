package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type MockResponse struct {
	Name  string
	Email string
}

func main() {

	g := gin.Default()
	requestCount := 0
	randomResponse := MockResponse{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	}
	g.GET("/test", func(c *gin.Context) {
		requestCount++
		if requestCount%3 == 0 {
			randomResponse = MockResponse{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
			}
		}

		body, err := json.Marshal(randomResponse)
		if err != nil {
			c.Error(err)
		}
		etag := fmt.Sprintf("%x", md5.Sum(body))
		c.Header("Cache-Control", "public, max-age=31536000")
		c.Header("ETag", etag)
		if match := c.GetHeader("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				c.Status(http.StatusNotModified)
				return
			}
		}
		c.JSON(200, randomResponse)
	})
	g.Run(":9000")
}
