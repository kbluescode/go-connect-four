package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	gridSize    int  = 5
	turnMax     int  = 5
	defaultMark rune = ' '
)

var (
	grid    [gridSize][gridSize]*Cell
	scanner = bufio.NewScanner(os.Stdin)
)

// Cell - A single cell in the grid
type Cell struct {
	row, col int
	val      rune
}

func (c *Cell) String() string {
	return fmt.Sprintf("[%c]", c.val)
}

func main() {
	resetGrid()
	gameLoop()
}

func gameLoop() {
	currentMark := 'x'
	for turn := 0; turn < turnMax; turn++ {
		fmt.Printf("Current player: '%c'\n\n", currentMark)
		displayGrid()
		takeTurn(currentMark)
		if currentMark == 'x' {
			currentMark = 'o'
		} else {
			currentMark = 'x'
		}
	}
}

func resetGrid() {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			grid[row][col] = &Cell{row: row, col: col, val: defaultMark}
		}
	}
}

func displayGrid() {
	buffer := bytes.NewBufferString("   0  1  2  3  4\n")
	for row := 0; row < gridSize; row++ {
		buffer.WriteString(fmt.Sprintf("%d ", row))
		for col := 0; col < gridSize; col++ {
			buffer.WriteString(grid[row][col].String())
		}
		buffer.WriteRune('\n')
	}
	fmt.Println(buffer.String())
}

func takeTurn(mark rune) {
	if scanner.Scan() {
		text := scanner.Text()
		if strings.ContainsRune(text, ' ') {
			chunks := strings.Split(text, " ")
			if len(chunks) > 1 {
				row, err1 := strconv.Atoi(chunks[0])
				col, err2 := strconv.Atoi(chunks[1])
				if err1 == nil && err2 == nil &&
					row >= 0 && row < gridSize &&
					col >= 0 && col < gridSize {
					cell := grid[row][col]
					if cell.val == defaultMark {
						cell.val = mark
					}
				}
			}
		}
	}
}
