# gotchas

Go gotchas, surprises, puzzles.

Other lists and articles:

* [50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)
* [Do you make these Go coding mistakes?](https://yourbasic.org/golang/gotcha/)
* [Common Gotchas in Go](https://deadbeef.me/2018/01/go-gotchas)

# Other gotchas and questions

## Embedded struct literal

The JSON unmarshaler works with embedded struct, while literal initialization will not.

* [https://play.golang.org/p/AeB6d-5HXrE](https://play.golang.org/p/AeB6d-5HXrE)

```
prog.go:31:13: cannot use promoted field A.A1 in struct literal of type B
```

## Missing high index causes out of range

The high index of a slice is not what it seems.

Question: Does the following snippet compile? Run? Panic?

```go
package main

func main() {
    v := [6]int{0, 1, 2, 3, 4, 5}
    w := v[:]

    w = w[:4]
    w = w[:0]

    // Question: Before you hit run! What do the following two lines result in?
    w = w[1:]
    w = w[1:3]
    // fmt.Printf("%v - len: %d, cap: %d\n", w, len(w), cap(w))
}
```

Try it out yourself:

* [https://play.golang.org/p/m-6IW67eJh2](https://play.golang.org/p/m-6IW67eJh2)

```
panic: runtime error: slice bounds out of range
```

## A struct pointer return value, that looks like nil, but is not

The zero value for a pointer is nil.

> Each element of such a variable or value is set to the zero value for its
> type: false for booleans, 0 for numeric types, "" for strings, and nil for
> pointers, functions, interfaces, slices, channels, and maps. -- [Zero value](https://golang.org/ref/spec#The_zero_value)

If the return value is pointer to a struct, e.g. `*E` and we return `nil`, we
do not actually return `nil`.

References:

* [Hiding nil values, understanding why golang fails here](https://stackoverflow.com/questions/29138591/hiding-nil-values-understanding-why-golang-fails-here)
* [https://play.golang.org/p/fj24_qUdm45](https://play.golang.org/p/fj24_qUdm45)

TBC.

```go
// go run main.go
//
// 2019/08/26 17:19:28 (*main.E)(nil): some error message
// exit status 1
//
package main

import (
        "log"
)

type E struct{}

func (e *E) Error() string {
        return "some error message"
}

func mayFail(f float32) *E {
        if f < 0.5 {
                return nil
        } else {
                return &E{}
        }
}

func main() {
        var err error
        err = mayFail(0.4)
        if err != nil {
                log.Fatalf("%#v: %s", err, err.Error())
        }
}
```
