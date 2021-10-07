// package lcs provides functions to calculate Longest Common Subsequence (LCS)
// values from two arbitrary arrays.
package lcs

import (
	"context"
	"reflect"
)

// Lcs is the interface to calculate the LCS of two arrays.
type Lcs interface {
	// Values calculates the LCS value of the two arrays.
	Values() (values []interface{})
	// ValueContext is a context aware version of Values()
	ValuesContext(ctx context.Context) (values []interface{}, err error)
	// IndexPairs calculates paris of indices which have the same value in LCS.
	IndexPairs() (pairs []IndexPair)
	// IndexPairsContext is a context aware version of IndexPairs()
	IndexPairsContext(ctx context.Context) (pairs []IndexPair, err error)
	// Length calculates the length of the LCS.
	Length() (length int)
	// LengthContext is a context aware version of Length()
	LengthContext(ctx context.Context) (length int, err error)
	// Left returns one of the two arrays to be compared.
	Left() (leftValues []interface{})
	// Right returns the other of the two arrays to be compared.
	Right() (righttValues []interface{})
}

// IndexPair represents an pair of indeices in the Left and Right arrays found in the LCS value.
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

// New creates a new LCS calculator from two arrays.
func New(left, right []interface{}) Lcs {
	return &lcs{
		left:       left,
		right:      right,
		table:      nil,
		indexPairs: nil,
		values:     nil,
	}
}

// Table implements Lcs.Table()
func (lcs *lcs) Table() (table [][]int) {
	table, _ = lcs.TableContext(context.Background())
	return table
}

// Table implements Lcs.TableContext()
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

// Table implements Lcs.Length()
func (lcs *lcs) Length() (length int) {
	length, _ = lcs.LengthContext(context.Background())
	return length
}

func (lcs *lcs) lengthContext(ctx context.Context) (length int, err error) {
	m := len(lcs.left)
	n := len(lcs.right)

	// allocate storage for one-dimensional array `curr`
	prev := 0
	curr := make([]int, n+1)

	// fill the lookup table in a bottom-up manner
	for i := 0; i <= m; i++ {
		select { // check in each y to save some time
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			// nop
		}
		prev = curr[0]
		for j := 0; j <= n; j++ {
			backup := curr[j]
			if i == 0 || j == 0 {
				curr[j] = 0
			} else if reflect.DeepEqual(lcs.left[i-1], lcs.right[j-1]) {
				// if the current character of `X` and `Y` matches
				curr[j] = prev + 1
			} else {
				// otherwise, if the current character of `X` and `Y` don't match
				curr[j] = max(curr[j], curr[j-1])
			}
			prev = backup
		}
	}
	// LCS will be the last entry in the lookup table
	return curr[n], nil
}

// Table implements Lcs.LengthContext()
func (lcs *lcs) LengthContext(ctx context.Context) (length int, err error) {
	if len(lcs.right) > len(lcs.left) {
		lcs.left, lcs.right = lcs.right, lcs.left
	}
	return lcs.lengthContext(ctx)
}

// Table implements Lcs.IndexPairs()
func (lcs *lcs) IndexPairs() (pairs []IndexPair) {
	pairs, _ = lcs.IndexPairsContext(context.Background())
	return pairs
}

// Table implements Lcs.IndexPairsContext()
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

// Table implements Lcs.Values()
func (lcs *lcs) Values() (values []interface{}) {
	values, _ = lcs.ValuesContext(context.Background())
	return values
}

// Table implements Lcs.ValuesContext()
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

// Table implements Lcs.Left()
func (lcs *lcs) Left() (leftValues []interface{}) {
	leftValues = lcs.left
	return
}

// Table implements Lcs.Right()
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
