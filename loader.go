package main

import (
	"bufio"
	"strings"
	"fmt"
	"strconv"
	"os"
	"log"
)

func load(filename string, rows, cols uint) *Board {
	f, err := os.Open(filename); if err != nil {
		panic(err)
	}
	defer f.Close()

	return loadFile(f, rows, cols)
}

func loadFile(f *os.File, rows, cols uint) *Board {
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

			effRows := uint(size)
			effCols := uint(size)

			if rows > 0 {
				effRows = rows
			}

			if cols > 0 {
				effCols = cols
			}

			board = NewBoard(int(effRows), int(effCols))
		} else {
			for col, c := range strings.Trim(line, "\n\r") {
				if col >= board.cols {
					log.Print(fmt.Sprintf("Invalid column %d", c))
				}

				if row >= board.rows {
					log.Print(fmt.Sprintf("Invalid column %d", row))
				}

				if c != '.' {
					board.Init(row, col, ALIVE)	
				}
			}
		}

		row++
	}

	return board
}