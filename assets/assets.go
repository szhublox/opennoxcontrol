package assets

import "embed"

//go:embed *.gohtml
var templates embed.FS

//go:embed css/*
var static embed.FS

func Templates() embed.FS {
	return templates
}

func Static() embed.FS {
	return static
}
