package main

const DEAD = 0;
const ALIVE = 1;

type Cell int64
const BITS = 64

type Board struct {
	rows int
	cols int
	data [][]Cell
	generation int
}

func NewBoard(size int) *Board {
	const generations = 2

	var board = Board{
		rows: size,
		cols: size,
		data: make([][]Cell, generations),
	}

	for i := 0; i < generations; i++ {
		board.data[i] = make([]Cell, (size + BITS) / BITS * size);
	}

	return &board
}

func (b Board) isValid(row, col int) bool {
	return row >= 0 && col >= 0 &&
		row < b.rows && col < b.cols;
}

func (b Board) getIndex(row, col int) int {
	return row * b.cols / BITS + col / BITS
}

func (b Board) Get(row, col int) Cell {
	if !b.isValid(row, col) {
		return DEAD
	}

	val := b.data[b.generation % 2][b.getIndex(row, col)] >> (uint(col) % BITS);

	val &= 0x1;
	return Cell(val)
}

func (b Board) Set(row, col int, val Cell) {
	if val == ALIVE {
		b.data[(b.generation+1)%2][b.getIndex(row, col)] |= 1 << (uint(col) % BITS)
	} else {
		b.data[(b.generation+1)%2][b.getIndex(row, col)] &^= (1 << (uint(col) % BITS))
	}
}

func (b Board) Init(row, col int, val Cell) {
	b.data[0][b.getIndex(row, col)] |= val << (uint(col) % BITS)
	b.data[1][b.getIndex(row, col)] |= val << (uint(col) % BITS)
}

func (b Board) AliveNeighbours(row, col, size int) int {
	var count = 0;

	for r := -size; r <= size; r++ {
		for c := -size; c <= size; c++ {
			if b.Get(row + r, col + c) == ALIVE {
				count++
			}
		}
	}

	return count - int(b.Get(row, col))
}

func calc(b *Board, start, to int) {
	for r := start; r < to; r++ {
		for c := 0; c < b.cols; c++ {
			var aliveNeighbours = b.AliveNeighbours(r, c, 1)
			var state = b.Get(r, c)
			if state == ALIVE {
				if aliveNeighbours < 2 || aliveNeighbours > 3 {
					state = DEAD
				} else {
					state = ALIVE
				}
			} else if aliveNeighbours == 3 {
				state = ALIVE
			}

			b.Set(r, c, state)
		}
	}
}