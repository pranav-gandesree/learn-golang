# Arrays and Slices

> See [arrays/main.go](../arrays/main.go) and [slicegrowth/main.go](../slicegrowth/main.go) for runnable demos.

## Arrays — fixed size, baked into the type

```go
var a [3]int                  // [0 0 0]   — zero-initialized
b := [3]int{1, 2, 3}
c := [...]int{1, 2, 3}        // size inferred from elements → still [3]int
```

- **Size is part of the type.** `[3]int` ≠ `[4]int`. You CANNOT change an array's size.
- Assignment COPIES the whole array.
- Passing to a function COPIES it (expensive for big arrays).

Arrays are rare in idiomatic Go. Slices are what you'll use 99% of the time.

## Slices — dynamic, three-field header

Internally, a slice is a tiny struct:

```
slice = { pointer to backing array, len, cap }
```

- `len` — how many elements you can access
- `cap` — size of the underlying backing array
- They're independent; `len ≤ cap` always

## Creating slices

```go
s := []int{4, 5, 6}              // literal — len=3, cap=3
s := make([]int, 5)              // [0 0 0 0 0] — len=5, cap=5
s := make([]int, 0, 10)          // [] — len=0, cap=10
s := make([]int, 5, 100)         // [0 0 0 0 0] — len=5, cap=100
```

**Critical gotcha:** `make([]int, 5, 100)` already contains 5 zeros — not an "empty slice with room for 5." If you want empty-but-with-capacity, use `make([]int, 0, 100)`.

## Reading and writing

```go
s[0]                  // read at index 0
s[0] = 99             // write at index 0
s[len(s)]             // PANIC — out of range
```

Valid indices are `0` through `len-1`. Capacity does NOT make extra indices accessible.

## `append` — adding elements

```go
s = append(s, 7)                     // add one
s = append(s, 8, 9, 10)              // add multiple
s = append(s, otherSlice...)         // spread another slice (... is required)
```

**ALWAYS reassign the result.** `append` returns a new slice header. If `len < cap`, the underlying array is reused; if not, a new bigger array is allocated.

## Slice growth (the doubling pattern)

When `append` exceeds `cap`, Go:
1. Allocates a NEW backing array (roughly 2× the old cap for small slices)
2. Copies all existing elements over
3. Adds the new element

Demo output:
```
start: len=0 cap=0
after append  1: len=1 cap=4
after append  4: len=4 cap=4
after append  5: len=5 cap=8       ← doubled
after append  8: len=8 cap=8
after append  9: len=9 cap=16      ← doubled
after append 17: len=17 cap=32     ← doubled
```

This is the same amortized-O(1) trick used by C++ `std::vector`, Java `ArrayList`, Python `list`. Capacity is reserved space so future appends don't always trigger reallocation.

## `append` doesn't inherit source capacity

```go
src := make([]int, 5, 111)   // len=5, cap=111
dst := []int{1, 2, 3}        // len=3, cap=3
dst = append(dst, src...)    // dst's cap is now ~8-16, NOT 111+
```

`append` uses the DESTINATION's growth rules. The source's capacity is irrelevant.

## Slices share memory

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:4]               // [2 3 4] — shares a's backing array
b[0] = 999
fmt.Println(a)            // [1 999 3 4 5] — a is modified too
```

To get an independent copy, use the built-in `copy`:

```go
b := make([]int, 3)
copy(b, a[1:4])           // b is independent
```

## Quick reference

```go
// Create
s := []int{1, 2, 3}
s := make([]int, len, cap)

// Read/write
s[i]
s[i] = v

// Grow
s = append(s, v)
s = append(s, otherSlice...)

// Stats
len(s)
cap(s)

// Sub-slice
s[i:j]                    // elements i to j-1
s[:n]                     // first n
s[n:]                     // from n to end

// Independent copy
copied := make([]T, len(s))
copy(copied, s)
```

## Arrays vs slices at a glance

| | Array `[3]int` | Slice `[]int` |
|---|---|---|
| Size fixed? | Yes — part of the type | No — grows via `append` |
| Assignment | COPY of all elements | Shares backing array |
| Pass to function | COPY (expensive for big arrays) | Cheap (just the 3-word header) |
| Use when | Rarely | Default |
