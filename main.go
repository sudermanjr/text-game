package main

import (
	"math"
	"strconv"

	tl "github.com/JoelOtter/termloop"

	"github.com/sudermanjr/text-game/pkg/level"
)

var (
	player1Color = tl.ColorRed
	player2Color = tl.ColorGreen
	goalColor    = tl.ColorBlue
	wallColor    = tl.ColorWhite
)

////////////////////////
// Maze generation stuff
////////////////////////
// type Point struct {
// 	x int
// 	y int
// 	p *Point
// }

// func (p *Point) Opposite() *Point {
// 	if p.x != p.p.x {
// 		return &Point{x: p.x + (p.x - p.p.x), y: p.y, p: p}
// 	}
// 	if p.y != p.p.y {
// 		return &Point{x: p.x, y: p.y + (p.y - p.p.y), p: p}
// 	}
// 	return nil
// }

// func adjacents(point *Point, maze [][]rune) []*Point {
// 	res := make([]*Point, 0)
// 	for i := -1; i < 2; i++ {
// 		for j := -1; j < 2; j++ {
// 			if (i == 0 && j == 0) || (i != 0 && j != 0) {
// 				continue
// 			}
// 			if !isInMaze(point.x+i, point.y+j, len(maze), len(maze[0])) {
// 				continue
// 			}
// 			if maze[point.x+i][point.y+j] == '*' {
// 				res = append(res, &Point{point.x + i, point.y + j, point})
// 			}
// 		}
// 	}
// 	return res
// }

func isInMaze(x, y int, w, h int) bool {
	return x >= 0 && x < w &&
		y >= 0 && y < h
}

// Generates a maze using Prim's Algorithm
// https://en.wikipedia.org/wiki/Maze_generation_algorithm#Randomized_Prim.27s_algorithm
// func generateMaze(w, h int) [][]rune {
// 	maze := make([][]rune, w)
// 	for row := range maze {
// 		maze[row] = make([]rune, h)
// 		for ch := range maze[row] {
// 			maze[row][ch] = '*'
// 		}
// 	}
// 	rand.Seed(time.Now().UnixNano())
// 	point := &Point{x: rand.Intn(w), y: rand.Intn(h)}
// 	point2 := &Point{x: rand.Intn(w), y: rand.Intn(h)}
// 	maze[point.x][point.y] = '1'
// 	maze[point2.x][point2.y] = '2'
// 	var last *Point
// 	walls := adjacents(point, maze)
// 	for len(walls) > 0 {
// 		rand.Seed(time.Now().UnixNano())
// 		wall := walls[rand.Intn(len(walls))]
// 		for i, w := range walls {
// 			if w.x == wall.x && w.y == wall.y {
// 				walls = append(walls[:i], walls[i+1:]...)
// 				break
// 			}
// 		}
// 		opp := wall.Opposite()
// 		if isInMaze(opp.x, opp.y, w, h) && maze[opp.x][opp.y] == '*' {
// 			maze[wall.x][wall.y] = '.'
// 			maze[opp.x][opp.y] = '.'
// 			walls = append(walls, adjacents(opp, maze)...)
// 			last = opp
// 		}
// 	}
// 	maze[last.x][last.y] = 'G'
// 	bordered := make([][]rune, len(maze)+2)
// 	for r := range bordered {
// 		bordered[r] = make([]rune, len(maze[0])+2)
// 		for c := range bordered[r] {
// 			if r == 0 || r == len(maze)+1 || c == 0 || c == len(maze[0])+1 {
// 				bordered[r][c] = '*'
// 			} else {
// 				bordered[r][c] = maze[r-1][c-1]
// 			}
// 		}
// 	}
// 	return bordered
// }

/////////////////
// Termloop stuff
/////////////////

type Block struct {
	*tl.Rectangle
	px        int // Previous x
	py        int // Previous y
	move      bool
	g         *tl.Game
	w         int // Width of maze
	h         int // Height of maze
	score     int
	scoretext *tl.Text
	player    int
}

