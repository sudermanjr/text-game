package map

import "github.com/SolarLune/dngn"



func newGameMap(w, h int) {
	var GameMap *dngn.Room

    // This line creates a new Room. The size is 10x10.
    GameMap = dngn.NewRoom(h, w)

    // This will select the cells the map has, and then fill the selection with "x"s.
    GameMap.Select().Fill('x')

    // Selections are structs, so we can store Selections in variables to store the "view" of the data.
    selection := GameMap.Select()

    // This will run a drunk-walk generation algorithm on the Room. It starts at a random point
    // in the Room, and walks around the Room, placing the value specified (0, in this case)
    // until the room is the percentage provided (0.5, or 50%, in this case) filled.
    GameMap.GenerateDrunkWalk(' ', 0.5)

    // This function will degrade the map slightly, making cells with a ' ' in them randomly turn into a cell with a 'x',
    // depending on how heavily the cell is surrounded by 'x's.
    GameMap.Select().ByRune(' ').Degrade('x')

    // Room.DataToString() will present the data in a nice, easy-to-understand visual format, useful when debugging.
	fmt.Println(GameMap.DataToString())

    // Now we're done! We can use the Room.
	return GameMap
}

func New(w, h) [][]rune {
    map := newGameMap(w, h)
    return [][]rune{}
}