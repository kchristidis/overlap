package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type id string

type interval struct {
	start, end float64
}

func main() {
	// Read a file with tuples [id, start, end]
	b, err := ioutil.ReadFile("unix_tuples.txt")
	if err != nil {
		panic(err)
	}

	// Format the input
	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	var line string
	var elem interval
	table := make(map[id]interval)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		line = scanner.Text()
		// We now have a line
		// Does it have 3 fields?
		vals := strings.Split(line, "\t")
		if len(vals) != 3 {
			panic(errors.New("Expected 3 fields per line"))
		}
		// The line has 3 fields as expected
		// Let's convert to the appropriate types
		elem.start, err = strconv.ParseFloat(vals[1], 64)
		if err != nil {
			panic(err)
		}
		elem.end, err = strconv.ParseFloat(vals[2], 64)
		if err != nil {
			panic(err)
		}
		table[id(vals[0])] = elem
	}

	// Debug
	var count int
	for k := range table {
		count++
		fmt.Println(k)
	}
	fmt.Printf("Number of lines: %d\n", count)
}
