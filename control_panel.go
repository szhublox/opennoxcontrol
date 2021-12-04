package opennoxcontrol

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/szhublox/opennoxcontrol/assets"
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

func NewControlPanel(game Game, opts *Options) *ControlPanel {
	if opts == nil {
		// everything defaults to false
		opts = &Options{}
	}
	tmpl, err := template.ParseFS(assets.Templates(), "*.gohtml")
	if err != nil {
		panic(err)
	}
	cp := &ControlPanel{
		g:    game,
		mux:  http.NewServeMux(),
		tmpl: tmpl,
		opts: *opts,
		maps: defaultMaps,
	}
	cp.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(assets.Static()))))
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
	tmpl *template.Template
	opts Options
	maps []string
}

func (cp *ControlPanel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cp.mux.ServeHTTP(w, r)
}

func (cp *ControlPanel) redirectToRoot(w http.ResponseWriter) {
	fmt.Fprintf(w, "<html><head><meta http-equiv=\"Refresh\""+
		" content=\"0; url=/\" /></head></html>")
}

func (cp *ControlPanel) rootHandler(w http.ResponseWriter, r *http.Request) {
	info, err := cp.g.GameInfo()
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Couldn't get game data.")
		return
	}
	err = cp.tmpl.ExecuteTemplate(w, "root", struct {
		Opts Options
		Maps []string
		Info Info
	}{
		Opts: cp.opts,
		Maps: cp.maps,
		Info: info,
	})
	if err != nil {
		panic(err)
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

	cp.redirectToRoot(w)
}

func (cp *ControlPanel) commandHandler(w http.ResponseWriter, r *http.Request) {
	if cp.opts.AllowCommands {
		r.ParseForm()
		var data = r.Form.Get("data")

		cp.g.Command(data)
	}

	cp.redirectToRoot(w)
}
