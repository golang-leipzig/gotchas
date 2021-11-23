# gotchas

Go gotchas, surprises, puzzles.

Other lists and articles:

* [50 Shades of Go: Traps, Gotchas, and Common Mistakes for New Golang Devs](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)
* [Do you make these Go coding mistakes?](https://yourbasic.org/golang/gotcha/)
* [Common Gotchas in Go](https://deadbeef.me/2018/01/go-gotchas)
* [Dark Corners of Go](https://rytisbiel.com/2021/03/06/darker-corners-of-go/)
* [Gotchas in the Go Network Packages Defaults](https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/)

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

We do not know exactly, but maybe the semantics around the limit are less clear
when writing - or there [has not been enough
demand](https://groups.google.com/d/msg/golang-nuts/FB4OtiRm28E/RI9-GHlecr4J).

Kubernetes has an
[ioutils](https://godoc.org/k8s.io/kubernetes/pkg/kubelet/util/ioutils)
package, which contains
a [LimitWriter](https://github.com/kubernetes/kubernetes/blob/579e0c74c150085b3fac01f6a33b66db96922f93/pkg/kubelet/util/ioutils/ioutils.go#L39-L70). Minio has a [LimitWriter](https://github.com/minio/minio/blob/dfadf70a7f418798feab0cf013a7765f1adb21f5/pkg/ioutil/ioutil.go#L70-L113) as well.

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

## Returning values from defer requires named result parameters

In Go we often defer function calls, e.g. `defer f.Close()` after we opened a file or a `defer tx.Rollback()` to rollback a failed SQL transaction.
But what if you want to return an error or result from a deferred function call?
Consider this example where we want to add additional context to an error.  One might think that we can simply assign to `err` since it is in the closure of the deferred function.

```go
func g() error {
    err := fmt.Errorf("new error")
    defer func() {
        err = fmt.Errorf("more context: %w", err)
    }()
    return err
}
```

Surprisingly, this will return the unmodified error.

[Effective Go#Recover](https://golang.org/doc/effective_go.html#recover) explains a simple panic recovery mechanism that points out what we are doing wrong:

> ... deferred functions can modify named return values.

After reading the section about [Defer statements in Go's language specification](https://golang.org/ref/spec#Defer_statements) it becomes pretty clear why we have to a named result parameter:

> ... deferred functions are executed after any result parameters are set by that return statement but before the function returns to its caller.

Here's a working example ([playground](https://play.golang.org/p/Sumi3P-bsUO)):

```go
func f() (err error) {
    err = fmt.Errorf("new error")
    defer func() {
        err = fmt.Errorf("more context: %w", err)
    }()
    return err
}
```

## Dot or underscore prefixed files

The go tools ignore files starting with either a `.` or an `_`, e.g. `_main.go` would not be considered by the Go tools.
Folders called `testdata` are also ignored.  Here is the statement from Go's documentation:

> Directory and file names that begin with "." or "_" are ignored by the go tool, as are directories named "testdata".

[source](https://golang.org/cmd/go/#hdr-Package_lists_and_patterns)

Note that `embed`, which was included with Go 1.16, also follows the behavior described above when determining which files to include.

> If a pattern names a directory, all files in the subtree rooted at that directory are embedded (recursively), except that files with names beginning with ‘.’ or ‘_’ are excluded.

[source](https://golang.org/pkg/embed/#hdr-Directives)

## `filepath.Clean` only works on absolute paths

Let's say you implement a file server and thus need to map a request URL path to local path, e.g. `https://<some-domain>/some/file` should serve `/<some-base-directory>/some/file`.  Without cleaning up the request path you will be open to [path-traversal attacks](https://owasp.org/www-community/attacks/Path_Traversal).  That means if someone requests `https://<some-domain>/../../../etc/passwd` your server will happily respond with `passwd` (if it has permission to read it).  To prevent this problem you can just do a `filepath.Clean(req.URL.Path)`, **but**, this will not work if `req.URL.Path` is not absolute.  So, when using `filepath.Clean` make sure that the path you're passing into is absolute, when in doubt just add a `/` in front.

[playground](https://play.golang.org/p/Q-qKnLVrh0P)

## `http.Client` is following up to 10 redirects by default

It was surprising to me that `http.Client` and `http.DefaultClient` are following up to 10 redirects by default.
I discovered this when doing a `HEAD` request on a resource that responds with a redirect, giving a 302 status code,
and a `Location` header with the response.  But, with the `DefaultClient` you will receive a `200` instead, and no `Location` header of course.

You can override the default behavior by setting a custom [`CheckRedirect` function in the client](https://pkg.go.dev/net/http#Client) like this:

```go
httpCl := http.DefaultClient
httpCl.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
        return http.ErrUseLastResponse
}
```

Note that [`ErrUseLastResponse`](https://pkg.go.dev/net/http#ErrUseLastResponse) will just signal the client to stop following any more redirects.

## `recover` won't catch panics in spawned goroutines

This one is subtle and caused many servers to crash.
One thing in advance, using a recovery middleware in your web server **won't** safe you!

The following function will not recover and cause your program to crash:

```go
func wontRecover() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recovered from panic:", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		panic("panic'ed in goroutine")
	}()

	fmt.Println("waiting for goroutine to finish")
	wg.Wait()
}
```

Imagine `wontRecover` to be an `http.Handler` and you will see how widespread this problem is.
The workaround is to `recover()` in any subroutine that might panic.
