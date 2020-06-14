package game

import (
	tl "github.com/JoelOtter/termloop"
	"k8s.io/klog"
)

// Package variables to control the look and feel
var (
	floorChar      = '.'
	wallChar       = '#'
	playerChar     = '@'
	closedDoorChar = '+'
	openDoorChar   = '/'
	doorColor      = tl.ColorWhite
	stairColor     = tl.ColorGreen
	wallColor      = tl.ColorWhite
	bgColor        = tl.ColorBlack
)

// New builds a new game and returns it
func New(w int, h int, fps float64, mapType string) *tl.Game {
	instance := tl.NewGame()

	instance.Screen().SetFps(fps)
	level := newLevel(instance, w, h, mapType)

	instance.Screen().SetLevel(level)
	klog.V(6).Info("returning instance of game to be started")
	return instance
}
