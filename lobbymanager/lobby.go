package lobbymanager

import (
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"
	//"strconv"
)

//Lobby is a list of players currently in a game
type Lobby struct {
	Users        []User
	Code         string
	LastNotified string
	LastNotifiedTime string
}

type User struct {
	Name string
}

/*GenerateLobby creates a new lobby with an empty list of Users and unique lobby Code

 */
func GenerateLobby() *Lobby {
	l := Lobby{Users: make([]User, 0), Code: createLobbyCode()}

	return &l
}

/*
PrintLobby prints the information fields of the structure for debugging
*/
func (lob_Instance *Lobby) PrintLobby() {
	fmt.Println("Lobby Code" + lob_Instance.Code)

	for _, i := range lob_Instance.Users {
		fmt.Println("User " + i.Name + " ")
	}
}

func (lob_Instance *Lobby) AddUser(user User) {
	doAppend := true
	for _, ele := range lob_Instance.Users {
		if ele.Name == user.Name {
			doAppend = false
			break
		}
	}
	if doAppend {
		lob_Instance.Users = append(lob_Instance.Users, user)
	}
}

func (lob_Instance *Lobby) RemoveUser(name string) {
	for i, ele := range lob_Instance.Users {
		if ele.Name == name {
			lob_Instance.Users = append(lob_Instance.Users[:i], lob_Instance.Users[i+1:]...)
		}
	}
}

func (lob_Instance *Lobby) GetUserPosition(name string) int {
	for i, s := range lob_Instance.Users {
		if s.Name == name {
			return i + 1
		}
	}

	return -1

}

/*
TODO: create randomized lobby Code
*/
func createLobbyCode() string {

	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// The tabwriter here helps us generate aligned output.
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	defer w.Flush()

	code := ""

	rands := []int{r.Intn(25), r.Intn(25), r.Intn(25), r.Intn(25)}

	for _, i := range rands {
		code = code + string(toChar(i))
	}
	return code
}

func toChar(i int) rune {
	return rune('A' + i)
}

func (lob_Instance *Lobby) NotifyUser() {
	if len(lob_Instance.Users) > 0 {
		lob_Instance.LastNotified = lob_Instance.Users[0].Name

		//time format string must use Mon Jan 2 15:04:05 MST 2006
		lob_Instance.LastNotifiedTime = time.Now().Format("Mon 03:04 PM EST")
		//remove the first user from the lobby
		lob_Instance.RemoveUser(lob_Instance.Users[0].Name)
	}
}
