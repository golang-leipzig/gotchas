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
