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

	//Cleanup the host folder
	os.RemoveAll("/tmp/")
	os.MkdirAll("/tmp/",777)
	
	//Spin up a null file in the host directory
	var webpage bytes.Buffer
	var OutputFile = "/tmp/nullHost.tmpl.html"
	d1 := webpage.Bytes()
	err := ioutil.WriteFile(OutputFile, d1, 0644)
	check(err)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	//router.Use(gin.ErrorLogger())

	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		router.LoadHTMLGlob("templates/*.tmpl.html")		
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/hostgame", func(c *gin.Context) {
		//If host already has a cookie
		//Serve the old lobby page
		
		//Load in host pages
		router.LoadHTMLGlob("tmp/*.tmpl.html")

		roomCode := getCookieValue(c)
		if roomCode != "" {
			fmt.Println("User has cookie " + roomCode)

			//check webpage is in the map
			//	if lobbymanager.GetInstance().Contains(val.Value) {

			//check that a webpage exists for this lobby
			code := getCookieValue(c)
			c.HTML(http.StatusOK, "host_page"+code+".tmpl.html", nil)

			if c.Errors != nil {
				fmt.Println("uhoh")
				fixMissingPage(c)
				c.Errors = nil
			} else {
				fmt.Println("Path FOUND")
				return
				//Serve the original page
			}

		} else {
			//User does not have a cookie, create a new page
			fmt.Println("User no cookie")
			fixMissingPage(c)
		}

		return
	})

	router.GET("/joingame", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/remake", func(c *gin.Context){
		fmt.Println("/tmp/host_page"+getCookieValue(c)+".tmpl.html")
		os.Remove("tmp/host_page"+getCookieValue(c)+".tmpl.html")
		assignCookie("",c)
		c.Redirect(302, "/hostgame")
	})

	router.Run(":" + port)
}

func assignCookie(code string, c *gin.Context) {
	//Give that user a new cookie for the replacement page
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "roomcode", Value: code, Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
}

func fixMissingPage(c *gin.Context) {
	l := createNewLobbyandPage()
	assignCookie(l.Code, c)

	fmt.Println("Created page " + l.Code)

	//serve the page
	c.HTML(http.StatusOK, "host_page"+l.Code+".tmpl.html", nil)
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

func getCookieValue(c *gin.Context) string {
	//Get value of cookie
	myCookie, err := c.Request.Cookie("roomcode")
	if err != nil {
		return ""
	}

	return myCookie.Value
}

func GetHostPage(l lobbymanager.Lobby) {
	// Create a new template and parse the letter into it.
	t := template.Must(template.New("host_page_template.tmpl").ParseFiles("templates/host_page_template.tmpl"))

	var webpage bytes.Buffer

	err := t.Execute(&webpage, l)
	if err != nil {
		log.Println("executing template:", err)
	}

	var OutputFile = "tmp/host_page" + l.Code + ".tmpl.html"
	d1 := webpage.Bytes()
	err = ioutil.WriteFile(OutputFile, d1, 0644)
	check(err)
}
