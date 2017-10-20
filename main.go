package main

import (
	"fmt"
	"time"
)

const DEAD = 0;
const ALIVE = 1;

type Board struct {
	rows int
	cols int
	data [][]int
	generation int
}

func NewBoard(size int) *Board {
	const generations = 2

	var board = Board{
		rows: size,
		cols: size,
		data: make([][]int, generations),
	}

	for i := 0; i < generations; i++ {
		board.data[i] = make([]int, size * size);
	}

	return &board
}

func (b Board) isValid(row, col int) bool {
	return row >= 0 && col >= 0 &&
		row < b.rows && col < b.cols;
}

func (b Board) getIndex(row, col int) int {
	return row * b.cols + col
}

func (b Board) Get(row, col int) int {
	if !b.isValid(row, col) {
		return DEAD
	}

	return b.data[b.generation % 2][b.getIndex(row, col)];
}

func (b Board) Set(row, col, val int) {
	//if b.isValid(row, col) {
		b.data[(b.generation + 1) % 2][b.getIndex(row, col)] = val
	//}
}

func (b Board) Init(row, col, val int) {
	b.data[0][b.getIndex(row, col)] = val
	b.data[1][b.getIndex(row, col)] = val
}

func (b Board) AliveNeighbours(row, col int) int {
	var count = 0;

	for r := -1; r <= 1; r++ {
		for c := -1; c <= 1; c++ {
			if b.Get(row + r, col + c) == ALIVE {
				count++
			}
		}
	}

	return count - b.Get(row, col)
}

func (b Board) Print() {
	for r := 0; r < b.rows; r++ {
		for c := 0; c < b.cols; c++ {
			if b.Get(r, c) == ALIVE {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func main() {
	const size = 10;
	var b = NewBoard(size)


	b.generation = 0
	b.Init(5, 5, ALIVE)
	b.Init(6, 5, ALIVE)
	b.Init(7, 5, ALIVE)

	b.Init(1, 2, ALIVE)
	b.Init(1, 3, ALIVE)
	b.Init(2, 2, ALIVE)
	b.Init(2, 3, ALIVE)

	//b.AliveNeighbours(6, 4)

	for {
		fmt.Print("\033[H\033[2J")
		for r := 0; r < size; r++ {
			for c := 0; c < size; c++ {
				var aliveNeighbours = b.AliveNeighbours(r, c)
				if b.Get(r, c) == ALIVE {
					if aliveNeighbours < 2 || aliveNeighbours > 3 {
						b.Set(r, c, DEAD)
					} else {
						b.Set(r, c, ALIVE)
					}
				} else if aliveNeighbours == 3 {
					b.Set(r, c, ALIVE);
				}
			}
		}
		b.generation++
		b.Print()
		time.Sleep(500 * time.Millisecond)
	}
}