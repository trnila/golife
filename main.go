package main

import (
	"sync"
	"time"
	"os"
	"fmt"
)

func createBoard() *Board {
	var b *Board

	if len(os.Args) == 1 {
		b = NewBoard(1000)

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
	} else {
		b = load(os.Args[1])
	}

	return b
}

func CreateRenderer(b *Board, quit chan struct{}) Renderer {
	var renderer Renderer

	renderer, err := NewTuiRenderer(b, quit);
	if err == nil {
		return renderer
	}
	fmt.Print(err)

	// fallback to ascii renderer
	renderer, err = NewAsciiRenderer(b, quit); if err != nil {
		panic(err)
	}

	return renderer
}

func main() {
	const threads = 4;

	b := createBoard()
	quit := make(chan struct{})

	renderer := CreateRenderer(b, quit)
	defer renderer.Close()

	go compute(threads, b, renderer)

	<- quit
}


func compute(threads int, b *Board, renderer Renderer) {
	const delay = 100 * time.Millisecond

	for {
		start := time.Now()

		var wait sync.WaitGroup
		wait.Add(threads)

		for t := 0; t < threads; t++ {
			go func(t int) {
				defer wait.Done()

				var from = (b.rows + threads) / threads * t;
				var to = from + (b.rows + threads) / threads
				if to > b.rows {
					to = b.rows
				}

				calc(b, from, to)
			}(t)
		}

		wait.Wait()

		b.generation++
		elapsed := time.Since(start)
		renderer.Render(b, elapsed)

		time.Sleep(delay - elapsed)
	}
}