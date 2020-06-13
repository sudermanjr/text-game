package game

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/SolarLune/dngn"
	"k8s.io/klog"

	"github.com/sudermanjr/text-game/pkg/utils"
)

func newBSPLevel(level *dngn.Room, splits int) {
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
	selection := level.Select()
	selection.Fill('W')

	level.GenerateDrunkWalk('F', pct)
	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, level.Width-2, level.Height-2)).Fill('W')
}

func newRoomLevel(level *dngn.Room) {
	selection := level.Select()
	selection.Fill('W')

	// numRooms := rand.Intn(10) + 10
	numRooms := 2
	roomPositions := level.GenerateRandomRooms('F', numRooms, 5, 5, 20, 20, false)

	klog.V(5).Info("attempting to connect rooms")
	for p := 0; p < len(roomPositions); p++ {
		if p < len(roomPositions)-1 {
			var (
				x  int
				y  int
				x2 int
				y2 int
			)

			room1Selection := level.SelectContiguous(roomPositions[p][0], roomPositions[p][1])
			room2Selection := level.SelectContiguous(roomPositions[p+1][0], roomPositions[p+1][1])
			var roomsConnected = false
			for _, room1Coord := range room1Selection.Cells {
				for _, room2Coord := range room2Selection.Cells {
					if room1Coord[0] == room2Coord[0] {
						x = room1Coord[0]
						y = room1Coord[1]
						x2 = room1Coord[0]
						y2 = room2Coord[1]
						level.DrawLine(x, y, x2, y2, 'F', 1, false)
						roomsConnected = true
					}
					if room1Coord[1] == room2Coord[1] {
						x = room1Coord[0]
						y = room1Coord[1]
						x2 = room2Coord[0]
						y2 = room2Coord[1]
						level.DrawLine(x, y, x2, y2, 'F', 1, false)
						roomsConnected = true
					}
					if roomsConnected {
						break
					}
				}
				if roomsConnected {
					break
				}
			}
			// x := roomPositions[p][0]
			// y := roomPositions[p][1]

			// x2 := roomPositions[p+1][0]
			// y2 := roomPositions[p+1][1]

		}
	}
	klog.V(5).Info("done connecting rooms")

	// Build the outer walls
	selection.RemoveSelection(selection.ByArea(1, 1, level.Width-2, level.Height-2)).Fill('W')
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
