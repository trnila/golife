package main

import (
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

func CreateRenderer(b *Board) Renderer {
	var renderer Renderer

	renderer, err := NewTuiRenderer(b)
	if err == nil {
		return renderer
	}
	fmt.Print(err)

	// fallback to ascii renderer
	renderer, err = NewAsciiRenderer(b); if err != nil {
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
	nextRender := make(chan RenderInfo)

	renderer := CreateRenderer(b)
	defer renderer.Close()

	go compute(threads, b, nextRender)

	renderer.Start(nextRender)
}


func compute(threads int, b *Board, nextRender chan RenderInfo) {
	const delay = 100 * time.Millisecond

	alivesChan := make(chan int)
	for {
		start := time.Now()

		for t := 0; t < threads; t++ {
			go func(t int) {
				var from = (b.rows + threads) / threads * t
				var to = from + (b.rows + threads) / threads
				if to > b.rows {
					to = b.rows
				}

				calc(b, from, to, alivesChan)
			}(t)
		}

		totalAlives := 0
		for t := 0; t < threads; t++ {
			totalAlives += <- alivesChan
		}

		b.generation++
		elapsed := time.Since(start)
		nextRender <- RenderInfo{
			elapsed: elapsed,
			alive: totalAlives,
		}

		time.Sleep(delay - elapsed)
	}
}