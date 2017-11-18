package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	// Read a file with tuples [id, start, end]
	f, err := ioutil.ReadFile("unix_tuples.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", f)
	println()
}
