package queue

import "fmt"

//Lobby is a list of players currently in a game
type Lobby struct {
	users []string
	code  string
}

/*GenerateLobby creates a new lobby with an empty list of users and unique lobby code

 */
func GenerateLobby() Lobby {
	l := Lobby{users: make([]string, 3), code: createLobbyCode()}

	return l
}

func (r *Lobby) addUser(name string) int {
	fmt.Println(name)
	return 0
}

/*
TODO: create randomized lobby code
*/
func createLobbyCode() string {
	return "ASDF"
}

/*
TODO: REMOVE THIS TESTING FUNCTION
*/

func PrintUsers() {
	fmt.Println("something happened")
}
