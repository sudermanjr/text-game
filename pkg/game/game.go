package game

import (
	tl "github.com/JoelOtter/termloop"
	"k8s.io/klog"
)

// Instance is all the config necessary for running the game
type Instance struct {
	Width   int
	Height  int
	Fps     float64
	MapType string
	ShowFPS bool
}

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

// Run sets up and starts a game instance
func (instance *Instance) Run() {
	g := tl.NewGame()
	g.Screen().SetFps(instance.Fps)
	level := newLevel(g, instance.Width, instance.Height, instance.MapType)

	g.Screen().SetLevel(level)

	if instance.ShowFPS {
		g.Screen().AddEntity(tl.NewFpsText(0, 0, tl.ColorRed, tl.ColorDefault, 0.5))
	}
	klog.V(6).Info("starting the game")
	g.Start()
}
