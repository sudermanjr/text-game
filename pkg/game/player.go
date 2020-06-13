package game

import (
	tl "github.com/JoelOtter/termloop"
)

// Player is the player struct for the main player
type Player struct {
	*tl.Entity
	prevX int
	prevY int
	level *tl.BaseLevel
}

// Collide is the player's collision processing
func (player *Player) Collide(collision tl.Physical) {
	// Wall Collision
	if _, ok := collision.(*WallBlock); ok {
		player.SetPosition(player.prevX, player.prevY)
	}
}

// Draw is the draw function for the player
func (player *Player) Draw(screen *tl.Screen) {
	// screenWidth, screenHeight := screen.Size()
	// x, y := player.Position()
	// player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)
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