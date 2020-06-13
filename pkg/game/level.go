package game

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/SolarLune/dngn"

	"github.com/sudermanjr/text-game/pkg/utils"
)

func newLevelData(w, h int, splits int) [][]rune {
	var GameMap *dngn.Room
	GameMap = dngn.NewRoom(w, h)

	// Selections are structs, so we can store Selections in variables to store the "view" of the data.
	selection := GameMap.Select()
	selection.Fill('F')

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, w-2, h-2)).Fill('W')

	// BSP
	GameMap.GenerateBSP('W', 'D', splits)

	// Make the walls thicker
	//GameMap.Select().ByRune('F').ByNeighbor('W', 1, false).Fill('W')

	// Pick a player starting position
	// Currently this is just the center or as near to it as we can get.
	// TODO: Find a better way to do this
	centerX, centerY := GameMap.Center()
	for {
		if GameMap.Get(centerX, centerY) == 'F' {
			GameMap.Set(centerX, centerY, '@')
			break
		}
		centerX = centerX + 1
		centerY = centerY + 1
	}
	return GameMap.Data
}

// newLevel builds a new level for the game
func newLevel(g *tl.Game, w, h int, splits int) *tl.BaseLevel {
	layout := newLevelData(w, h, splits)
	l := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	g.Screen().SetLevel(l)
	g.Log("Building level with width %d and height %d", w, h)

	for i, row := range layout {
		for j, y := range row {
			switch y {
			case 'W':
				wall := &WallBlock{
					Entity: tl.NewEntity(j, i, 1, 1),
					level:  l,
				}
				l.AddEntity(wall)
			case '@':
				player := Player{
					Entity: tl.NewEntity(j, i, 1, 1),
					level:  l,
					color:  tl.ColorRed,
				}
				l.AddEntity(&player)
			case 'D':
				door := &Door{
					Entity: tl.NewEntity(j, i, 1, 1),
					level:  l,
					open:   utils.RandomBool(),
				}
				l.AddEntity(door)
			}
		}
	}
	return l
}
