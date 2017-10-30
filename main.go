package main

import (
	"sync"
	"time"
	"os"
	"fmt"
	"math/rand"
	"runtime"
	"flag"
)

const defaultSize = 100

func createBoard(rows, cols uint) *Board {
	var b *Board

	if len(flag.Args()) == 0 {
		if rows <= 0 {
			rows = defaultSize
		}
		if cols <= 0 {
			cols = defaultSize
		}

		b = NewBoard(int(rows), int(cols))

		rand.Seed(time.Now().Unix())
		for i := 0; i < b.rows * b.cols / 2; i++ {
			b.Init(rand.Intn(b.rows), rand.Intn(b.cols), Cell(rand.Intn(2)))
		}
	} else {
		filename := flag.Arg(0)

		if filename == "-" {
			b = loadFile(os.Stdin, rows, cols)
		} else {
			b = load(filename, rows, cols)
		}
	}

	return b
}

func CreateRenderer(b *Board, quit chan struct{}) Renderer {
	var renderer Renderer

	renderer, err := NewTuiRenderer(b, quit)
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
	threads := runtime.GOMAXPROCS(0)
	rows := flag.Uint("rows", 0, "number of rows in world")
	cols := flag.Uint("cols", 0, "number of cols in world")

	flag.Parse()

	b := createBoard(*rows, *cols)
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

				var from = (b.rows + threads) / threads * t
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