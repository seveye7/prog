# prog
go program runner


``` go
package main

import (
	"context"
	"log"

	"github.com/seveye7/prog"
)

func main() {
    // simple
	prog.NewProgram(
		context.Background(),
		prog.WithMainFunc(func() error {
			log.Println("hello world")
			return nil
		}),
	).Run()

    // with init and stop
	prog.NewProgram(
		context.Background(),
		// main func
		func() error {
			log.Println("hello world")
			return nil
		},
		// init func
		prog.WithInit(func() error {
			log.Println("init")
			return nil
		}),
		// stop func
		prog.WithStopFunc(func() {
			log.Println("stop")
		}),
	).Run()
}
```