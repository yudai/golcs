package lcs

import (
	"errors"
	"reflect"
)

func Lcs(left, right []interface{}) (lcs []interface{}) {
	table := Table(left, right)
	lcs, err := LcsFromTable(left, right, table)
	if err != nil {
		// error cannot happen here
		panic(err.Error())
	}
	return
}

func Length(left, right []interface{}) (length int) {
	table := Table(left, right)
	length = table[len(left)][len(right)]
	return
}

func Table(left, right []interface{}) (table [][]int) {
	sizeX := len(left) + 1
	sizeY := len(right) + 1

	table = make([][]int, sizeX)
	for x := 0; x < sizeX; x++ {
		table[x] = make([]int, sizeY)
	}

	for y := 1; y < sizeY; y++ {
		for x := 1; x < sizeX; x++ {
			increment := 0
			if reflect.DeepEqual(left[x-1], right[y-1]) {
				increment = 1
			}
			table[x][y] = max(table[x-1][y-1]+increment, table[x-1][y], table[x][y-1])
		}
	}

	return
}

func LcsFromTable(left, right []interface{}, table [][]int) (lcs []interface{}, err error) {
	if len(left)+1 != len(table) {
		return nil, errors.New("Table size and length of first argument doesn't match")
	}
	if len(left) == 0 {
		return []interface{}{}, nil
	}

	lcs = make([]interface{}, table[len(table)-1][len(table[0])-1])

	for x, y := len(left), len(right); x > 0 && y > 0; {
		decrease := false
		if reflect.DeepEqual(left[x-1], right[y-1]) {
			lcs[table[x][y]-1] = left[x-1]
			decrease = true
		}

		xDelta := 0
		if decrease || table[x-1][y] >= table[x][y-1] {
			xDelta++
		}

		if decrease || table[x-1][y] < table[x][y-1] {
			y--

			if y > len(table[x]) {
				return nil, errors.New("Table size and length of second argument dosn't match")
			}
		}
		x -= xDelta
	}

	return
}

func max(first int, rest ...int) (max int) {
	max = first
	for _, value := range rest {
		if value > max {
			max = value
		}
	}
	return
}
