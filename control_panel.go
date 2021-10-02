package main

import "net/http"

func NewControlPanel(allowCmd bool) *ControlPanel {
	cp := &ControlPanel{
		mux:      http.NewServeMux(),
		allowCmd: allowCmd,
	}
	cp.mux.HandleFunc("/", cp.rootHandler)
	cp.mux.HandleFunc("/map/", cp.mapHandler)
	if allowCmd {
		cp.mux.HandleFunc("/cmd/", cp.commandHandler)
	}
	return cp
}

type ControlPanel struct {
	mux      *http.ServeMux
	allowCmd bool
}

func (cp *ControlPanel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cp.mux.ServeHTTP(w, r)
}
