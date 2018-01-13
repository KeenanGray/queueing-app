package queue

import "fmt"

//Lobby is a list of players currently in a game
type Lobby struct {
	Users []string
	Code  string
}

/*GenerateLobby creates a new lobby with an empty list of Users and unique lobby Code

 */
func GenerateLobby() Lobby {
	l := Lobby{Users: make([]string, 3), Code: createLobbyCode()}

	return l
}

/*
PrintLobby prints the information fields of the structure for debugging
*/
func (lob_Instance *Lobby) PrintLobby() {
	fmt.Println(lob_Instance.Code)
}

func (lob_Instance *Lobby) addUser(name string) int {
	fmt.Println(name)
	return 0
}

/*
TODO: create randomized lobby Code
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
