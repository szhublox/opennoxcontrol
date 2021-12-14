package opennoxcontrol

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

var defaultMaps = []string{
	"autumn", "beneath", "blsphmy", "BluDeath", "Bunker",
	"Bywaters", "caverna", "Connect4", "courtyrd", "darkLib", "Decrypt",
	"deserts", "Dewls", "DownDown", "DthTmple", "Duel", "Ennead", "Estate",
	"FortNox", "Fortress", "gloomy", "grotto", "headache", "Impulse",
	"Inferno", "ktdefend", "Library", "LostTomb", "luckspin", "ManaMine",
	"minetomb", "MiniMine", "MnaVault", "Notso", "nowhere", "Oasis",
	"outtaB", "pullrope", "sjdream", "sjglav", "sjhom", "SJJC", "sjscary",
	"skycolor", "smth", "SpyFort", "tbhold", "TheGuild", "TreeHaus",
	"TriLevel", "tropix", "uden", "uwcastle", "waterlib", "Waterwar",
	"Whirl", "winter", "WorldEnd",
}

type Options struct {
	AllowCommands  bool
	AllowMapChange bool
}

// NewControlPanel creates a new Nox game control panel.
func NewControlPanel(game Game, opts *Options) *ControlPanel {
	if opts == nil {
		// everything defaults to false
		opts = &Options{}
	}
	cp := &ControlPanel{
		g:    game,
		mux:  http.NewServeMux(),
		opts: *opts,
		maps: defaultMaps,
	}
	cp.mux.HandleFunc("/", cp.rootHandler)
	if opts.AllowMapChange || opts.AllowCommands {
		cp.mux.HandleFunc("/map/", cp.mapHandler)
	}
	if opts.AllowCommands {
		cp.mux.HandleFunc("/cmd/", cp.commandHandler)
	}
	if list, err := game.ListMaps(); err == nil {
		cp.maps = nil
		for _, m := range list {
			cp.maps = append(cp.maps, m.Name)
		}
	}
	return cp
}

type ControlPanel struct {
	g    Game
	mux  *http.ServeMux
	opts Options
	maps []string
}

func (cp *ControlPanel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cp.mux.ServeHTTP(w, r)
}

func (cp *ControlPanel) refresh_to_root(w http.ResponseWriter) {
	fmt.Fprintf(w, "<html><head><meta http-equiv=\"Refresh\""+
		" content=\"0; url=/\" /></head></html>")
}

func (cp *ControlPanel) print_players_table(w http.ResponseWriter, info Info) {
	fmt.Fprintf(w, `<table summary="server details">
<tr><td>Server Name</td><td>%s</td></tr>
<tr><td>Current Mode</td><td>%s</td></tr>
<tr><td>Current Map</td><td>%s</td></tr>
<tr><td>Player Count</td><td>%d / %d</td></tr>`,
		info.Name, info.Mode, info.Map, info.PlayerInfo.Cur,
		info.PlayerInfo.Max)

	if info.PlayerInfo.Cur > 0 {
		fmt.Fprintf(w, "<tr><td>Players</td><td>")
		for i := 0; i < info.PlayerInfo.Cur; i++ {
			fmt.Fprintf(w, "%s the %s",
				info.PlayerInfo.List[i].Name,
				info.PlayerInfo.List[i].Class)
			if i < info.PlayerInfo.Cur-1 {
				fmt.Fprintf(w, "<br />\n")
			}
		}
		fmt.Fprintf(w, "</td></tr>")
	}

	fmt.Fprintf(w, "</table>\n")
}

func (cp *ControlPanel) print_map_form(w http.ResponseWriter, info Info) {
	fmt.Fprintf(w, "<br />\n")
	if !cp.opts.AllowCommands {
		fmt.Fprintf(w,
			"\n<b>Map change only allowed when "+
				"the server is empty.</b>")
	}
	fmt.Fprintf(w, `<form action="/map/" method="POST">
<label>Change Map</label>
<select name="data">`)
	for _, name := range cp.maps {
		fmt.Fprintf(w, `<option value="%s"`, name)
		if strings.EqualFold(name, info.Map) {
			fmt.Fprintf(w, " selected")
		}
		fmt.Fprintf(w, ">%s</option>\n", name)
	}
	fmt.Fprintf(w, `</select><input type="submit" value="Submit" /></form>`)
}

func (cp *ControlPanel) print_command_form(w http.ResponseWriter) {
	fmt.Fprintf(w, `<br /><form action="/cmd/" method="post">
<label>Command</label>
<input type="text" name="data" />
<input type="submit" value="Submit" />
</form>
`)
}

func (cp *ControlPanel) rootHandler(w http.ResponseWriter, r *http.Request) {
	var info Info

	fmt.Fprintf(w, "<!DOCTYPE html>\n"+
		"<html><head><title>OpenNox Server Control</title>\n"+
		"<style>td { padding-right: 15px; }</style>\n"+
		"</head>\n"+
		"<body>\n")
	defer fmt.Fprintln(w, `</body></html>`)

	info, err := cp.g.GameInfo()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Couldn't get game data.")
		return
	}

	cp.print_players_table(w, info)
	if cp.opts.AllowMapChange || cp.opts.AllowCommands {
		cp.print_map_form(w, info)
	}
	if cp.opts.AllowCommands {
		cp.print_command_form(w)
	}
}

func (cp *ControlPanel) mapHandler(w http.ResponseWriter, r *http.Request) {
	if cp.opts.AllowMapChange || cp.opts.AllowCommands {
		r.ParseForm()

		info, err := cp.g.GameInfo()
		if err != nil {
			// silently return since we can't print and expect refresh to work
			return
		}
		var data = r.Form.Get("data")

		if (cp.opts.AllowCommands || info.PlayerInfo.Cur == 0) && len(data) > 0 {
			cp.g.ChangeMap(data)
		}
	}

	cp.refresh_to_root(w)
}

func (cp *ControlPanel) commandHandler(w http.ResponseWriter, r *http.Request) {
	if cp.opts.AllowCommands {
		r.ParseForm()
		var data = r.Form.Get("data")

		cp.g.Command(data)
	}

	cp.refresh_to_root(w)
}
