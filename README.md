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

A similar example: [play.golang.org/p/ZfY7AN687ah](https://play.golang.org/p/ZfY7AN687ah)

## Literal initialization of promoted fields

> cannot use promoted field __ in struct literal of type

Example: One GitHub API (v3) wrapper defines different option types. The
[SearchOptions](https://godoc.org/github.com/google/go-github/github#SearchOptions)
embed
a [ListOptions](https://godoc.org/github.com/google/go-github/github#ListOptions)
type for pagination.

The following would not work: "cannot use promoted field ..."

```go
opt := &github.SearchOptions{Sort: "stars", PerPage: 10}
```

Workaround is to create an options value, then assign to the promoted field.

```go
opt := &github.SearchOptions{Sort: "stars"}
opt.PerPage = 10
```

## Three-star programming in Go

There is the concept of a star in programming languages and the concept of
a programmer who uses them, e.g.
[ThreeStarProgrammer](https://wiki.c2.com/?ThreeStarProgrammer).

It's easy to get to three stars in Go using arithmetic and indirection.

```go
package main

import (
	"fmt"
)

func main() {
	var (
		a = new(int)
		b = &a
	)
	**b = 3
	fmt.Printf("%d", *a***b)
}
```

Can you see the result? If not, just [try it out](https://play.golang.org/p/G06tzZ2mJAO)!

## Why is there no io.LimitWriter?

We do not know exactly, but maybe the semantics around the limit are less clear when writing.

Kubernetes has an
[ioutils](https://godoc.org/k8s.io/kubernetes/pkg/kubelet/util/ioutils)
package, which contains
a [LimitReader](https://github.com/kubernetes/kubernetes/blob/579e0c74c150085b3fac01f6a33b66db96922f93/pkg/kubelet/util/ioutils/ioutils.go#L39-L70).

Illustrating the point above, currently an `ErrShortWrite` will be returned, if
the limit is hit (in some previous version it was called `ErrMaximumWrite`).
However, we could just stop writing without returning an error, so `io.Copy` would work:

* [https://play.golang.org/p/LIs3CdFi59H](https://play.golang.org/p/LIs3CdFi59H)

```go
// Branch off a limited number of bytes from resp.Body into a buffer.
var (
    buf bytes.Buffer
    tee = io.TeeReader(resp.Body, LimitWriter(&buf, 512))
)
if _, err := io.Copy(ioutil.Discard, tee); err != nil {
	log.Fatal(err)
}
```

