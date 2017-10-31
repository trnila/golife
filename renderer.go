package main

import "time"

type RenderInfo struct {
	elapsed time.Duration
	alive int
}

type Renderer interface {
	Start(renderFrame chan RenderInfo)
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
