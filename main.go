package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type segment struct {
	id         string
	start, end float64
}

type pointImpl struct {
	location float64
	in       []string
}

type point interface {
	belongsTo() []string
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
	var start, end float64
	var segments []segment
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
		start, err = strconv.ParseFloat(vals[1], 64)
		if err != nil {
			panic(err)
		}
		end, err = strconv.ParseFloat(vals[2], 64)
		if err != nil {
			panic(err)
		}
		segments = append(segments, segment{id: vals[0], start: start, end: end})
	}

	// Concatenate the start and end points
	// Use a map so as to remove duplicates
	pointsMap := make(map[float64]bool)
	for _, v := range segments {
		pointsMap[v.start] = true
		pointsMap[v.end] = true
	}
	// Turn into a slice so that you can sort it in ascending order.
	var pointsList []float64
	for k := range pointsMap {
		pointsList = append(pointsList, k)
	}
	// Sort them in ascending order
	sort.Float64s(pointsList)
	// Record the counts for each point
	var counts []int
	var count int
	for _, p := range pointsList {
		count = 0
		// Given a point, let's figure out the segments it belongs to
		for _, s := range segments {
			if p >= s.start && p <= s.end {
				count++
			}
		}
		counts = append(counts, count)
	}

	// Debug
	for i := 0; i < len(pointsList); i++ {
		fmt.Printf("[%f] %d\n", pointsList[i], counts[i])
	}
	fmt.Printf("\nNumber of points overall: %d\nNumber of unique points: %d\n", 2*len(segments), len(pointsMap))
}
