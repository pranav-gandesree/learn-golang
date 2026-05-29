# Go Learning Notes

Quick-reference docs covering Go fundamentals I've worked through. Each file is a self-contained reference for one topic.

| # | Topic | Concepts |
|---|---|---|
| [01](01-packages.md) | Packages | one folder = one package, imports, exported/unexported names |
| [02](02-arrays-and-slices.md) | Arrays & Slices | fixed-size arrays vs dynamic slices, `len`/`cap`, `append`, `make`, slice growth |
| [03](03-maps.md) | Maps | declaration, comma-ok idiom, deletion, random iteration order, reference types |
| [04](04-goroutines.md) | Goroutines & WaitGroups | concurrent functions, `go` keyword, synchronization |
| [05](05-mutexes.md) | Mutexes | data races, `sync.Mutex`, `sync.RWMutex`, race detector |
| [06](06-channels.md) | Channels | unbuffered/buffered, close + range, direction, `select`, worker pools |

## Companion code

Each topic has a runnable demo in its own folder:

```
../arrays/         # arrays + slices basics
../slicegrowth/    # visualizing slice capacity doubling
../maps/           # map operations
../goroutines/     # goroutines + WaitGroup
../mutexes/        # data races + mutexes
../channels/       # channels, select, worker pool
../mathutils/      # imported by ../main.go to demo packages
```

Run any of them from the `learning/` directory:

```bash
go run ./maps
go run ./goroutines
go run -race ./mutexes   # with the race detector
```
