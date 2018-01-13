package lobbymanager

/*HostManager is a singleton that will store a map of lobbies and room codes.
implemented based on sync-safe code by
Marcia Castilho http://marcio.io/2015/07/singleton-pattern-in-go/""
*/

import (
	"sync"
)

/*HostManager holds the map of room codes to lobby objects*/
type HostManager struct {
	LobbyMap map[string]*Lobby
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

func (hm *HostManager) AddLobby(l *Lobby) {
	if hm.LobbyMap == nil {
		hm.LobbyMap = make(map[string]*Lobby)
	}
	hm.LobbyMap[l.Code] = l

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
