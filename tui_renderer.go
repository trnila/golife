package main

import (
	"github.com/gdamore/tcell"
	"fmt"
)

const deadCell = ' '
const aliveCell = '■'

const (
	cmdRight = 1
	cmdLeft = 2
	cmdUp = 3
	cmdDown = 4
	cmdZoomIn = 5
	cmdZoomOut = 6
	cmdQuit = 7
)

type TuiRenderer struct {
	board *Board
	view Viewport
	screen tcell.Screen
	cmd chan int
}

func (renderer *TuiRenderer) Close() {
	renderer.screen.Fini()
}

func NewTuiRenderer(b* Board) (*TuiRenderer, error) {
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
		board: b,
		view: view,
		screen: s,
		cmd: make(chan int),
	}

	go handleEvents(renderer, b)

	return renderer, nil
}

func (renderer *TuiRenderer) Start(renderFrame chan RenderInfo) {
	var state RenderInfo
	for {
		select {
			case state = <- renderFrame:
			case cmd := <- renderer.cmd:
				view := &renderer.view
				switch cmd {
				case cmdRight:
					view.centerCol = min(view.centerCol + 1, renderer.board.cols - view.cols / 2)
					break

				case cmdLeft:
					view.centerCol = max(view.centerCol - 1, view.cols / 2)
					break

				case cmdUp:
					view.centerRow = max(view.centerRow - 1, view.rows / 2)
					break

				case cmdDown:
					view.centerRow = min(view.centerRow + 1, renderer.board.rows - view.rows / 2)
					break

				case cmdZoomIn:
					view.zoom = max(view.zoom - 1, 1)
					break

				case cmdZoomOut:
					view.zoom++
					break

				case cmdQuit:
					return
				}
		}

		renderer.Render(renderer.board, state)
	}
}

func key(k tcell.Key, c byte) int {
	return (int(k) << 16) | int(c)
}

func handleEvents(renderer *TuiRenderer, b *Board) {
	keys := map[int]int {
		key(tcell.KeyCtrlC, 3): cmdQuit,
		key(tcell.KeyRight, 0): cmdRight,
		key(tcell.KeyLeft, 0): cmdLeft,
		key(tcell.KeyUp, 0): cmdUp,
		key(tcell.KeyDown, 0): cmdDown,
		key(tcell.KeyRune, '+'): cmdZoomIn,
		key(tcell.KeyRune, '-'): cmdZoomOut,
	}

	for {
		ev := renderer.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			k := key(ev.Key(), byte(ev.Rune()))
			action, ok := keys[k]
			if ok {
				renderer.cmd <- action
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

func (renderer* TuiRenderer) Render(b *Board, info RenderInfo) {
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