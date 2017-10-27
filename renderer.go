package main

import "time"

type Renderer interface {
	Render(b *Board, elapsed time.Duration)
}
