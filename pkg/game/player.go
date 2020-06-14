package game

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
	"k8s.io/klog"
)

// Player is the player struct for the main player
// The player is the base of the entire game
type Player struct {
	*tl.Entity
	// Name is the name of the character
	Name  string `json:"name"`
	prevX int
	prevY int
	// Char is the character used to represent the player. Usually @
	Char rune `json:"char"`
	// Color is the color of the character player
	Color tl.Attr `json:"color"`
	// Text is the text object that holds status
	Text *tl.Text `json:"text"`
	// BaseText is what always shows up in the message
	BaseText string `json:"baseText"`
	// CurrentLevel is what level the character is on
	CurrentLevel int `json:"currentLevel"`
	// Game is the game object for the character
	Game *tl.Game `json:"game"`
	// Width the level width
	Width int `json:"width"`
	// Height is the level height
	Height int `json:"height"`
	// ShowFPS controls whether the fps message is shown
	ShowFPS bool `json:"showFPS"`
	// FPS is the game setting fps
	Fps float64 `json:"fps"`
	// MapType is the type of levels generated. Noramlly rooms
	MapType string `json:"mapType"`
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
			player.newLevel()
		case false:
			player.setMessage("a staircase leading up")
		}
	}

	// //Hallway
	// if _, ok := collision.(*HallwayTile); ok {
	// 	player.setMessage("an empty hallway space")
	// }

	// //Room
	// if _, ok := collision.(*RoomTile); ok {
	// 	player.setMessage("an empty room space")
	// }
}

// Draw is the draw function for the player
func (player *Player) Draw(screen *tl.Screen) {
	// Camera centering (disabled because it's annoying):
	// screenWidth, screenHeight := screen.Size()
	// x, y := player.Position()
	// player.level.SetOffset(screenWidth/2-x, screenHeight/2-y)

	// Set the player color
	player.SetCell(0, 0, &tl.Cell{
		Fg: player.Color,
		Ch: player.Char,
	})
	player.Entity.Draw(screen)
}

// Tick is the player control
func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		klog.V(7).Infof("pressed key: %v %s", event.Key, string(event.Ch))
		player.prevX, player.prevY = player.Position()
		switch event.Key {
		case tl.KeyArrowRight:
			player.SetPosition(player.prevX+1, player.prevY)
		case tl.KeyArrowLeft:
			player.SetPosition(player.prevX-1, player.prevY)
		case tl.KeyArrowUp:
			player.SetPosition(player.prevX, player.prevY-1)
		case tl.KeyArrowDown:
			player.SetPosition(player.prevX, player.prevY+1)
		case 0:
			switch event.Ch {
			case '>':
				player.setMessage("going down")
			case '<':
				player.setMessage("going up")
			case 'c':
				player.setMessage("close a door")
			}
		}
	}
}

// NewPlayer generates a new character
func NewPlayer(char rune, height int, width int, name string, mapType string, fps float64) *Player {
	player := &Player{
		Entity:       tl.NewEntity(0, 0, 1, 1),
		Color:        tl.ColorRed,
		CurrentLevel: 0,
		BaseText:     string(char) + name,
		Text:         tl.NewText(0, height+1, "", tl.ColorCyan, tl.ColorBlack),
		Game:         tl.NewGame(),
		Width:        width,
		Height:       height,
		MapType:      mapType,
		Fps:          fps,
		Char:         char,
	}
	player.Game.Screen().SetFps(fps)
	player.newLevel()
	player.Game.Screen().AddEntity(player.Text)

	player.setMessage("Welcome to the game!")
	return player
}

func (player *Player) setMessage(message string) {
	text := fmt.Sprintf("%s Level:%d       %s", player.BaseText, player.CurrentLevel, message)
	player.Text.SetText(text)
}

// Start starts the game
func (player *Player) Start() {

	if player.ShowFPS {
		player.Game.Screen().AddEntity(tl.NewFpsText(0, 0, tl.ColorRed, tl.ColorDefault, 0.5))
	}
	klog.V(6).Info("starting the game")
	player.Game.Start()
}
