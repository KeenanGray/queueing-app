package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"queueing-app/queue"
)

func main() {
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
		l := queue.GenerateLobby()

		fmt.Println("Hey")
		l.PrintLobby()

		GetHostPage(l)

		c.HTML(http.StatusOK, "host_page.tmpl.html", nil)
	})

	router.GET("/joingame", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.Run(":" + port)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetHostPage(l queue.Lobby) {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("host_page_template.tmpl").ParseFiles("templates/host_page_template.tmpl"))

	var webpage bytes.Buffer

	// Execute the template for each recipient.
	/*
		for _, r := range Projects {
			err := t.Execute(&webpage, r)
			if err != nil {
				log.Println("executing template:", err)
			}
		}
	*/

	err := t.Execute(&webpage, l)
	if err != nil {
		log.Println("executing template:", err)
	}

	var OutputFile = "templates/host_page.tmpl.html"
	d1 := webpage.Bytes()
	err = ioutil.WriteFile(OutputFile, d1, 0644)
	check(err)
}
