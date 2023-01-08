package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type Cell struct {
	Row int
	Col int
}

// Generation Represents 2D arrays of Generations.
type Generation [][]string

const numOfNeighbors int = 8
const numOfGenerations int = 10

func main() {
	var curGen, nextGen Generation
	N := getInput()
	curGen = fillMatrix(N)

	// create new generations
	for i := 1; i <= numOfGenerations; i++ {
		nextGen = evolve(curGen)

		fmt.Printf("Generation #%d\n", i)
		fmt.Printf("Alive: %d\n", countAliveInGeneration(&nextGen))
		copy(curGen, nextGen)
		printMatrix(curGen)
	}
}

func evolve(g Generation) Generation {
	var cell = Cell{}
	var nextGen = make(Generation, len(g))
	//fmt.Printf("evolve(): len=%v, cap=%v\n", len(nextGen), cap(nextGen))

	cellMap := make(map[Cell]int)
	for i := 0; i < len(g); i++ {
		nextGen[i] = make([]string, 0, len(nextGen))
		//fmt.Printf("evolve(): len[i]=%v, cap[i]=%v\n", len(nextGen[i]), cap(nextGen[i]))
		for j := 0; j < len(g); j++ {
			cell = Cell{i, j}
			value := g[cell.Row][cell.Col]
			if _, ok := cellMap[cell]; ok {
				continue
			}
			neighbors := getNeighbors(len(g), len(g[i]), cell)
			cellMap[cell] = countAliveInNeighbors(&g, &neighbors)
			////fmt.Printf("evolve(): cellMap=%v\n", cellMap)

			switch value {
			case "O":
				// A live cell survives if it has two or three
				// live neighbors; otherwise, it dies of boredom
				// (<2) or overpopulation (>3).
				if cellMap[cell] == 2 || cellMap[cell] == 3 {
					//fmt.Printf("\tcell: %v=%q SURVIVES, alive_count=%v\n",
					//	cell, value, cellMap[cell])
					nextGen[i] = append(nextGen[i], "O")
				} else {
					//fmt.Printf("\tcell: %v=%q DIES, alive_count=%v\n",
					//	cell, value, cellMap[cell])
					nextGen[i] = append(nextGen[i], " ")
				}
			case " ":
				// A dead cell is reborn if it has exactly three
				// live neighbors.
				if value == " " && cellMap[cell] == 3 {
					//fmt.Printf("\tcell: %v=%q REBORN, alive_count=%v\n",
					//	cell, value, cellMap[cell])
					nextGen[i] = append(nextGen[i], "O")
				} else {
					// stays dead
					//fmt.Printf("\tcell: %v=%q DIES, alive_count=%v\n",
					//	cell, value, cellMap[cell])
					nextGen[i] = append(nextGen[i], " ")
				}
			}
		}
	}
	//fmt.Printf("EVOLVE END: %v\n", nextGen)
	return nextGen
}

func countAliveInGeneration(m *Generation) int {
	alive := 0
	for _, slice := range *m {
		s := strings.Join(slice, "")
		alive += strings.Count(s, "O")
	}
	return alive
}

func countAliveInNeighbors(m *Generation, cs *[]Cell) int {
	alive := 0
	for _, v := range *cs {
		if alive == numOfNeighbors {
			// cannot have more than 8 alive neighbors
			return alive
		}
		if (*m)[v.Row][v.Col] == "O" {
			alive++
			//fmt.Printf("countAlive(): m[%v][%v]=%v, val=%v\n",
			//	v.Row, v.Col, (*m)[v.Row][v.Col], v)
		}
	}
	return alive
}

func getNeighbors(lrow, lcol int, c Cell) []Cell {
	rows, cols := lrow, lcol
	//fmt.Printf("\ngetNei(): CELL: %v\n", c)
	neighbors := make([]Cell, 0, numOfNeighbors) // each cell has 8 neighbor

	for i := c.Row - 1; i <= c.Row+1; i++ {
		for j := c.Col - 1; j <= c.Col+1; j++ {
			if (i >= 0 && i < rows) &&
				(j >= 0 && j < cols) &&
				!(i == c.Row && j == c.Col) {
				//fmt.Printf("\t-> cell is inbetween: Cell{%d, %d}\n", i, j)
				//fmt.Printf("\t\t-> adding: %v\n", Cell{i, j})
				neighbors = append(neighbors, Cell{i, j})
			}
		}
	}

	// STOP processing if the length of `neighbors` is 8!
	if len(neighbors) != numOfNeighbors {
		// need to pass the original matrix length
		neighbors = processDirection(neighbors, c, lrow, lcol)
	}
	// found all neighbors

	sort.Slice(neighbors, func(i, j int) bool {
		if neighbors[i].Row != neighbors[j].Row {
			return neighbors[i].Row < neighbors[j].Row
		}
		return neighbors[i].Col < neighbors[j].Col
	})
	//fmt.Printf("getNei(): neighbor=%v\n", neighbors)
	return neighbors
}

