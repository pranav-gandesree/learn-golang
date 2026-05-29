# Maps

> See [maps/main.go](../maps/main.go) for a runnable demo.

A **map** is Go's hash table — unordered key→value pairs with unique keys. Declared as `map[KeyType]ValueType`.

## Declaring maps — three ways

```go
// Way A — make() — standard way for an empty writable map
m := make(map[string]int)

// Way B — map literal — declare and populate at once
m := map[string]int{
    "alice": 90,
    "bob":   75,                  // trailing comma REQUIRED on multi-line
}

// Way C — var — creates a NIL MAP (gotcha!)
var m map[string]int
// Reading is OK (returns zero value).
// Writing to a nil map PANICS at runtime.
```

**Rule:** Always use `make` or a literal. `var m map[K]V` produces a nil map you cannot write to.

## Operations

```go
// Insert / update — same syntax
m["alice"] = 95

// Read (always safe)
v := m["alice"]                   // zero value if absent
v, ok := m["alice"]               // comma-ok idiom

// Delete
delete(m, "alice")                // no-op if key absent — no error

// Size
len(m)
```

## The comma-ok idiom

Reading a missing key returns the **zero value** of the value type — `0`, `""`, `false`, or `nil`. There's NO error, NO panic.

```go
m := map[string]int{"alice": 90}
m["charlie"]                      // 0 — not in map, but no error
```

To distinguish "key present with value 0" from "key absent," use comma-ok:

```go
v, ok := m["charlie"]
// v = 0, ok = false → key was absent

v, ok = m["alice"]
// v = 90, ok = true → key was present
```

Use this anywhere zero is a valid value (counters, scores, prices).

## Iteration — ORDER IS RANDOM

```go
for k, v := range m {
    fmt.Println(k, "=>", v)
}
```

Go **intentionally** randomizes map iteration order. Run twice → different orders. Never depend on a particular order.

Variants:
```go
for k := range m         { ... }   // keys only
for _, v := range m      { ... }   // values only
```

## Sorted iteration

If you need a deterministic order, sort the keys separately:

```go
import "sort"

keys := make([]string, 0, len(m))
for k := range m { keys = append(keys, k) }
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, "=>", m[k])
}
```

## Maps are REFERENCE TYPES

Assignment doesn't copy data — both variables refer to the same hash table:

```go
original := map[string]int{"x": 1}
alias    := original
alias["x"] = 999
// original["x"] is now 999 too
```

Same applies to passing maps to functions — the function can mutate the caller's map. This is usually what you want, but be deliberate.

To deep-copy:

```go
clone := make(map[string]int, len(original))
for k, v := range original {
    clone[k] = v
}
```

## Comparison with other types

| Type | Assignment copies data? |
|---|---|
| Array `[3]int` | YES |
| Slice `[]int` | NO (shares backing array) |
| Map `map[K]V` | NO (shares the hash table) |

## Quick reference

| What you want | How |
|---|---|
| Empty writable map | `make(map[K]V)` |
| Map with initial entries | `map[K]V{...}` |
| Check if key exists | `v, ok := m[k]` |
| Get with zero-value default | `v := m[k]` |
| Insert or update | `m[k] = v` |
| Remove a key | `delete(m, k)` |
| Count entries | `len(m)` |
| Iterate (random order) | `for k, v := range m` |
| Iterate sorted | Collect keys → sort → loop |

## Common gotchas

1. **Nil maps** — `var m map[K]V` can be read but writing panics. Use `make`.
2. **Random iteration** — never rely on order.
3. **Reference semantics** — `alias := original` shares state. Copy explicitly.
4. **Missing keys return zero values** — use comma-ok when zero is a legitimate value.
