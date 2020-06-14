package game

import (
	"math/rand"

	tl "github.com/JoelOtter/termloop"
	"github.com/SolarLune/dngn"
	"k8s.io/klog"

	"github.com/sudermanjr/text-game/pkg/utils"
)

func newBSPLevel(level *dngn.Room, splits int) {
	klog.V(5).Info("building new bsp level")
	// Selections are structs, so we can store Selections in variables to store the "view" of the data.
	selection := level.Select()
	selection.Fill('F')

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, level.Width-2, level.Height-2)).Fill('W')

	// BSP
	level.GenerateBSP('W', 'D', splits)

	// Make the walls thicker
	//GameMap.Select().ByRune('F').ByNeighbor('W', 1, false).Fill('W')
}

func newDrunkWalkLevel(level *dngn.Room, pct float32) {
	klog.V(5).Info("building new drunkwalk level")
	selection := level.Select()
	selection.Fill('W')

	level.GenerateDrunkWalk('F', pct)

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, level.Width-2, level.Height-2)).Fill('W')
}

func newRoomLevel(level *dngn.Room) {
	xidx := 0
	yidx := 1
	klog.V(5).Info("building new rooms level")
	selection := level.Select()
	selection.Fill('W')

	// level.SetSeed(10)

	numRooms := rand.Intn(12) + 6
	// numRooms = 3
	roomPositions := level.GenerateRandomRooms('F', numRooms, 5, 5, 10, 10, false)

	klog.V(9).Infof("room positions: %v", roomPositions)
	klog.V(5).Info("attempting to connect rooms")

	for a, room1 := range roomPositions {
		var (
			x  int
			y  int
			x2 int
			y2 int
		)
		for b, room2 := range roomPositions {
			if a == b {
				break
			}
			var connected bool = false
			klog.V(7).Infof("checking room connections between room %d and %d", a, b)
			klog.V(6).Infof("a: %v, b: %v", room1, room2)
			room1Cells := level.SelectContiguous(room1[xidx], room1[yidx]).Cells
			room2Cells := level.SelectContiguous(room2[xidx], room2[yidx]).Cells

			for _, room1Coord := range room1Cells {
				for _, room2Coord := range room2Cells {
					klog.V(8).Infof("room1: %v - room2: %v", room1Coord, room2Coord)
					if room1Coord[xidx] == room2Coord[xidx] || room1Coord[yidx] == room2Coord[yidx] {
						x = room1Coord[xidx]
						y = room1Coord[yidx]
						x2 = room2Coord[xidx]
						y2 = room2Coord[yidx]
						level.DrawLine(x, y, x2, y2, 'F', 1, false)
						klog.V(5).Infof("connected room %d to room %d via %v %v", a, b, room1Coord, room2Coord)
						connected = true
						break
					}
				}
				if connected {
					break
				}
			}
		}
	}

	// klog.V(5).Info("done connecting rooms")

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, level.Width-2, level.Height-2)).Fill('W')
	klog.V(8).Info("built outer walls. room generation complete")
}

func placePlayer(level *dngn.Room) {
	openFloor := level.Select().ByRune('F')
	randomFloor := rand.Intn(len(openFloor.Cells))
	randomX := openFloor.Cells[randomFloor][0]
	randomY := openFloor.Cells[randomFloor][1]

	klog.V(4).Infof("Placing player at random location: %d, %d", randomX, randomY)
	level.Set(randomX, randomY, '@')
}

func newLevelData(w, h int, levelType string) [][]rune {
	var GameMap *dngn.Room
	GameMap = dngn.NewRoom(w, h)

	switch levelType {
	case "bsp":
		newBSPLevel(GameMap, 20)
	case "drunkwalk":
		newDrunkWalkLevel(GameMap, 0.5)
	case "rooms":
		newRoomLevel(GameMap)
	}
	placePlayer(GameMap)
	return GameMap.Data
}

// newLevel builds a new level for the game
func newLevel(g *tl.Game, w, h int, mapType string) *tl.BaseLevel {
	layout := newLevelData(w, h, mapType)
	l := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	g.Screen().SetLevel(l)
	g.Log("Building level with width %d and height %d", w, h)

	for i, row := range layout {
		for j, y := range row {
			switch y {
			case 'W':
				klog.V(10).Infof("adding wall at %d, %d", j, i)
				wall := &WallBlock{
					Entity: tl.NewEntity(j, i, 1, 1),
					level:  l,
				}
				l.AddEntity(wall)
			case '@':
				klog.V(8).Infof("setting player at %d, %d", j, i)
				player := Player{
					Entity: tl.NewEntity(j, i, 1, 1),
					level:  l,
					color:  tl.ColorRed,
				}
				l.AddEntity(&player)
			case 'D':
				klog.V(8).Infof("adding door at %d, %d", j, i)
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
