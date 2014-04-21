package main

import "github.com/nsf/termbox-go"
import "time"
import "math/rand"

type GCell struct {
	Mark      rune
	X         int
	Y         int
	Fg        termbox.Attribute
	Bg        termbox.Attribute
	Breakable bool
	Passable  bool
}

type Painter interface {
	draw()
}

type Player struct {
	Point GCell
}

func (p *Player) draw() {
	point := p.Point
	termbox.SetCell(point.X, point.Y, point.Mark, point.Fg, point.Bg)
}

type Map struct {
	Points []GCell
	Xmax   int
	Ymax   int
}

func (m *Map) canMove(x int, y int, ox int, oy int) (int, int) {
	if( m.pointAt(x, y).Passable ) {
		return x, y
	} else {
		return ox, oy
	}
}

func (m *Map) draw() {
	for _, cell := range m.Points {
		termbox.SetCell(cell.X, cell.Y, cell.Mark, cell.Fg, cell.Bg)
	}
}

func (m *Map) reset(x int, y int) {
	op := m.Points[(y * m.Xmax) + x]
	termbox.SetCell(x, y, op.Mark, op.Fg, op.Bg)
}

func (m *Map) pointAt(x int, y int) GCell {
	return m.Points[m.index(x, y)]
}

func (m *Map) index(x int, y int) int {
	return (y * m.Xmax) + x
}

func (m *Map) paint(x int, y int, mark rune, fg termbox.Attribute, bg termbox.Attribute, pass bool, brk bool) {
	point := GCell{mark, x, y, fg, bg, pass, brk}
	m.Points[m.index(x, y)] = point
}

func (m *Map) fillRandom() {
	wallmark := '#'
	openmark := '.'
	fg, bg := termbox.ColorDefault, termbox.ColorDefault
	wfg, wbg := termbox.ColorGreen, termbox.ColorDefault

	for x := 0; x < m.Xmax; x++ {
		for y := 0; y < m.Ymax; y++ {
			if( rand.Intn(100) > 30 ){
				m.paint(x, y, openmark, fg, bg, false, true)
			} else {
				m.paint(x, y, wallmark, wfg, wbg, true, false)
			}
		}
	}
	m.fillBorders()
}

func (m *Map) fillMap() {
	mark := '.'
	fg, bg := termbox.ColorDefault, termbox.ColorDefault

	for x := 0; x < m.Xmax; x++ {
		for y := 0; y < m.Ymax; y++ {
			m.paint(x, y, mark, fg, bg, false, true)
		}
	}
	m.fillBorders()
}

func (m *Map) fillBorders() {
	mark := '#'
	fg, bg := termbox.ColorGreen, termbox.ColorDefault

	for x := 0; x < m.Xmax; x++ {
		m.paint(x, m.Ymax-1, mark, fg, bg, false, false)
		m.paint(x, 0       , mark, fg, bg, false, false)
	}
	for y := 0; y < m.Ymax; y++ {
		m.paint(m.Xmax-1, y, mark, fg, bg, false, false)
		m.paint(0       , y, mark, fg, bg, false, false)
	}
}

func flush(painters []Painter) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, p := range painters {
		p.draw()
	}
	termbox.Flush()
}

func flushSingle(painter Painter){
	painter.draw()
	termbox.Flush()
}


func main() {
	x, y := 5, 5

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	xmax, ymax := termbox.Size()
	points := make([]GCell, xmax * ymax)
	world := Map{points, xmax, ymax}

	point := GCell{'@', 5, 5, termbox.ColorGreen, termbox.ColorDefault, true, true}
	player := Player{point}

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	world.fillRandom()

	painters := []Painter{&world, &player}
	flush(painters)
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
			if ev.Type == termbox.EventKey {
				ox, oy := x, y
				switch ev.Ch {
				case 'j': y++
				case 'k': y--
				case 'l': x++
				case 'h': x--
				}
				x, y = world.canMove(x, y, ox, oy)
				world.reset(ox,oy)
			}
			player.Point.X = x
			player.Point.Y = y
		}
		flushSingle(&player)
		time.Sleep(10*time.Millisecond)
	}
}