func NewBlock(x, y int, color tl.Attr, g *tl.Game, w, h, score int, scoretext *tl.Text, player int) *Block {
	b := &Block{
		g:         g,
		w:         w,
		h:         h,
		score:     score,
		scoretext: scoretext,
		player:    player,
	}
	b.Rectangle = tl.NewRectangle(x, y, 1, 1, color)
	return b
}

func (b *Block) Draw(s *tl.Screen) {
	if l, ok := b.g.Screen().Level().(*tl.BaseLevel); ok {
		// Set the level offset so the maze is always in the
		// center of the screen.
		sw, sh := s.Size()
		// x, y := b.Position())
		l.SetOffset(int(math.Round(float64(sw)*.05)), int(math.Round(float64(sh)*.1)))
	}

	b.Rectangle.Draw(s)
}

func (b *Block) Tick(ev tl.Event) {
	// Enable arrow key movement for player 1
	if b.player == 1 {
		if ev.Type == tl.EventKey {
			b.px, b.py = b.Position()
			switch ev.Key {
			case tl.KeyArrowRight:
				b.SetPosition(b.px+1, b.py)
			case tl.KeyArrowLeft:
				b.SetPosition(b.px-1, b.py)
			case tl.KeyArrowUp:
				b.SetPosition(b.px, b.py-1)
			case tl.KeyArrowDown:
				b.SetPosition(b.px, b.py+1)
			}
		}
	}

	if b.player == 2 {
		if ev.Type == tl.EventKey {
			b.px, b.py = b.Position()
			switch ev.Key {
			case tl.KeyCtrlD:
				b.SetPosition(b.px+1, b.py)
			case tl.KeyCtrlA:
				b.SetPosition(b.px-1, b.py)
			case tl.KeyCtrlW:
				b.SetPosition(b.px, b.py-1)
			case tl.KeyCtrlS:
				b.SetPosition(b.px, b.py+1)
			}
		}
	}
}

func (b *Block) Collide(c tl.Physical) {
	if r, ok := c.(*tl.Rectangle); ok {
		switch r.Color() {
		case wallColor:
			b.SetPosition(b.px, b.py)
		case goalColor:
			b.w += 1
			b.h += 1
			b.score += 1
			buildLevel(b.g, b.w, b.h, b.score)
		}
	}
}

func buildLevel(g *tl.Game, w, h, score int) {
	maze := level.New(w, h)
	l := tl.NewBaseLevel(tl.Cell{})
	g.Screen().SetLevel(l)
	g.Log("Building level with width %d and height %d", w, h)
	scoreTextPlayer1 := tl.NewText(0, 1, "Player 1 Score: "+strconv.Itoa(score),
		tl.ColorBlue, tl.ColorBlack)
	scoreTextPlayer2 := tl.NewText(0, 2, "Player 2 Score: "+strconv.Itoa(score),
		tl.ColorBlue, tl.ColorBlack)
	g.Screen().AddEntity(tl.NewText(0, 0, "Pyramid!", tl.ColorBlue, tl.ColorBlack))
	g.Screen().AddEntity(scoreTextPlayer1)
	g.Screen().AddEntity(scoreTextPlayer2)
	for i, row := range maze {
		for j, path := range row {
			switch path {
			case '*':
				l.AddEntity(tl.NewRectangle(i, j, 1, 1, wallColor))
			case '1':
				l.AddEntity(NewBlock(i, j, player1Color, g, w, h, score, scoreTextPlayer1, 1))
			case '2':
				l.AddEntity(NewBlock(i, j, player2Color, g, w, h, score, scoreTextPlayer2, 2))
			case 'G':
				l.AddEntity(tl.NewRectangle(i, j, 1, 1, goalColor))
			}
		}
	}
}

func main() {
	g := tl.NewGame()
	g.Screen().SetFps(60)
	buildLevel(g, 36, 12, 0)
	g.SetDebugOn(true)
	g.Start()
}
