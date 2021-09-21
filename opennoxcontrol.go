package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var bind_local = true // 127.0.0.1 instead of 0.0.0.0
var bind_port = "8080"

var mapnames = [58]string{"autumn", "beneath", "blsphmy", "BluDeath", "Bunker",
	"Bywaters", "caverna", "Connect4", "courtyrd", "darkLib", "Decrypt",
	"deserts", "Dewls", "DownDown", "DthTmple", "Duel", "Ennead", "Estate",
	"FortNox", "Fortress", "gloomy", "grotto", "headache", "Impulse",
	"Inferno", "ktdefend", "Library", "LostTomb", "luckspin", "ManaMine",
	"minetomb", "MiniMine", "MnaVault", "Notso", "nowhere", "Oasis",
	"outtaB", "pullrope", "sjdream", "sjglav", "sjhom", "SJJC", "sjscary",
	"skycolor", "smth", "SpyFort", "tbhold", "TheGuild", "TreeHaus",
	"TriLevel", "tropix", "uden", "uwcastle", "waterlib", "Waterwar",
	"Whirl", "winter", "WorldEnd"}

func refresh_to_root(w http.ResponseWriter) {
	fmt.Fprintf(w, "<html><head><meta http-equiv=\"Refresh\""+
		" content=\"0; url=/\" /></head></html>")
}

func print_players_table(w http.ResponseWriter, info Info) {
	fmt.Fprintf(w, "<table summary=\"server details\">\n"+
		"<tr><td>Server Name</td><td>%v</td></tr>\n"+
		"<tr><td>Current Mode</td><td>%v</td></tr>\n"+
		"<tr><td>Current Map</td><td>%v</td></tr>\n"+
		"<tr><td>Player Count</td><td>%v / %v</td></tr>\n",
		info.Name, info.Mode, info.Map, info.PlayerInfo.Cur,
		info.PlayerInfo.Max)

	if info.PlayerInfo.Cur > 0 {
		fmt.Fprintf(w, "<tr><td>Players</td><td>")
		for i := 0; i < info.PlayerInfo.Cur; i++ {
			fmt.Fprintf(w, "%v the %v",
				info.PlayerInfo.List[i].Name,
				info.PlayerInfo.List[i].Class)
			if i < info.PlayerInfo.Cur-1 {
				fmt.Fprintf(w, "<br />\n")
			}
		}
	}

	fmt.Fprintf(w, "</td></tr></table>\n")
}

func print_map_form(w http.ResponseWriter, info Info) {
	fmt.Fprintf(w, "<br />\n")
	if !bind_local {
		fmt.Fprintf(w,
			"\n<b>Map change only allowed when "+
				"the server is empty.</b>")
	}
	fmt.Fprintf(w, "<form action=\"/map/\" method=\"POST\">"+
		"<label>Change Map</label>"+
		"<select name=\"data\">")
	for i := 0; i < len(mapnames); i++ {
		fmt.Fprintf(w, "<option value=\"%v\"", mapnames[i])
		if strings.EqualFold(mapnames[i], info.Map) {
			fmt.Fprintf(w, " selected")
		}
		fmt.Fprintf(w, ">%v</option>\n", mapnames[i])
	}
	fmt.Fprintf(w, "</select>\n"+
		"<input type=\"submit\" value=\"Submit\" />\n"+
		"</form>\n")
}

func print_command_form(w http.ResponseWriter) {
	fmt.Fprintf(w, "<form action=\"/cmd/\" method=\"post\">\n"+
		"<label>Command</label>\n"+
		"<input type=\"text\" name=\"data\" />\n"+
		"<input type=\"submit\" value=\"Submit\" />\n"+
		"</form>\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var info Info

	info = get_info()

	fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \""+
		"-//W3C//DTD HTML 4.01 Transitional//EN\">\n"+
		"<html><head><title>OpenNox Server Control</title>\n"+
		"</head>\n"+
		"<body>\n")

	print_players_table(w, info)
	print_map_form(w, info)

	if bind_local {
		print_command_form(w)
	}
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var info Info = get_info()
	var data = r.Form.Get("data")

	if (bind_local || info.PlayerInfo.Cur == 0) && len(data) > 0 {
		nox_curl_post("map", data)
	}

	refresh_to_root(w)
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	if bind_local {
		r.ParseForm()
		var data = r.Form.Get("data")

		nox_curl_post("cmd", data)
	}

	refresh_to_root(w)
}

func main() {
	var bind_host string

	if bind_local {
		bind_host = "127.0.0.1:" + bind_port
	} else {
		bind_host = "0.0.0.0:" + bind_port
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/map/", mapHandler)

	if bind_local {
		http.HandleFunc("/cmd/", commandHandler)
	}

	log.Fatal(http.ListenAndServe(bind_host, nil))
}
