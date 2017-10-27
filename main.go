package main

import (
	"fmt"
	"sync"
	_ "time"
	"math"
	"github.com/gdamore/tcell"
	"os"
	"math/rand"
)

const DEAD = 0;
const ALIVE = 1;
type Viewport struct {
	rows int
	cols int
	centerRow int
	centerCol int
	zoom int
	screen tcell.Screen
}

type Cell int

type Board struct {
	rows int
	cols int
	data [][]Cell
	generation int
}

func NewBoard(size int) *Board {
	const generations = 2

	var board = Board{
		rows: size,
		cols: size,
		data: make([][]Cell, generations),
	}

	for i := 0; i < generations; i++ {
		board.data[i] = make([]Cell, size * size);
	}

	return &board
}

func (b Board) isValid(row, col int) bool {
	return row >= 0 && col >= 0 &&
		row < b.rows && col < b.cols;
}

func (b Board) getIndex(row, col int) int {
	return row * b.cols + col
}

func (b Board) Get(row, col int) Cell {
	if !b.isValid(row, col) {
		return DEAD
	}

	return b.data[b.generation % 2][b.getIndex(row, col)];
}

func (b Board) Set(row, col int, val Cell) {
	//if b.isValid(row, col) {
		b.data[(b.generation + 1) % 2][b.getIndex(row, col)] = val
	//}
}

func (b Board) Init(row, col int, val Cell) {
	b.data[0][b.getIndex(row, col)] = val
	b.data[1][b.getIndex(row, col)] = val
}

func (b Board) AliveNeighbours(row, col, size int) int {
	var count = 0;

	for r := -size; r <= size; r++ {
		for c := -size; c <= size; c++ {
			if b.Get(row + r, col + c) == ALIVE {
				count++
			}
		}
	}

	return count - int(b.Get(row, col))
}

func (b Board) Print(view *Viewport) {
	var rowFrom = int(math.Max(0, float64(view.centerRow - view.rows / 2)))
	var colFrom = int(math.Max(0, float64(view.centerCol - view.cols / 2)))

	view.screen.Clear()
	var z = view.zoom
	for r := 0; r < view.rows; r += 1 {
		for c := 0; c < view.cols; c += 1 {
			a := ' ';

			if view.zoom == 1 && b.Get(colFrom + c * z, rowFrom + r * z) == ALIVE {
				a = '*'
			} else {
				count := b.AliveNeighbours(colFrom+c*z, rowFrom+r*z, view.zoom-1)
				if count > view.zoom*view.zoom/2 || (view.zoom == 1 && count == 1) {
					a = '*'
				}
			}

			view.screen.SetContent(
				c,
				r,
				rune(a),
				[]rune(""),
				tcell.StyleDefault,
			)
		}
	}
	view.screen.Show()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func calc(b *Board, start, to int) {
	for r := start; r < to; r++ {
		for c := 0; c < b.cols; c++ {
			var aliveNeighbours = b.AliveNeighbours(r, c, 1)
			var state = b.Get(r, c)
			if state == ALIVE {
				if aliveNeighbours < 2 || aliveNeighbours > 3 {
					state = DEAD
				} else {
					state = ALIVE
				}
			} else if aliveNeighbours == 3 {
				state = ALIVE
			}

			b.Set(r, c, state)
		}
	}
}

func main() {
	const size = 1000;
	const threads = 4;
	var b = NewBoard(size)


	b.generation = 0
	b.Init(5, 5, ALIVE)
	b.Init(6, 5, ALIVE)
	b.Init(7, 5, ALIVE)

	b.Init(1, 2, ALIVE)
	b.Init(1, 3, ALIVE)
	b.Init(2, 2, ALIVE)
	b.Init(2, 3, ALIVE)

	b.Init(20, 5, ALIVE)
	b.Init(21, 5, ALIVE)
	b.Init(22, 5, ALIVE)

	b.Init(20, 170, ALIVE)
	b.Init(21, 170, ALIVE)
	b.Init(22, 170, ALIVE)

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}


	cols, rows := s.Size()
	var view = Viewport{
		rows: rows,
		cols: cols,
		centerRow: rows / 2,
		centerCol: cols / 2,
		zoom: 1,
	}
	view.screen = s

	for i := 0; i < 1000000; i++ {
		b.Init(rand.Intn(size), rand.Intn(size), Cell(rand.Intn(2)))
	}


	quit := make(chan struct{})



	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC {
					close(quit)
					return
				} else if ev.Key() == tcell.KeyRight {
					view.centerCol += 1
				} else if ev.Key() == tcell.KeyLeft {
					view.centerCol -= 1
				} else if ev.Key() == tcell.KeyUp {
					view.centerRow -= 1
				} else if ev.Key() == tcell.KeyDown {
					view.centerRow += 1
				} else if ev.Key() == tcell.KeyRune {
					if ev.Rune() == '-' {
						view.zoom += 1
					} else if ev.Rune() == '+' {
						view.zoom = max(view.zoom - 1, 1);
					}
				}
			}
		}
	}()


	go compute(threads, size, b, &view)

	<- quit
	s.Fini()
}


func compute(threads int, size int, b *Board, view *Viewport) {
	for {
		//start := time.Now()

		var wait sync.WaitGroup
		wait.Add(threads)

		for t := 0; t < threads; t++ {
			go func(t int) {
				defer wait.Done()

				var from = (size + threads) / threads * t;
				var to = from + (size+threads)/threads
				if to > size {
					to = size
				}

				calc(b, from, to)
			}(t)
		}

		wait.Wait()

		b.generation++
		b.Print(view)

		//elapsed := time.Since(start)
		//fmt.Println(elapsed)
	}
}