func processDirection(n []Cell, c Cell, mRow, lCol int) []Cell {
	directionMap := make(map[Cell]bool, 0)
	for i := range n {
		rc := new(Cell)

		if onTopBorder(c) {
			*rc = Cell{mRow - 1, n[i].Col}
			//fmt.Printf("\t\t-> adding: %v\n", *rc)
		} else if onBottomBorder(mRow, c) {
			*rc = Cell{0, n[i].Col}
			//fmt.Printf("\t\t-> adding: %v\n", *rc)
		} else if onRightBorder(lCol, c) {
			*rc = Cell{n[i].Row, 0}
			//fmt.Printf("\t\t-> adding: %v\n", *rc)
		} else if onLeftBorder(c) {
			*rc = Cell{n[i].Row, lCol - 1}
			//fmt.Printf("\t\t-> adding: %v\n", *rc)
		}

		if _, ok := directionMap[*rc]; !ok {
			directionMap[*rc] = true
			n = append(n, *rc)
		}
	}

	if cs, ok := onCorner(mRow, lCol, c); ok {
		//fmt.Printf("\t\t-> adding: %v\n", cs)
		n = append(n, cs...)
	}

	return n
}

// If cell is right-border, its right (east) neighbor
// is leftmost cell in the same row.
func onRightBorder(l int, c Cell) bool {
	cols := l
	if c.Col == cols-1 && c.Row >= 0 {
		//fmt.Printf("\t-> cell is on RIGHT: %v\n", c)
		return true
	}
	return false
}

func onLeftBorder(c Cell) bool {
	if c.Col == 0 && c.Row >= 0 {
		//fmt.Printf("\t-> cell is on LEFT: %v\n", c)
		return true
	}
	return false
}

func onTopBorder(c Cell) bool {
	if c.Row == 0 && c.Col >= 0 {
		//fmt.Printf("\t-> cell is on TOP: %v\n", c)
		return true
	}
	return false
}

// If cell is bottom-border, its bottom (south) neighbor
// is topmost cell in the same column.
func onBottomBorder(l int, c Cell) bool {
	rows := l
	if c.Row == rows-1 && c.Col >= 0 {
		//fmt.Printf("\t-> cell is on BOTTOM: %v\n", c)
		return true
	}
	return false
}

func onCorner(lr, lc int, c Cell) ([]Cell, bool) {
	rows, cols := lr, lc
	//fmt.Printf("onCorner(): lr=%v, lc=%v\n", lr, lc)
	topRight := Cell{0, cols - 1}
	topLeft := Cell{0, 0}
	bottomRight := Cell{rows - 1, cols - 1}
	bottomLeft := Cell{rows - 1, 0}

	var corner bool
	var cs = make([]Cell, 0, cols)

	switch c {
	case topRight:
		corner = true
		//fmt.Printf("\t-> cell is on TOP RIGHT corner: %v\n", topRight)
		cs = append(cs, Cell{0, 0}, Cell{1, 0}, Cell{rows - 1, 0})
	case topLeft:
		corner = true
		//fmt.Printf("\t-> cell is on TOP LEFT corner: %v\n", topLeft)
		cs = append(cs, Cell{0, cols - 1}, Cell{1, cols - 1}, Cell{rows - 1, cols - 1})
	case bottomRight:
		corner = true
		//fmt.Printf("\t-> cell is on BOTTOM RIGHT corner: %v\n", bottomRight)
		cs = append(cs, Cell{0, 0}, Cell{rows - 2, 0}, Cell{rows - 1, 0})
	case bottomLeft:
		corner = true
		//fmt.Printf("\t-> cell is on BOTTOM LEFT corner: %v\n", bottomLeft)
		cs = append(cs, Cell{0, cols - 1}, Cell{rows - 2, cols - 1}, Cell{rows - 1, cols - 1})
	default:
		corner = false
		//fmt.Printf("onCorner(): CELL: %v -> NOT IN CORNER -> %t\n", c, corner)
		cs = nil
	}

	return cs, corner
}

func printMatrix(m Generation) {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m); j++ {
			fmt.Print(m[i][j])
		}
		fmt.Println()
	}
}

func fillMatrix(N int) Generation {
	matrix := make(Generation, N)

	for i := 0; i < N; i++ {
		matrix[i] = make([]string, N)
		for j := 0; j < N; j++ {
			if rand.Intn(2) == 1 {
				matrix[i][j] = "O"
			} else {
				matrix[i][j] = " "
			}
		}
	}

	return matrix
}

func getInput() int {
	var N int // size of the universe
	fmt.Scanf("%d", &N)
	return N
}
