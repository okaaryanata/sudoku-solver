package main

import (
	"fmt"
	"time"
)

type SudokuRequest struct {
	Body map[int][]int
}

func main() {
	start := time.Now()
	fmt.Println("welcome to sudoku solver")

	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("\nThe function took %v seconds", elapsed.Seconds())
	}()

	sr := initSudokuBody()
	sr.PrintBody()

	var (
		isSolved bool
	)
	for !isSolved {
		err := sr.solver()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("===================================================")
		sr.PrintBody()
		fmt.Println("===================================================")
		isSolved = sr.isSudokuCompleted()
	}

	if isSolved {
		fmt.Println("===================================================")
		sr.PrintBody()
		fmt.Println("===================================================")
	}
}

func initSudokuBody() *SudokuRequest {
	// req := &SudokuRequest{
	// 	Body: map[int][]int{
	// 		0: {4, 5, 6, 0, 2, 3, 7, 8, 9},
	// 		1: {1, 2, 3, 4, 0, 6, 7, 8, 9},
	// 		2: {7, 8, 9, 7, 8, 9, 7, 8, 9},
	// 		3: {1, 2, 3, 2, 1, 6, 7, 8, 9},
	// 		4: {1, 2, 3, 3, 3, 6, 7, 8, 9},
	// 		5: {1, 2, 3, 9, 4, 6, 7, 8, 9},
	// 		6: {1, 2, 3, 5, 6, 6, 7, 8, 9},
	// 		7: {1, 2, 3, 6, 7, 6, 7, 8, 9},
	// 		8: {1, 2, 3, 8, 9, 6, 7, 8, 9},
	// 	},
	// }

	req := &SudokuRequest{
		Body: map[int][]int{
			0: {0, 8, 4, 0, 0, 0, 0, 3, 1},
			1: {6, 1, 0, 0, 0, 5, 0, 2, 0},
			2: {0, 5, 7, 0, 3, 8, 0, 0, 9},
			3: {8, 3, 2, 7, 5, 0, 4, 9, 0},
			4: {0, 4, 0, 0, 9, 0, 0, 1, 8},
			5: {0, 0, 6, 0, 8, 2, 0, 0, 3},
			6: {0, 0, 6, 0, 8, 2, 0, 0, 3},
			7: {3, 7, 8, 0, 1, 0, 9, 0, 0},
			8: {5, 6, 0, 0, 0, 0, 0, 0, 4},
		},
	}

	return req
}

func (s *SudokuRequest) GetRangeChecker(xAxis, yAxis int) ([]int, []int, error) {
	rangeX, err := s.getRangeMapper(xAxis)
	if err != nil {
		return nil, nil, err
	}

	rangeY, err := s.getRangeMapper(yAxis)
	if err != nil {
		return nil, nil, err
	}

	return rangeX, rangeY, nil
}

func (s *SudokuRequest) getRangeMapper(axis int) ([]int, error) {
	if axis < 0 || axis > 8 {
		return nil, fmt.Errorf("invalid axis value")
	}

	var xRange []int
	switch {
	case axis <= 2:
		xRange = []int{0, 1, 2}
	case axis > 2 && axis <= 5:
		xRange = []int{3, 4, 5}
	default:
		xRange = []int{6, 7, 8}
	}

	return xRange, nil
}

func initMapper() map[int]bool {
	return map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
		6: false,
		7: false,
		8: false,
		9: false,
	}
}

func intersection(arrays ...[]int) []int {
	if len(arrays) == 0 {
		return []int{}
	}

	freqMap := make(map[int]int)
	for i, array := range arrays {
		tempMap := make(map[int]bool)
		for _, num := range array {
			if i == 0 || !tempMap[num] {
				freqMap[num]++
				tempMap[num] = true
			}
		}
	}

	result := []int{}
	for num, count := range freqMap {
		if count == len(arrays) {
			result = append(result, num)
		}
	}

	return result
}

func (s *SudokuRequest) solver() error {
	for y, xList := range s.Body {
		for x, value := range xList {
			if value < 0 && value > 9 {
				return fmt.Errorf("not valid")
			}

			if s.Body[y][x] > 0 {
				continue
			}

			xRange, yRange, err := s.GetRangeChecker(x, y)
			if err != nil {
				return err
			}

			tmp := map[int][]int{}
			for _, yAxis := range yRange {
				if _, ok := tmp[yAxis]; !ok {
					tmp[yAxis] = []int{}
				}
				for _, xAxis := range xRange {
					tmp[yAxis] = append(tmp[yAxis], s.Body[yAxis][xAxis])
				}
			}

			// mapper
			yMapper := initMapper()
			xMapper := initMapper()
			boxMapper := initMapper()

			// X
			for _, value := range s.Body[y] {
				xMapper[value] = true
			}

			// Y
			for i := 0; i < 9; i++ {
				value := s.Body[i][x]
				yMapper[value] = true
			}

			// BOX
			for _, ints := range tmp {
				for _, v := range ints {
					boxMapper[v] = true
				}
			}

			// Probability
			yPossible := []int{}
			xPossible := []int{}
			boxPossible := []int{}
			for i := 1; i <= 9; i++ {
				if !yMapper[i] {
					yPossible = append(yPossible, i)
				}

				if !xMapper[i] {
					xPossible = append(xPossible, i)
				}

				if !boxMapper[i] {
					boxPossible = append(boxPossible, i)
				}
			}

			possibleValues := intersection(yPossible, xPossible, boxPossible)
			if len(possibleValues) > 0 {
				if len(possibleValues) == 1 {
					s.Body[y][x] = possibleValues[0]
				}

				fmt.Printf("Y:%d X:%d -> %d is possible\n", y, x, possibleValues)
				fmt.Println("===================================================")
			}

			if len(possibleValues) == 0 {
				return nil
			}
		}
	}

	return nil
}

func (s *SudokuRequest) PrintBody() {
	for i := 0; i < 9; i++ {
		fmt.Printf("%d : %d\n", i, s.Body[i])

	}
}

func (s *SudokuRequest) isSudokuCompleted() bool {
	// Helper to check if a slice contains exactly 1–9
	isValidGroup := func(group []int) bool {
		if len(group) != 9 {
			return false
		}
		seen := make(map[int]bool)
		for _, num := range group {
			if num < 1 || num > 9 || seen[num] {
				return false
			}
			seen[num] = true
		}
		return true
	}

	// Check all rows
	for _, row := range s.Body {
		if !isValidGroup(row) {
			return false
		}
	}

	// Check all columns
	for col := 0; col < 9; col++ {
		column := []int{}
		for row := 0; row < 9; row++ {
			column = append(column, s.Body[row][col])
		}
		if !isValidGroup(column) {
			return false
		}
	}

	// Check all 3x3 subgrids
	for gridRow := 0; gridRow < 3; gridRow++ {
		for gridCol := 0; gridCol < 3; gridCol++ {
			subgrid := []int{}
			for row := gridRow * 3; row < (gridRow+1)*3; row++ {
				for col := gridCol * 3; col < (gridCol+1)*3; col++ {
					subgrid = append(subgrid, s.Body[row][col])
				}
			}
			if !isValidGroup(subgrid) {
				return false
			}
		}
	}

	return true
}
