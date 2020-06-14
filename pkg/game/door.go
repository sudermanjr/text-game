package game

import (
	tl "github.com/JoelOtter/termloop"
)

// Door is a door entity
type Door struct {
	*tl.Entity
	level  *tl.BaseLevel
	open   bool
	locked bool
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

// Collide is what happens when a player hits a door
func (door *Door) Collide(collision tl.Physical) {
	// Player Collision
	// if _, ok := collision.(*Player); ok {
	// 	if !door.open && !door.locked {
	// 		collision.(*Player).setMessage("the door is unlocked. you open it")
	// 		door.open = true
	// 	}
	// }
}