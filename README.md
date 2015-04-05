# Go Longest Common Subsequence (LCS)

## Usage

```sh
go get github.com/yudai/golcs
```

```go
left = []interface{}{1, 2, 5, 3, 1, 1, 5, 8, 3}
right = []interface{}{1, 2, 3, 3, 4, 4, 5, 1, 6}

golcs.Lcs(left, right) // => []interface{}{1, 2, 5, 1}

// If you just need their LCS length
golcs.Length(left, right) // => 4
```

## FAQ

### How can I give `[]byte` values to `Lcs()` as its arguments?

As `[]interface{}` is incompatible with `[]othertype` like `[]byte`, you need to create a `[]interface{}` slice and copy the values in your `[]byte` slice into it. Unfortunately, Go doesn't provide any mesure to cast a slice into `[]interface{}` with zero cost. Your copy costs O(n).

```go
leftBytes := []byte("TGAGTA")
left = make([]interface{}, len(leftBytes))
for i, v := range leftBytes {
	left[i] = v
}

rightBytes := []byte("GATA")
right = make([]interface{}, len(rightBytes))
for i, v := range rightBytes {
	right[i] = v
}

Lcs(left, right)
```


## LICENSE

The MIT license (See `LICENSE` for detail)
