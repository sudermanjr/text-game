package game

import (
	tl "github.com/JoelOtter/termloop"
)

// WallBlock is a single wall block element
type WallBlock struct {
	*tl.Entity
	level *tl.BaseLevel
}

// Draw is draw function that creates wall blocks
func (wallblock *WallBlock) Draw(screen *tl.Screen) {
	wallblock.SetCell(0, 0,
		&tl.Cell{
			Ch: wallChar,
			Fg: tl.ColorWhite,
		})
	wallblock.Entity.Draw(screen)
}

// Door is a door entity
type Door struct {
	*tl.Entity
	level *tl.BaseLevel
	open  bool
}

func (door *Door) Draw(screen *tl.Screen) {
	switch door.open {
	case true:
		door.SetCell(0, 0, &tl.Cell{
			Ch: openDoorChar,
			Fg: tl.ColorWhite,
		})
	case false:
		door.SetCell(0, 0, &tl.Cell{
			Ch: closedDoorChar,
			Fg: tl.ColorWhite,
		})
	}
	door.Entity.Draw(screen)
}
