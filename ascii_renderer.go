package main

import (
	"fmt"
)

type AsciiRenderer struct {
	board *Board
	renderFrame chan RenderInfo
}

func (renderer AsciiRenderer) Start(renderFrame chan RenderInfo) {
	b := renderer.board

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

func NewAsciiRenderer(b* Board) (AsciiRenderer, error) {
	return AsciiRenderer{
		board:b,
	}, nil
}

