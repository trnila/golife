package main

import "time"

type Renderer interface {
	Render(b *Board, elapsed time.Duration)
	Close()
}

func NewViewport(rows, cols int) Viewport  {
	return Viewport{
		rows: rows,
		cols: cols,
		centerRow: rows / 2,
		centerCol: cols / 2,
		zoom: 1,
	}
}
