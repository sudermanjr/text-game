package game

import (
	tl "github.com/JoelOtter/termloop"
)

// WallBlock is a single wall block element
type WallBlock struct {
	*tl.Entity
	level     *tl.BaseLevel
	breakable bool
}

// Staircase is down > or up >
type StairCase struct {
	*tl.Entity
	level *tl.BaseLevel
	down  bool
}

// HallwayTile is a single space in a hallway
type HallwayTile struct {
	*tl.Entity
	level *tl.BaseLevel
}

// RoomTile is one square in a room
type RoomTile struct {
	*tl.Entity
	level *tl.BaseLevel
}

// Draw the hallway tiles
func (hallwayTile *HallwayTile) Draw(screen *tl.Screen) {
	cell := &tl.Cell{
		Bg: bgColor,
	}
	hallwayTile.SetCell(0, 0, cell)
	hallwayTile.Entity.Draw(screen)
}

// Draw the Room Tiles
func (roomTile *RoomTile) Draw(screen *tl.Screen) {
	cell := &tl.Cell{
		Bg: bgColor,
	}
	roomTile.SetCell(0, 0, cell)
	roomTile.Entity.Draw(screen)
}

// Draw draws a staircase
func (stair *StairCase) Draw(screen *tl.Screen) {
	cell := &tl.Cell{
		Fg: stairColor,
	}
	if stair.down {
		cell.Ch = '>'
	} else {
		cell.Ch = '<'
	}
	stair.SetCell(0, 0, cell)
	stair.Entity.Draw(screen)
}

// Draw is draw function that creates wall blocks
func (wallblock *WallBlock) Draw(screen *tl.Screen) {
	wallblock.SetCell(0, 0,
		&tl.Cell{
			Ch: wallChar,
			Fg: wallColor,
		})
	wallblock.Entity.Draw(screen)
}

// Door is a door entity
type Door struct {
	*tl.Entity
	level *tl.BaseLevel
	open  bool
}

// Draw is the draw function for doors
func (door *Door) Draw(screen *tl.Screen) {
	switch door.open {
	case true:
		door.SetCell(0, 0, &tl.Cell{
			Ch: openDoorChar,
			Fg: doorColor,
		})
	case false:
		door.SetCell(0, 0, &tl.Cell{
			Ch: closedDoorChar,
			Fg: doorColor,
		})
	}
	door.Entity.Draw(screen)
}
