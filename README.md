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
	prog.NewProgram(context.Background(), func() {
		log.Println("Main")
	}).Run()

    // with init and stop
	prog.NewProgram(context.Background(), func() error {
		log.Println("Main")
		return nil
	}).Init(func() error {
		log.Println("Init")
		return nil
	}).Stop(func() {
		log.Println("Stop")
	}).Run()
}
```