package main

import (
	"bufio"
	"strings"
	"fmt"
	"strconv"
	"os"
)

func load(filename string) *Board {
	f, err := os.Open(filename); if err != nil {
		panic(err)
	}
	defer f.Close()

	return loadFile(f)
}

func loadFile(f *os.File) *Board {
	reader := bufio.NewReader(f)

	var board *Board
	row := 0
	first := true
	for {
		line, err := reader.ReadString('\n'); if err != nil {
			break
		}

		if first {
			first = false
			size, err := strconv.ParseInt(strings.Trim(line, " \n\r"), 10, 32); if err != nil {
				panic(err)
			}

			board = NewBoard(int(size))
		} else {
			for col, c := range(strings.Trim(line, "\n\r")) {
				if col >= board.rows {
					panic(fmt.Sprintf("Invalid column %d", c))
				}

				if row >= board.rows {
					panic(fmt.Sprintf("Invalid column %d", row))
				}

				if c != '.' {
					board.Init(col, row, ALIVE)
				}
			}
		}

		row += 1
	}

	return board
}