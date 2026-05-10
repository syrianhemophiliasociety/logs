package assets

import "embed"

//go:embed *
var fs embed.FS

func FS() embed.FS {
	return fs
}
