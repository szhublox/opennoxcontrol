package main

type Player struct {
	Name  string `json:"name"`
	Class string `json:"class"`
}

type PlayerInfo struct {
	Cur  int      `json:"cur"`
	Max  int      `json:"max"`
	List []Player `json:"list"`
}

type Info struct {
	Name       string     `json:"name"`
	Map        string     `json:"map"`
	Mode       string     `json:"mode"`
	Vers       string     `json:"vers"`
	PlayerInfo PlayerInfo `json:"players"`
}

type Game interface {
	GameInfo() (Info, error)
	ChangeMap(name string) error
	Command(cmd string) error
}
