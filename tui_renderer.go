package main

import (
	"github.com/gdamore/tcell"
	"fmt"
)

const deadCell = ' '
const aliveCell = '■'

type TuiRenderer struct {
	view Viewport
	screen tcell.Screen
}

func (renderer *TuiRenderer) Close() {
	renderer.screen.Fini()
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
			switch ev.Key() {
			case tcell.KeyCtrlC:
				close(quit)
				return

			case tcell.KeyRight:
				view.centerCol = min(view.centerCol + 1, b.cols - view.cols / 2)
				break

			case tcell.KeyLeft:
				view.centerCol = max(view.centerCol - 1, view.cols / 2)
				break

			case tcell.KeyUp:
				view.centerRow = max(view.centerRow - 1, view.rows / 2)
				break

			case tcell.KeyDown:
				view.centerRow = min(view.centerRow + 1, b.rows - view.rows / 2)
				break

			case tcell.KeyRune:
				if ev.Rune() == '-' {
					view.zoom++
				} else if ev.Rune() == '+' {
					view.zoom = max(view.zoom - 1, 1)
				}
			}
		}
	}
}

func puts(s tcell.Screen, row, col int, str string) {
	for _, c := range str {
		s.SetContent(col, row, rune(c), []rune(""), tcell.StyleDefault)
		col++
	}
}

func (renderer TuiRenderer) Render(b *Board, info RenderInfo) {
	view := renderer.view

	var rowFrom = max(0, view.centerRow - view.rows / 2)
	var colFrom = max(0, view.centerCol - view.cols / 2)

	renderer.screen.Clear()
	for row := 0; row < view.rows; row++ {
		for col := 0; col < view.cols; col++ {
			cell := getCellState(view, b, rowFrom + row * view.zoom, colFrom + col * view.zoom)

			renderer.screen.SetContent(
				col,
				row,
				rune(cell),
				[]rune(""),
				tcell.StyleDefault,
			)
		}
	}

	puts(renderer.screen, view.rows - 4, 0, fmt.Sprintf("%d alive", info.alive))
	puts(renderer.screen, view.rows - 3, 0, fmt.Sprintf("%dx%d, %d", b.rows, b.cols, b.generation))
	puts(renderer.screen, view.rows - 2, 0, fmt.Sprintf("[%d, %d], %dx", view.centerRow, view.centerCol, view.zoom))
	puts(renderer.screen, view.rows - 1, 0, fmt.Sprintf("%s", info.elapsed))
	renderer.screen.Show()
}

func getCellState(view Viewport, b *Board, globalRow, globalCol int) int32 {
	if view.zoom == 1 {
		if b.Get(globalRow, globalCol) == ALIVE {
			return aliveCell
		}
	} else {
		count := b.AliveNeighbours(globalRow, globalCol, view.zoom - 1)
		if count > view.zoom * view.zoom / 2 {
			return aliveCell
		}
	}
	return deadCell
}