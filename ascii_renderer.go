package main

import (
	"time"
	"fmt"
)

type AsciiRenderer struct {
	view Viewport
}

type Viewport struct {
	rows int
	cols int
	centerRow int
	centerCol int
	zoom int
}

func (renderer AsciiRenderer) Render(b *Board, elapsed time.Duration) {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			c := " "
			if b.Get(row, col) == ALIVE {
				c = "*"
			}
			fmt.Print(c)
		}
		fmt.Print("\n")
	}
}

func (renderer AsciiRenderer) Close() {

}

func NewAsciiRenderer(b* Board, quit chan struct{}) (AsciiRenderer, error) {
	return AsciiRenderer{
		view: NewViewport(10, 10),
	}, nil
}

