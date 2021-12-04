package opennoxcontrol

// Player represents information about a specific player.
type Player struct {
	Name  string `json:"name"`
	Class string `json:"class"`
}

// PlayerInfo lists information about all players on a server.
type PlayerInfo struct {
	Cur  int      `json:"cur"`
	Max  int      `json:"max"`
	List []Player `json:"list"`
}

// Info represents current game info.
type Info struct {
	Name       string     `json:"name"`
	Map        string     `json:"map"`
	Mode       string     `json:"mode"`
	Vers       string     `json:"vers"`
	PlayerInfo PlayerInfo `json:"players"`
}

// Map represent information about a Nox map.
type Map struct {
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	Flags      int    `json:"flags"`
	MinPlayers int    `json:"min_players"`
	MaxPlayers int    `json:"max_players"`
}

type Game interface {
	// GameInfo returns current game info, see Info.
	GameInfo() (Info, error)
	// ListMaps lists maps available on the server.
	ListMaps() ([]Map, error)
	// ChangeMap stops current game and loads a given map instead.
	ChangeMap(name string) error
	// Command executes a server command.
	Command(cmd string) error
}
