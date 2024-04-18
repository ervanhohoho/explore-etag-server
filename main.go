package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {

	g := gin.Default()
	requestCount := 0
	currentUser := User{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	}
	g.GET("/", func(c *gin.Context) {
		requestCount++
		//if requestCount%3 == 0 {
		//	currentUser = User{
		//		Name:  gofakeit.Name(),
		//		Email: gofakeit.Email(),
		//	}
		//}

		body, err := json.Marshal(currentUser)
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
		c.JSON(200, currentUser)
	})

	g.PUT("/", func(c *gin.Context) {
		var newUser User
		requestPayloadBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if err := json.Unmarshal(requestPayloadBytes, &newUser); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		currentUserBytes, err := json.Marshal(currentUser)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		currentEtag := fmt.Sprintf("%x", md5.Sum(currentUserBytes))
		requestEtag := c.GetHeader("If-Match")
		if !strings.Contains(requestEtag, currentEtag) {
			c.Status(http.StatusPreconditionFailed)
			return
		}

		currentUser = newUser
		newEtag := fmt.Sprintf("%x", md5.Sum(requestPayloadBytes))

		c.Header("Cache-Control", "public, max-age=31536000")
		c.Header("ETag", newEtag)
		c.JSON(200, currentUser)
	})
	g.Run(":9000")
}
