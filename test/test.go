package main

import (
	"fmt"

	"github.com/Kitsunetic/recursivelist"
)

func main() {
	files, errs, done := recursivelist.ListRecurrent("./")
L:
	for {
		select {
		case file := <-files:
			fmt.Println(file)
		case err := <-errs:
			fmt.Println(err)
		case <-done:
			break L
		}
	}
}
