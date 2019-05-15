# gotchas

Go gotchas, surprises, puzzles.

Obligatory readings:

* [50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)

# Other gotchas and questions

## Embedded struct literal

The JSON unmarshaler works with embedded struct, while literal initialization will not.

* [https://play.golang.org/p/AeB6d-5HXrE](https://play.golang.org/p/AeB6d-5HXrE)

```
prog.go:31:13: cannot use promoted field A.A1 in struct literal of type B
```

## Missing high index causes out of range

The high index of a slice is not what it seems.

* [https://play.golang.org/p/m-6IW67eJh2](https://play.golang.org/p/m-6IW67eJh2)

```
panic: runtime error: slice bounds out of range
```
