package game

import (
	tl "github.com/JoelOtter/termloop"
)

// Package variables to control the look and feel
var (
	floorChar      = '.'
	wallChar       = '#'
	playerChar     = '@'
	closedDoorChar = '+'
	openDoorChar   = '/'
)

// New builds a new game and returns it
func New(w int, h int, fps float64) *tl.Game {
	instance := tl.NewGame()
	instance.Screen().SetFps(fps)
	level := newLevel(instance, w, h, 20)

	instance.Screen().SetLevel(level)
	return instance
}
