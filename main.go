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
func (c Cell) String() string {
	return fmt.Sprintf("[%c]", c.val)
}

// Grid - a grid of cells
type Grid [gridSize][gridSize]*Cell

// Cell - retrieves a cell from grid or returns an error
func (g Grid) Cell(row, col int) (cell *Cell, err error) {
	if (row < 0 || row >= gridSize) || (col < 0 || col >= gridSize) {
		err = fmt.Errorf("Invalid row or column given, must be between 0 and %v", gridSize)
	}
	if err == nil {
		cell = grid[row][col]
	}
	return cell, err
}

// String - implements stringer interface for Grid
func (g Grid) String() string {
	buffer := bytes.NewBufferString("   0  1  2  3  4\n")
	for row := 0; row < gridSize; row++ {
		buffer.WriteString(fmt.Sprintf("%d ", row))
		for col := 0; col < gridSize; col++ {
			buffer.WriteString(grid[row][col].String())
		}
		buffer.WriteRune('\n')
	}
	return buffer.String()
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
		fmt.Println(grid)
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
	fmt.Println(grid)
}

// resetGrid - reverts the grid to its base state of blank cells
func resetGrid() {
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			cell := grid[row][col]
			if cell == nil {
				grid[row][col] = &Cell{row: row, col: col, val: defaultMark}
			} else {
				cell.val = defaultMark
			}
		}
	}
}

// takeTurn - will attempt to get input from the user and display errors if the input is invalid
func takeTurn(mark rune) {
	fmt.Println("[row],[col]")
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

	if !strings.ContainsRune(command, ',') {
		return chunks, errors.New("Command must contain a ',' between row and column")
	}

	stringChunks := strings.Split(command, ",")
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
			contiguousCount = handleMatch(contiguousCount, grid[row][col].val, symbol)
		}
		if contiguousCount == 4 {
			return true
		}
	}

	// check vertical
	for col := 0; col < gridSize; col++ {
		contiguousCount := 0
		for row := 0; row < gridSize && contiguousCount < 4; row++ {
			contiguousCount = handleMatch(contiguousCount, grid[row][col].val, symbol)
		}
		if contiguousCount == 4 {
			return true
		}
	}

	// check down-right angle
	// 0,0 1,1 2,2 3,3 4,4, etc.
	contiguousCount := 0
	for col, row := 0, 0; col < gridSize && contiguousCount < 4; col, row = col+1, row+1 {
		contiguousCount = handleMatch(contiguousCount, grid[row][col].val, symbol)
		if contiguousCount == 4 {
			return true
		}
	}

	//check down-left angles
	// 0,4 1,3 2,2 3,1 4,0, etc.
	contiguousCount = 0
	for col, row := 4, 0; col >= 0 && contiguousCount < 4; col, row = col-1, row+1 {
		contiguousCount = handleMatch(contiguousCount, grid[row][col].val, symbol)
		if contiguousCount == 4 {
			return true
		}
	}

	return false
}

// handleMatch - takes current contiguous count, cell value and active player's symbol
//	-> returning the updated contiguous count
func handleMatch(currentCount int, cellVal, symbol rune) int {
	if cellVal == symbol {
		currentCount++
	} else {
		currentCount = 0
	}
	return currentCount
}
