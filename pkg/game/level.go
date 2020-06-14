package game

import (
	"fmt"
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
	selection.Fill('R')

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, level.Width-2, level.Height-2)).Fill(' ')

	// BSP
	level.GenerateBSP(' ', 'D', splits)

	// Make the walls thicker
	//GameMap.Select().ByRune('F').ByNeighbor('W', 1, false).Fill('W')
}

func newDrunkWalkLevel(level *dngn.Room, pct float32) {
	klog.V(5).Info("building new drunkwalk level")
	selection := level.Select()
	selection.Fill(' ')

	level.GenerateDrunkWalk('F', pct)

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, level.Width-2, level.Height-2)).Fill('W')
}

func newRoomLevel(level *dngn.Room) {
	xidx := 0
	yidx := 1
	klog.V(5).Info("building new rooms level")
	selection := level.Select()
	selection.Fill(' ')

	numRooms := rand.Intn(12) + 6

	// Set these for debugging room creation
	// level.SetSeed(11)
	// numRooms = 3
	roomPositions := level.GenerateRandomRooms('R', numRooms, 5, 5, 10, 10, false)

	klog.V(9).Infof("room positions: %v", roomPositions)
	klog.V(5).Info("attempting to generate hallways")

	roomMap := make(map[int]dngn.Selection)
	for idx, room := range roomPositions {
		roomMap[idx] = level.SelectContiguous(room[xidx], room[yidx])
	}

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

			for _, room1Coord := range roomMap[a].Cells {
				for _, room2Coord := range roomMap[b].Cells {
					klog.V(8).Infof("room1: %v - room2: %v", room1Coord, room2Coord)
					if room1Coord[xidx] == room2Coord[xidx] || room1Coord[yidx] == room2Coord[yidx] {
						x = room1Coord[xidx]
						y = room1Coord[yidx]
						x2 = room2Coord[xidx]
						y2 = room2Coord[yidx]
						level.DrawLine(x, y, x2, y2, 'H', 1, false)
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
	// Fix up the rooms to clear out errant hallways
	for _, room := range roomMap {
		room.Fill('R')
	}

	klog.V(5).Info("done connecting rooms")

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, level.Width-2, level.Height-2)).Fill('W')

	klog.V(8).Info("room generation complete")
	fmt.Println(level.DataToString())
}

func placePlayer(level *dngn.Room) {
	openFloor := level.Select().ByRune('R').ByNeighbor('R', 8, true).Cells
	if len(openFloor) < 1 {
		klog.Error("no open room floor found")
		return
	}
	randomFloor := rand.Intn(len(openFloor))
	randomX := openFloor[randomFloor][0]
	randomY := openFloor[randomFloor][1]

	klog.V(4).Infof("Placing player at random location: %d, %d", randomX, randomY)
	level.Set(randomX, randomY, '@')
}

func placeStaircase(level *dngn.Room, down bool) {
	// A room with a hallway attached and a space away from the walls
	openFloor := level.Select().ByRune('R').ByNeighbor('R', 8, true).Cells
	if len(openFloor) < 1 {
		klog.Error("no open floor found")
		return
	}
	randomFloor := rand.Intn(len(openFloor))
	randomX := openFloor[randomFloor][0]
	randomY := openFloor[randomFloor][1]

	klog.V(4).Infof("Placing down staircase at random location: %d, %d", randomX, randomY)

	if down {
		level.Set(randomX, randomY, '>')
	} else {
		level.Set(randomX, randomY, '<')
	}
}

func placeDoors(level *dngn.Room) {
	// A hallway space with three room spaces next to it as well as two walls that are non-diagonal

	// Desired door placement
	//      R
	// HHHHDR
	//      R
	doorLocations := level.Select().ByRune('H').ByNeighbor('R', 3, true).By(func(x, y int) bool {
		return (level.Get(x+1, y) == ' ' && level.Get(x-1, y) == ' ') || (level.Get(x, y+1) == ' ' && level.Get(x, y-1) == ' ')
	})
	doorLocations.ByPercentage(.5).Fill('D')
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
	placeStaircase(GameMap, true)
	placeDoors(GameMap)
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
			case ' ':
				klog.V(10).Infof("adding wall at %d, %d", j, i)
				wall := &WallBlock{
					Entity:    tl.NewEntity(j, i, 1, 1),
					level:     l,
					breakable: false,
				}
				l.AddEntity(wall)
			case 'W':
				klog.V(10).Infof("adding wall at %d, %d", j, i)
				wall := &WallBlock{
					Entity:    tl.NewEntity(j, i, 1, 1),
					level:     l,
					breakable: true,
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
			case '>':
				klog.V(8).Infof("setting down staircase at %d, %d", j, i)
				downStair := &StairCase{
					Entity: tl.NewEntity(j, i, 1, 1),
					level:  l,
					down:   true,
				}
				l.AddEntity(downStair)
			case 'R':
				klog.V(8).Infof("drawing room tile at %d, %d", j, i)
				l.AddEntity(&RoomTile{Entity: tl.NewEntity(j, i, 1, 1), level: l})
			case 'H':
				klog.V(8).Infof("drawing hallway tile at %d, %d", j, i)
				l.AddEntity(&HallwayTile{Entity: tl.NewEntity(j, i, 1, 1), level: l})
			}
		}
	}
	return l
}
