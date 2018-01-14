package lobbymanager

/*HostManager is a singleton that will store a map of lobbies and room codes.
implemented based on sync-safe code by
Marcia Castilho http://marcio.io/2015/07/singleton-pattern-in-go/""
*/

import (
	"sync"
	"log"
	"fmt"
)

/*HostManager holds the map of room codes to lobby objects*/
type HostManager struct {
	LobbyMap map[string]*Lobby
	currentUID int
}

var instance *HostManager
var once sync.Once

/*GetInstance returns the single instance of this struct*/
func GetInstance() *HostManager {
	once.Do(func() {
		instance = &HostManager{}
	})
	return instance
}

func (hm *HostManager) Init(){
	if hm.LobbyMap == nil {
		hm.LobbyMap = make(map[string]*Lobby)
	}
	hm.currentUID=0;
}

func (hm *HostManager) PrintLobbies(){
	fmt.Println("Printing  Lobbies")
	for key, value := range hm.LobbyMap {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func (hm *HostManager) AddLobby(l *Lobby) {
	hm.LobbyMap[l.Code] = l
}

func (hm *HostManager) RemoveLobby(code string) {
	if hm.LobbyMap==nil {
		log.Fatal("No lobbies in being monitored by Host Manager")
	}
	delete(hm.LobbyMap,code)
}

func (hm *HostManager) AddUser(code string, user User) {
	if hm.LobbyMap == nil {
		log.Fatal("No game with this code exists")
	}
	hm.LobbyMap[code].AddUser(user)
	}

func (hm *HostManager) RemoveUser(code string, name string) {
	if hm.LobbyMap == nil {
		log.Fatal("No game with this code exists")
	}
	hm.LobbyMap[code].RemoveUser(name)
	}

func (hm *HostManager) Contains(code string) bool {
	if hm.LobbyMap == nil {
		return false
	}

	if (hm.LobbyMap[code] == nil) {
		return false
	} else {
		return true
	}

}

func (hm *HostManager) GetPositionInLobby(code string, name string) int {
	return hm.LobbyMap[code].GetUserPosition(name)	 
}

func (hm *HostManager) NotifyNextInQueue(code string){
	hm.LobbyMap[code].NotifyUser()
}