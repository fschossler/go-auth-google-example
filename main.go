package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

type JsonReturn struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		json := "{\"Homepage\"}"
		c.JSON(http.StatusOK, json)
	})

	router.GET("/login", func(c *gin.Context) {
		// Redirect to Google's consent page to ask for permission
		// for the scopes specified above.
		url := googleOauthConfig.AuthCodeURL("userstate")
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	router.GET("/callback", func(c *gin.Context) {

		state := c.Query("state")
		if state != "userstate" {
			log.Fatal("States not match.")
			return
		}

		code := c.Query("code")
		token, err := googleOauthConfig.Exchange(context.Background(), code)

		if err != nil {
			log.Fatal("The server could not exchange the token.")
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			log.Fatal(err)
		}

		userData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error in read response.")
		}

		var jsonReturn JsonReturn

		err = json.Unmarshal(userData, &jsonReturn)

		if err != nil {
			fmt.Println("error")
		}

		c.JSON(http.StatusOK, jsonReturn)

	})

	router.Run(":8080")
}
