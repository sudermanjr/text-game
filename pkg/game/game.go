package game

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/SolarLune/dngn"
)

var (
	floorChar  = '.'
	wallChar   = '#'
	playerChar = '@'
	doorChar   = '+'
)

func newLevelData(w, h int) [][]rune {
	var GameMap *dngn.Room
	GameMap = dngn.NewRoom(w, h)

	// Selections are structs, so we can store Selections in variables to store the "view" of the data.
	selection := GameMap.Select()
	selection.Fill(floorChar)

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, w-2, h-2)).Fill(wallChar)

	// BSP
	GameMap.GenerateBSP(wallChar, doorChar, 25)

	// Make the walls thicker
	GameMap.Select().ByRune('.').ByNeighbor(wallChar, 1, false).Fill(wallChar)

	// Pick a player starting position
	// Currently this is just the center or as near to it as we can get.
	// TODO: Find a better way to do this
	centerX, centerY := GameMap.Center()
	for {
		if GameMap.Get(centerX, centerY) == floorChar {
			GameMap.Set(centerX, centerY, playerChar)
			break
		}
		centerX = centerX + 1
		centerY = centerY + 1
	}
	return GameMap.Data
}

// newLevel builds a new level for the game
func newLevel(g *tl.Game, w, h int) *tl.BaseLevel {
	layout := newLevelData(w, h)
	l := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	g.Screen().SetLevel(l)
	g.Log("Building level with width %d and height %d", w, h)

	for i, row := range layout {
		for j, y := range row {
			switch y {
			case wallChar:
				wall := &WallBlock{
					Entity: tl.NewEntity(j, i, 1, 1),
					level:  l,
				}
				l.AddEntity(wall)
			case playerChar:
				player := Player{
					Entity: tl.NewEntity(j, i, 1, 1),
					level:  l,
				}
				// Set the character at position (0, 0) on the entity.
				player.SetCell(0, 0, &tl.Cell{Fg: tl.ColorRed, Ch: playerChar})
				l.AddEntity(&player)
			}
		}
	}
	return l
}

// New builds a new game and returns it
func New(w int, h int, fps float64) *tl.Game {
	instance := tl.NewGame()
	instance.Screen().SetFps(fps)
	level := newLevel(instance, w, h)

	instance.Screen().SetLevel(level)
	return instance
}
