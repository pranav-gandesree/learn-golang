# Go Packages

## The four rules

1. **One folder = one package.** Every `.go` file in a folder must declare the same `package` name.
2. **Every file needs `package X` as its first non-comment line.** No exceptions.
3. **Capitalized names are EXPORTED** (visible from other packages). Lowercase names are package-private.
4. **Import path = module name (from `go.mod`) + subfolder path.**

## File-to-package mapping

```
learning/
├── go.mod                  module learning
├── main.go    \
├── hello.go   |── all "package main" — share the same namespace
├── loops.go   /
├── mathutils/
│   └── math.go             package mathutils
└── arrays/
    └── main.go             package main (separate program)
```

Files inside the same folder share the same namespace — they don't need to import each other. Files in DIFFERENT folders need explicit imports.

## Same-package vs different-package calls

| In same package | In different package |
|---|---|
| `Add(2, 3)` | `mathutils.Add(2, 3)` |
| Direct call, no prefix | `package.Func(args)` |

## Exported vs unexported

```go
package mathutils

func Add(a, b int) int       { return a + b } // EXPORTED — callable from elsewhere
func multiply(a, b int) int  { return a * b } // unexported — only inside mathutils
```

Trying to call `mathutils.multiply(2, 3)` from another package:
```
cannot refer to unexported name mathutils.multiply
```

## Import syntax

```go
import "fmt"                          // single import

import (
    "fmt"

    "learning/mathutils"              // local imports separated by blank line by convention
)
```

## Special exception: `_test.go` files

The ONLY time two package names can coexist in one folder is for testing. Files ending in `_test.go` can use `package X_test` to write black-box tests that only see exported names.

## `package main`

The one package name not tied to its folder name. Used in every folder that produces an executable. Must contain a `func main()` entry point.

## Running multi-folder programs

From the module root:

```bash
go run .              # runs the root package main
go run ./arrays       # runs the arrays/ folder's program
go run ./mathutils    # FAILS — mathutils is not main
```

## Common errors

| Error | Cause |
|---|---|
| `found packages main and X in <dir>` | Two different package names in one folder |
| `expected 'package', found ...` | Missing `package X` declaration at top of file |
| `cannot refer to unexported name X.y` | Tried to call a lowercase function from another package |
| `package <path> is not a main package` | Tried `go run` on a non-main package |
| `main redeclared in this block` | Two `func main()` functions in the same folder |

## Mental model

Think of `package X` at the top of a file as a **label**. The compiler groups all files with the same label (in the same folder) into one compilation unit. That's why same-package files don't need to import each other — they're already glued together.

The unit of import in Go is the **package**, not the file.
