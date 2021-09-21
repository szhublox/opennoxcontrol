# opennoxcontrol

web server that runs on :8080 to talk to the opennox 1.8+ server api using token

    X-Token: xyz

for now will die on any error, including opennox server not running. included is an experimental global variable `bind_local` at the top of opennoxcontrol.go which changes whether the web server is externally accessible, and lowers its powers to match - since there's currently no authentication, anybody in the world shouldn't be able to change the map when people are playing
