package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"fmt"
	"queueing-app/lobbymanager"

	"time"
)

func main() {
	lobbymanager.GetInstance()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/hostgame", func(c *gin.Context) {

		//If host already has a cookie
		//Serve the old lobby page
		val, err := c.Request.Cookie("roomcode")

		if err != nil {
			fmt.Println("error in cookie request")
		}

		if val != nil {
			fmt.Println("User has cookie " + val.Value)

			defer func() {
				fmt.Println("defering")
				l := createNewLobbyandPage()

				//Give that user a cookie for this lobby
				expiration := time.Now().Add(365 * 24 * time.Hour)
				cookie := http.Cookie{Name: "roomcode", Value: l.Code, Expires: expiration}
				http.SetCookie(c.Writer, &cookie)

				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}

			}()

			//serve webpage based on cookie
			//check webpage is in the map
			//	if lobbymanager.GetInstance().Contains(val.Value) {
			fmt.Println("Serving page " + val.Value)

			if _, err := os.Stat("host_page" + val.Value + ".tmpl.html"); os.IsNotExist(err) {
				fmt.Println("Panicking!")
				panic(createNewLobbyandPage())
			}

			c.HTML(http.StatusOK, "host_page"+val.Value+".tmpl.html", nil)

			return
			//	}
		}

		/*
			//otherwise, create a new lobby
			l := createNewLobbyandPage()

			//Give that user a cookie for this lobby
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "roomcode", Value: l.Code, Expires: expiration}
			http.SetCookie(c.Writer, &cookie)

			router.LoadHTMLGlob("templates/hosts/*.tmpl.html")
			c.HTML(http.StatusOK, "host_page"+l.Code+".tmpl.html", nil)
		*/
		return
	})

	router.GET("/joingame", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}

func createNewLobbyandPage() *lobbymanager.Lobby {
	l := lobbymanager.GenerateLobby()
	GetHostPage(*l)
	lobbymanager.GetInstance().AddLobby(l)
	return l
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetHostPage(l lobbymanager.Lobby) {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("host_page_template.tmpl").ParseFiles("templates/host_page_template.tmpl"))

	var webpage bytes.Buffer

	err := t.Execute(&webpage, l)
	if err != nil {
		log.Println("executing template:", err)
	}

	var OutputFile = "templates/hosts/host_page" + l.Code + ".tmpl.html"
	d1 := webpage.Bytes()
	err = ioutil.WriteFile(OutputFile, d1, 0644)
	check(err)
}
