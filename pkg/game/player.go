package game

import (
	tl "github.com/JoelOtter/termloop"
)

// Player is the player struct for the main player
type Player struct {
	*tl.Entity
	prevX        int
	prevY        int
	color        tl.Attr
	level        *tl.BaseLevel
	text         *tl.Text
	baseText     string
	currentLevel int
}

// Collide is the player's collision processing
func (player *Player) Collide(collision tl.Physical) {
	// Wall Collision
	if _, ok := collision.(*WallBlock); ok {
		player.SetPosition(player.prevX, player.prevY)
		player.setMessage("you run into a wall")
	}

	// Door
	if _, ok := collision.(*Door); ok {
		door := collision.(*Door)
		switch door.open {
		case true:
			player.SetPosition(collision.Position())
			player.setMessage("an open door")
		case false:
			if door.locked {
				player.setMessage("this door is locked")
			} else {
				player.setMessage("the door is unlocked. you open it")
				door.open = true
			}
			player.SetPosition(player.prevX, player.prevY)
		}
	}

	// Staircase
	if _, ok := collision.(*StairCase); ok {
		switch collision.(*StairCase).down {
		case true:
			player.setMessage("a staircase leading down")
		case false:
			player.setMessage("a staircase leading up")
		}

	}
}

// Draw is the draw function for the player
func (player *Player) Draw(screen *tl.Screen) {
	// Camera centering (disabled because it's annoying):
	// screenWidth, screenHeight := screen.Size()
	// x, y := player.Position()
	// player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)

	// Set the player color
	player.SetCell(0, 0, &tl.Cell{
		Fg: player.color,
		Ch: playerChar,
	})
	player.Entity.Draw(screen)
}

// Tick is the player control
func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		player.prevX, player.prevY = player.Position()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			player.SetPosition(player.prevX+1, player.prevY)
		case tl.KeyArrowLeft:
			player.SetPosition(player.prevX-1, player.prevY)
		case tl.KeyArrowUp:
			player.SetPosition(player.prevX, player.prevY-1)
		case tl.KeyArrowDown:
			player.SetPosition(player.prevX, player.prevY+1)
		}
	}
}

// NewPlayer generates a new character
func NewPlayer(x int, y int, char rune, text *tl.Text) *Player {
	player := &Player{
		Entity:       tl.NewEntity(x, y, 1, 1),
		color:        tl.ColorRed,
		text:         text,
		currentLevel: 0,
		baseText:     "",
	}
	return player
}

func (player *Player) setMessage(message string) {
	player.text.SetText(player.baseText + message)
}
