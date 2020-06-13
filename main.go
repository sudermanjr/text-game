package main

import (
	"math/rand"
	"time"

	"github.com/sudermanjr/text-game/cmd"
)

var (
	// version is set during build
	version = "development"
	// commit is set during build
	commit = "n/a"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmd.Execute(version, commit)
}
