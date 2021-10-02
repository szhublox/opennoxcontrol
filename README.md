# OpenNox Server Control Panel

Web server that runs on :8080 to talk to the OpenNox v1.8+ server API using token.

## Building

```
go build ./cmd/opennox-control
```

## Running

1. Start OpenNox server with `NOX_API_TOKEN=xyz` environment variable, where `xyz` is a random string.
2. Start control panel: `./opennox-control -token="xyz" -host=127.0.0.1:8080`

This will start control panel on [127.0.0.1:8080](http://127.0.0.1:8080).

If you want to allow external access, pass `-host=:8080` option instead.

Make sure to adjust `-cmd=false` if needed. Since there's currently no authentication, anybody in the world could change the map when people are playing, if that flag is set to true.
