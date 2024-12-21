package main

import (
	"fmt"
	"time"
)

var MaxLength int = 9

type SudokuRequest struct {
	Body map[int][]int
}

func main() {
	fmt.Println("Welcome to Sudoku Solver!")

	// initialize the Sudoku puzzle
	sudoku := initSudokuBody()
	fmt.Println("\nInitial Sudoku Board:")
	sudoku.PrintBody()

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("\nThe function took %.6f seconds", elapsed.Seconds())
	}()
	// solve the Sudoku puzzle
	if err := sudoku.Solve(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nSolved Sudoku Board:")
	sudoku.PrintBody()
}

func initSudokuBody() *SudokuRequest {
	// example sudoku
	return &SudokuRequest{
		Body: map[int][]int{
			0: {0, 0, 3, 0, 2, 0, 6, 0, 0},
			1: {9, 0, 0, 3, 0, 5, 0, 0, 1},
			2: {0, 0, 1, 8, 0, 6, 4, 0, 0},
			3: {0, 0, 8, 1, 0, 2, 9, 0, 0},
			4: {7, 0, 0, 0, 0, 0, 0, 0, 8},
			5: {0, 0, 6, 7, 0, 8, 2, 0, 0},
			6: {0, 0, 2, 6, 0, 9, 5, 0, 0},
			7: {8, 0, 0, 2, 0, 3, 0, 0, 9},
			8: {0, 0, 5, 0, 1, 0, 3, 0, 0},
		},
	}
}

func (s *SudokuRequest) Solve() error {
	for {
		updated := false

		for y, row := range s.Body {
			for x, value := range row {
				if value != 0 {
					continue
				}

				possibleValues, err := s.findPossibleValues(x, y)
				if err != nil {
					return err
				}

				if len(possibleValues) == 1 {
					s.Body[y][x] = possibleValues[0]
					updated = true
				}
			}
		}

		if s.isSudokuCompleted() {
			return nil
		}

		if !updated {
			return fmt.Errorf("unable to solve Sudoku with current logic")
		}
	}
}

func (s *SudokuRequest) findPossibleValues(x, y int) ([]int, error) {
	if x < 0 || x >= MaxLength || y < 0 || y >= MaxLength {
		return nil, fmt.Errorf("invalid coordinates (%d, %d)", x, y)
	}

	xRange, yRange, err := s.GetRangeChecker(x, y)
	if err != nil {
		return nil, err
	}

	used := make(map[int]bool)

	// check row and column
	for i := 0; i < 9; i++ {
		used[s.Body[y][i]] = true
		used[s.Body[i][x]] = true
	}

	// check 3x3 grid
	for _, row := range yRange {
		for _, col := range xRange {
			used[s.Body[row][col]] = true
		}
	}

	// collect unused values
	possibleValues := []int{}
	for i := 1; i <= MaxLength; i++ {
		if !used[i] {
			possibleValues = append(possibleValues, i)
		}
	}

	return possibleValues, nil
}

func (s *SudokuRequest) PrintBody() {
	fmt.Println("=======================")
	for y := 0; y < 9; y++ {
		fmt.Println(s.Body[y][:3], s.Body[y][3:6], s.Body[y][6:])
		switch y {
		case 2, 5:
			fmt.Println("-----------------------")
		}
	}
	fmt.Println("=======================")
}

func (s *SudokuRequest) isSudokuCompleted() bool {
	// validate all rows
	for _, row := range s.Body {
		if !isValidGroup(row) {
			return false
		}
	}

	// validate all columns
	for x := 0; x < MaxLength; x++ {
		column := []int{}
		for y := 0; y < MaxLength; y++ {
			column = append(column, s.Body[y][x])
		}
		if !isValidGroup(column) {
			return false
		}
	}

	// validate 3x3 grids
	for gridRow := 0; gridRow < 3; gridRow++ {
		for gridCol := 0; gridCol < 3; gridCol++ {
			subgrid := []int{}
			for y := gridRow * 3; y < (gridRow+1)*3; y++ {
				for x := gridCol * 3; x < (gridCol+1)*3; x++ {
					subgrid = append(subgrid, s.Body[y][x])
				}
			}
			if !isValidGroup(subgrid) {
				return false
			}
		}
	}

	return true
}

func (s *SudokuRequest) GetRangeChecker(x, y int) ([]int, []int, error) {
	getRange := func(idx int) ([]int, error) {
		switch {
		case idx >= 0 && idx < 3:
			return []int{0, 1, 2}, nil
		case idx >= 3 && idx < 6:
			return []int{3, 4, 5}, nil
		case idx >= 6 && idx < 9:
			return []int{6, 7, 8}, nil
		default:
			return nil, fmt.Errorf("invalid index %d", idx)
		}
	}

	xRange, err := getRange(x)
	if err != nil {
		return nil, nil, err
	}

	yRange, err := getRange(y)
	if err != nil {
		return nil, nil, err
	}

	return xRange, yRange, nil
}

func isValidGroup(group []int) bool {
	if len(group) != MaxLength {
		return false
	}

	seen := make(map[int]bool)
	for _, value := range group {
		if value < 1 || value > MaxLength || seen[value] {
			return false
		}
		seen[value] = true
	}

	return true
}
