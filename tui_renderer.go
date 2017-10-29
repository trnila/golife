package main

import (
	"github.com/gdamore/tcell"
	"fmt"
	"math"
	"time"
)

type TuiRenderer struct {
	view Viewport
	screen tcell.Screen
}

func (r *TuiRenderer) Close() {
	r.screen.Fini()
}

func NewTuiRenderer(b* Board, quit chan struct{}) (*TuiRenderer, error) {
	s, e := tcell.NewScreen()
	if e != nil {
		return nil, e
	}
	if e = s.Init(); e != nil {
		return nil, e
	}

	cols, rows := s.Size()
	var view = NewViewport(rows, cols)

	renderer := &TuiRenderer{
		view: view,
		screen: s,
	}

	go handleEvents(renderer, b, quit)

	return renderer, nil
}

func handleEvents(renderer *TuiRenderer, b *Board, quit chan struct{}) {
	view := &renderer.view
	for {
		ev := renderer.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				close(quit)
				return
			} else if ev.Key() == tcell.KeyRight {
				view.centerCol = min(view.centerCol + 1, b.cols - view.cols / 2)
			} else if ev.Key() == tcell.KeyLeft {
				view.centerCol = max(view.centerCol - 1, view.cols / 2)
			} else if ev.Key() == tcell.KeyUp {
				view.centerRow = max(view.centerRow - 1, view.rows / 2)
			} else if ev.Key() == tcell.KeyDown {
				view.centerRow = min(view.centerRow + 1, b.rows - view.rows / 2)
			} else if ev.Key() == tcell.KeyRune {
				if ev.Rune() == '-' {
					view.zoom += 1
				} else if ev.Rune() == '+' {
					view.zoom = max(view.zoom - 1, 1);
				}
			}
		}
	}
}

func print(s tcell.Screen, row, col int, str string) {
	for _, c := range(str) {
		s.SetContent(col, row, rune(c), []rune(""), tcell.StyleDefault)
		col += 1
	}
}

func (renderer TuiRenderer) Render(b *Board, elapsed time.Duration) {
	view := renderer.view

	var rowFrom = int(math.Max(0, float64(view.centerRow - view.rows / 2)))
	var colFrom = int(math.Max(0, float64(view.centerCol - view.cols / 2)))

	renderer.screen.Clear()
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

			renderer.screen.SetContent(
				c,
				r,
				rune(a),
				[]rune(""),
				tcell.StyleDefault,
			)
		}
	}

	print(renderer.screen, view.rows - 3, 0, fmt.Sprintf("%dx%d, %d", b.rows, b.cols, b.generation))
	print(renderer.screen, view.rows - 2, 0, fmt.Sprintf("[%d, %d], %dx", view.centerRow, view.centerCol, view.zoom))
	print(renderer.screen, view.rows - 1, 0, fmt.Sprintf("%s", elapsed))
	renderer.screen.Show()
}