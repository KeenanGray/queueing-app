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
	Users []string
	Code string
}

/*GenerateLobby creates a new lobby with an empty list of Users and unique lobby Code

 */
func GenerateLobby() *Lobby {
	l := Lobby{Users: make([]string, 0), Code: createLobbyCode()}

	return &l
}

/*
PrintLobby prints the information fields of the structure for debugging
*/
func (lob_Instance *Lobby) PrintLobby() {
	fmt.Println(lob_Instance.Code)

	for _, i := range lob_Instance.Users{
		fmt.Println("User " + i + " ")	
	}
}

func (lob_Instance *Lobby) AddUser(name string) int {
	lob_Instance.Users = append(lob_Instance.Users, name)
	return len(lob_Instance.Users)
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

	code:=""

	rands := []int{r.Intn(25),r.Intn(25),r.Intn(25),r.Intn(25)}

	for _, i := range rands {
		code = code + string(toChar(i))
	}
	return code
}

func toChar(i int) rune {
	return rune('A' + i)
}
