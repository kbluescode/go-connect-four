package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	gridSize    int  = 5
	defaultMark rune = ' '
	p1Mark      rune = 'x'
	p2Mark      rune = 'o'
)

var (
	grid    Grid
	scanner = bufio.NewScanner(os.Stdin)
	turn    = 1
	maxTurn = gridSize * gridSize
)

// Cell - A single cell in the grid
type Cell struct {
	row, col int
	val      rune
}

// String - to satisfy the Stringer interface for fmt
func (c *Cell) String() string {
	return fmt.Sprintf("[%c]", c.val)
}

// Grid - a grid of cells
type Grid [gridSize][gridSize]*Cell

// Cell - retrieves a cell from grid or returns an error
func (g *Grid) Cell(row, col int) (cell *Cell, err error) {
	if (row < 0 || row >= gridSize) || (col < 0 || col >= gridSize) {
		err = fmt.Errorf("Invalid row or column given, must be between 0 and %v", gridSize)
	}
	if err == nil {
		cell = grid[row][col]
	}
	return cell, err
}

func main() {
	resetGrid()
	gameLoop()
}

// gameLoop - the core game loop, ending if a player wins
func gameLoop() {
	var won bool
	var winnerMark rune
	currentMark := p1Mark
	for ; !won && turn <= maxTurn; won, winnerMark = hasWon() {
		fmt.Printf("Current player: '%c'\n\n", currentMark)
		displayGrid()
		takeTurn(currentMark)
		if currentMark == p1Mark {
			currentMark = p2Mark
		} else {
			currentMark = p1Mark
		}
		turn++
	}
	if winnerMark == defaultMark {
		fmt.Println("\n\nGame Over: It's a draw!")
	} else {
		fmt.Printf("\n\nGame Over:\nPlayer %q has won!\n", winnerMark)
	}
	displayGrid()
}

// resetGrid - reverts the grid to its base state of blank cells
func resetGrid() {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			grid[row][col] = &Cell{row: row, col: col, val: defaultMark}
		}
	}
}

// displayGrid - prints out the grid in a human-readable format
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

// takeTurn - will attempt to get input from the user and display errors if the input is invalid
func takeTurn(mark rune) {
	for {
		if scanner.Scan() {
			text := scanner.Text()
			if commandChunks, err := parseCommand(text); err == nil {
				cell, err := grid.Cell(commandChunks[0], commandChunks[1])
				if err != nil {
					fmt.Println(err)
				} else if cell.val != defaultMark {
					fmt.Println("Invalid move: choose a space that hasn't been chosen")
				} else {
					cell.val = mark
					break
				}
			}
		}
	}
}

// parseCommand - takes a string and tries to produce a game command out of it. Produces valid coordinates if string passes, otherwise returns an error.
// ex. "0 1"
func parseCommand(command string) ([]int, error) {
	var chunks []int

	if !strings.ContainsRune(command, ' ') {
		return chunks, errors.New("Command must contain a ' ' between row and column")
	}

	stringChunks := strings.Split(command, " ")
	for _, val := range stringChunks {
		num, err := strconv.Atoi(val)
		if err != nil {
			return chunks, err
		}
		chunks = append(chunks, num)
	}

	return chunks, nil
}

// hasWon - returns true if any line of 4 of 'x' or 'o' has been made in the grid
func hasWon() (bool, rune) {
	if fourConnected(p1Mark) {
		return true, p1Mark
	} else if fourConnected(p2Mark) {
		return true, p2Mark
	}

	return false, defaultMark
}

// fourConnected - returns true if collection of cells has a contiguous line of 4 or false if no contiguous line is present
func fourConnected(symbol rune) bool {
	// check horizontal
	for row := 0; row < gridSize; row++ {
		contiguousCount := 0
		for col := 0; col < gridSize && contiguousCount < 4; col++ {
			cell := grid[row][col]
			if cell.val == symbol {
				contiguousCount++
			} else {
				contiguousCount = 0
			}
		}
		if contiguousCount == 4 {
			return true
		}
	}

	// check vertical
	for col := 0; col < gridSize; col++ {
		contiguousCount := 0
		for row := 0; row < gridSize && contiguousCount < 4; row++ {
			cell := grid[row][col]
			if cell.val == symbol {
				contiguousCount++
			} else {
				contiguousCount = 0
			}
		}
		if contiguousCount == 4 {
			return true
		}
	}

	//check down-right angles
	startingCells := make([]*Cell, 0)
	for col := 0; col < gridSize; col++ {
		startingCells = append(startingCells, grid[0][col])
	}
	for row := 0; row < gridSize; row++ {
		startingCells = append(startingCells, grid[row][0])
	}
	for _, start := range startingCells {
		contiguousCount := 0
		for row, col := start.row, start.col; row < gridSize && col < gridSize && contiguousCount < 4; row, col = row+1, col+1 {
			cell := grid[row][col]
			if cell.val == symbol {
				contiguousCount++
			} else {
				contiguousCount = 0
			}
		}
		if contiguousCount == 4 {
			return true
		}
	}

	//check down-left angles
	startingCells = make([]*Cell, 0)
	for col := gridSize - 1; col <= 0; col-- {
		startingCells = append(startingCells, grid[gridSize-1][col])
	}
	for row := gridSize - 1; row <= 0; row-- {
		startingCells = append(startingCells, grid[row][gridSize-1])
	}
	for _, start := range startingCells {
		contiguousCount := 0
		for row, col := start.row, start.col; row >= 0 && col >= 0 && contiguousCount < 4; row, col = row-1, col-1 {
			cell := grid[row][col]
			if cell.val == symbol {
				contiguousCount++
			} else {
				contiguousCount = 0
			}
		}
		if contiguousCount == 4 {
			return true
		}
	}

	return false
}
