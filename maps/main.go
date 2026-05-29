package main

import "fmt"

func main() {
	// =========================================================
	// 1. DECLARING MAPS — three ways
	// =========================================================

	// Way A: var declaration creates a NIL map.
	// You can READ from a nil map (returns zero value) but you
	// CANNOT WRITE to one — writing to a nil map panics at runtime.
	var nilMap map[string]int
	fmt.Println("nilMap:", nilMap, "  is nil?", nilMap == nil)

	// Way B: make() — the standard way to create a writable map.
	// Returns an empty map, ready to use.
	scores := make(map[string]int)
	scores["alice"] = 90
	scores["bob"] = 75
	fmt.Println("scores after make + writes:", scores)

	// Way C: map literal — declare AND populate in one expression.
	// Trailing comma after the last entry is REQUIRED in multi-line form.
	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#00ff00",
		"blue":  "#0000ff",
	}
	fmt.Println("colors literal:", colors)

	// =========================================================
	// 2. READING — and the "comma ok" idiom
	// =========================================================

	// Reading a key that EXISTS — straightforward indexing.
	fmt.Println("alice's score:", scores["alice"])

	// Reading a key that DOESN'T EXIST — Go returns the ZERO VALUE
	// of the map's value type. NO error, NO panic. This is a gotcha.
	// For int, the zero value is 0. For string it's "". For bool it's false.
	fmt.Println("charlie's score:", scores["charlie"], "(charlie isn't in the map)")

	// To distinguish "key not present" from "key present with zero value",
	// use the COMMA-OK idiom — Go returns (value, bool).
	val, ok := scores["charlie"]
	fmt.Printf("comma-ok: charlie => value=%d, exists=%v\n", val, ok)

	val, ok = scores["alice"]
	fmt.Printf("comma-ok: alice   => value=%d, exists=%v\n", val, ok)

	// =========================================================
	// 3. UPDATING — same syntax as inserting
	// =========================================================

	// Go does NOT distinguish "add" vs "update" — assignment does both.
	scores["alice"] = 95
	fmt.Println("scores after updating alice to 95:", scores)

	// =========================================================
	// 4. DELETING
	// =========================================================

	// delete(map, key) is a built-in.
	delete(scores, "bob")
	fmt.Println("scores after deleting bob:", scores)

	// Deleting a key that doesn't exist is SILENTLY OK — no panic.
	delete(scores, "nonexistent")
	fmt.Println("scores after deleting nonexistent (no error):", scores)

	// =========================================================
	// 5. LENGTH
	// =========================================================

	// len() works on maps — returns the number of key-value pairs.
	fmt.Println("number of entries in scores:", len(scores))

	// =========================================================
	// 6. ITERATION — order is INTENTIONALLY RANDOM
	// =========================================================

	// Go randomizes map iteration order on purpose so code doesn't
	// accidentally depend on a particular order. Run the program
	// multiple times — the order of colors below will change.
	fmt.Println("iterating over colors (order varies each run):")
	for key, value := range colors {
		fmt.Printf("  %s => %s\n", key, value)
	}

	// If you only want keys, omit the second variable:
	//   for key := range colors { ... }
	// If you only want values, use _ for the key:
	//   for _, value := range colors { ... }

	// =========================================================
	// 7. MAPS ARE REFERENCE TYPES — the assignment gotcha
	// =========================================================

	// Unlike arrays, assigning a map to another variable does NOT
	// copy the data — both variables refer to the SAME underlying map.
	// Modifying one modifies the other.
	original := map[string]int{"x": 1}
	alias := original
	alias["x"] = 999
	fmt.Println("original after modifying 'alias':", original)
	// Prints x:999 — they share the same data.
}
