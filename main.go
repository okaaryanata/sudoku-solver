package main

import (
	"fmt"
	"time"
)

const MaxLength = 9

type SudokuRequest struct {
	Body [][]int
}

func main() {
	fmt.Println("Welcome to Sudoku Solver!")

	sudoku := initSudokuBody()
	fmt.Println("\nInitial Sudoku Board:")
	sudoku.PrintBody()

	start := time.Now()
	defer func() {
		fmt.Printf("\nThe function took %.6f seconds\n", time.Since(start).Seconds())
	}()

	if sudoku.Solve() {
		fmt.Println("\nSolved Sudoku Board:")
		sudoku.PrintBody()
	} else {
		fmt.Println("\nUnable to solve Sudoku with the given input.")
	}
}

func initSudokuBody() *SudokuRequest {
	return &SudokuRequest{
		Body: [][]int{
			{0, 2, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0, 9, 7, 0, 0},
			{4, 0, 3, 0, 0, 0, 8, 0, 0},
			{3, 0, 2, 0, 0, 5, 0, 0, 0},
			{0, 0, 0, 0, 0, 7, 9, 1, 0},
			{0, 0, 6, 0, 0, 0, 0, 0, 4},
			{6, 8, 0, 0, 0, 0, 4, 0, 0},
			{0, 0, 0, 7, 0, 0, 0, 3, 0},
			{0, 0, 0, 5, 0, 0, 0, 0, 2},
		},
	}
}

func (s *SudokuRequest) Solve() bool {
	row, col, found := s.findEmptyCell()
	if !found {
		return true
	}

	for num := 1; num <= MaxLength; num++ {
		if s.isValidPlacement(row, col, num) {
			s.Body[row][col] = num
			if s.Solve() {
				return true
			}
			s.Body[row][col] = 0
		}
	}

	return false
}

func (s *SudokuRequest) findEmptyCell() (int, int, bool) {
	for row := 0; row < MaxLength; row++ {
		for col := 0; col < MaxLength; col++ {
			if s.Body[row][col] == 0 {
				return row, col, true
			}
		}
	}
	return -1, -1, false
}

func (s *SudokuRequest) isValidPlacement(row, col, num int) bool {
	for i := 0; i < MaxLength; i++ {
		if s.Body[row][i] == num || s.Body[i][col] == num {
			return false
		}
	}

	startRow, startCol := (row/3)*3, (col/3)*3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if s.Body[startRow+i][startCol+j] == num {
				return false
			}
		}
	}

	return true
}

func (s *SudokuRequest) PrintBody() {
	fmt.Println("=======================")
	for row := 0; row < MaxLength; row++ {
		fmt.Println(s.Body[row][:3], s.Body[row][3:6], s.Body[row][6:])
		if row == 2 || row == 5 {
			fmt.Println("-----------------------")
		}
	}
	fmt.Println("=======================")
}
