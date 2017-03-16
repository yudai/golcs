package lcs

import (
	"context"
	"reflect"
)

type Lcs interface {
	Values() (values []interface{})
	ValuesContext(ctx context.Context) (values []interface{}, err error)
	IndexPairs() (pairs []IndexPair)
	IndexPairsContext(ctx context.Context) (pairs []IndexPair, err error)
	Length() (length int)
	LengthContext(ctx context.Context) (length int, err error)
	Left() (leftValues []interface{})
	Right() (righttValues []interface{})
}

type IndexPair struct {
	Left  int
	Right int
}

type lcs struct {
	left  []interface{}
	right []interface{}
	/* for caching */
	table      [][]int
	indexPairs []IndexPair
	values     []interface{}
}

func New(left, right []interface{}) Lcs {
	return &lcs{
		left:       left,
		right:      right,
		table:      nil,
		indexPairs: nil,
		values:     nil,
	}
}

func (lcs *lcs) Table() (table [][]int) {
	table, _ = lcs.TableContext(context.Background())
	return table
}

func (lcs *lcs) TableContext(ctx context.Context) (table [][]int, err error) {
	if lcs.table != nil {
		return lcs.table, nil
	}

	sizeX := len(lcs.left) + 1
	sizeY := len(lcs.right) + 1

	table = make([][]int, sizeX)
	for x := 0; x < sizeX; x++ {
		table[x] = make([]int, sizeY)
	}

	for y := 1; y < sizeY; y++ {
		select { // check in each y to save some time
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// nop
		}
		for x := 1; x < sizeX; x++ {
			increment := 0
			if reflect.DeepEqual(lcs.left[x-1], lcs.right[y-1]) {
				increment = 1
			}
			table[x][y] = max(table[x-1][y-1]+increment, table[x-1][y], table[x][y-1])
		}
	}

	lcs.table = table
	return table, nil
}

func (lcs *lcs) Length() (length int) {
	length, _ = lcs.LengthContext(context.Background())
	return length
}

func (lcs *lcs) LengthContext(ctx context.Context) (length int, err error) {
	table, err := lcs.TableContext(ctx)
	if err != nil {
		return 0, err
	}
	return table[len(lcs.left)][len(lcs.right)], nil
}

func (lcs *lcs) IndexPairs() (pairs []IndexPair) {
	pairs, _ = lcs.IndexPairsContext(context.Background())
	return pairs
}

func (lcs *lcs) IndexPairsContext(ctx context.Context) (pairs []IndexPair, err error) {
	if lcs.indexPairs != nil {
		return lcs.indexPairs, nil
	}

	table, err := lcs.TableContext(ctx)
	if err != nil {
		return nil, err
	}

	pairs = make([]IndexPair, table[len(table)-1][len(table[0])-1])

	for x, y := len(lcs.left), len(lcs.right); x > 0 && y > 0; {
		if reflect.DeepEqual(lcs.left[x-1], lcs.right[y-1]) {
			pairs[table[x][y]-1] = IndexPair{Left: x - 1, Right: y - 1}
			x--
			y--
		} else {
			if table[x-1][y] >= table[x][y-1] {
				x--
			} else {
				y--
			}
		}
	}

	lcs.indexPairs = pairs

	return pairs, nil
}

func (lcs *lcs) Values() (values []interface{}) {
	values, _ = lcs.ValuesContext(context.Background())
	return values
}

func (lcs *lcs) ValuesContext(ctx context.Context) (values []interface{}, err error) {
	if lcs.values != nil {
		return lcs.values, nil
	}

	pairs, err := lcs.IndexPairsContext(ctx)
	if err != nil {
		return nil, err
	}

	values = make([]interface{}, len(pairs))
	for i, pair := range pairs {
		values[i] = lcs.left[pair.Left]
	}
	lcs.values = values

	return values, nil
}

func (lcs *lcs) Left() (leftValues []interface{}) {
	leftValues = lcs.left
	return
}

func (lcs *lcs) Right() (rightValues []interface{}) {
	rightValues = lcs.right
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
