package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	"queueing-app/lobbymanager"
	"strings"

	"time"
)

type User struct {
	Name string
	Pos  int
}

func main() {
	lobbymanager.GetInstance().Init()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	//router.Use(gin.ErrorLogger())
	router.LoadHTMLGlob("templates/*.tmpl.html")

	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		router.LoadHTMLGlob("templates/*.tmpl.html")
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/hostgame", func(c *gin.Context) {
		//Get the hosts cookie
		roomCode := getCookieValue("HostInfo", c)
		if roomCode != "" {
			//If host already has a cookie
			//Serve the lobby page with is cookie
			code := getCookieValue("HostInfo", c)
			l := lobbymanager.GetInstance().LobbyMap[code]
			if l != nil {
				c.HTML(http.StatusOK, "host_page.tmpl.html", l)
				return
			} else {
			}
		}
		//User does not have a cookie, create a new page
		createNewLobbyandAssign(c)

		c.Redirect(302, "/hostgame")
		return
	})

	router.POST("/hostgame", func(c *gin.Context) {
		//remove current host from map
		removeCodeFromMap(getCookieValue("HostInfo", c))
		//reassign the host cookie to null
		assignHostCookie("", c)
		//refresh the hostgame page
		c.Redirect(302, "/hostgame")
	})

	router.POST("/notify", func(c *gin.Context){
		RoomCode:=getCookieValue("HostInfo",c)
		lobbymanager.GetInstance().NotifyNextInQueue(RoomCode)

		c.Redirect(302, "/hostgame")
	})

	router.GET("/joingame", func(c *gin.Context) {
		c.HTML(http.StatusOK, "joingame.tmpl.html", nil)
	})	
	router.POST("/joingame", func(c *gin.Context) {
		joincode := strings.ToUpper(c.Request.FormValue("code"))
		joinname := strings.ToUpper(c.Request.FormValue("name"))

		if lobbymanager.GetInstance().Contains(joincode) {
			//else create a new user
			lobbymanager.GetInstance().AddUser(joincode, lobbymanager.User{Name: joinname})
			assignUserCookie(joinname+","+joincode, c)
			c.Redirect(302, "/ingame")		
		} else {
			c.Redirect(302, "/joingame")
			return
		}
	})

	router.GET("/ingame", func(c *gin.Context) {
		//get user info
		UserInfo := getCookieValue("UserInfo", c)
		UserName := strings.Split(UserInfo, ",")[0]
		RoomCode := strings.Split(UserInfo, ",")[1]

		//check that lobby exists in manager
		if lobbymanager.GetInstance().Contains(RoomCode) {
			c.HTML(http.StatusOK, "client_page.tmpl.html", User{UserName, lobbymanager.GetInstance().GetPositionInLobby(RoomCode, UserName)})
		} else {
			c.Redirect(302, "/joingame")
		}

		return
	})

	router.POST("/leave", func(c *gin.Context) {
		UserInfo := getCookieValue("UserInfo", c)
		if UserInfo != "" {
			UserName := strings.Split(UserInfo, ",")[0]
			RoomCode := strings.Split(UserInfo, ",")[1]

			//find the lobby
			lobbymanager.GetInstance().RemoveUser(RoomCode, UserName)
		}
		assignUserCookie("", c)
		c.Redirect(302, "/")
	})

	router.GET("/reset", func(c *gin.Context){
		//resets cookies and returns to home page
		assignHostCookie("",c)
		assignUserCookie("",c)
		c.Redirect(302, "/")
	})

	router.Run(":" + port)
}

func assignHostCookie(code string, c *gin.Context) {
	//Give that user a new cookie for the replacement page
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "HostInfo", Value: code, Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
}

func assignUserCookie(UserInfo string, c *gin.Context) {
	//Give that user a new cookie for the replacement page
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "UserInfo", Value: UserInfo, Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
}

func getCookieValue(cookieName string, c *gin.Context) string {
	//Get value of cookie
	myCookie, err := c.Request.Cookie(cookieName)
	if err != nil {
		return ""
	}

	return myCookie.Value
}

func createNewLobbyandAssign(c *gin.Context) *lobbymanager.Lobby {
	l := lobbymanager.GenerateLobby()
	lobbymanager.GetInstance().AddLobby(l)

	assignHostCookie(l.Code, c)

	return l
}

func removeCodeFromMap(code string) {
	lobbymanager.GetInstance().RemoveLobby(code)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
