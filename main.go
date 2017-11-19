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

type id string

type segment struct {
	start, end int64
}

// Implement sort.Interface
// https://golang.org/pkg/sort/#Interface
type int64arr []int64

func (a int64arr) Len() int           { return len(a) }
func (a int64arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }

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
	var s segment
	var start, end float64
	segments := make(map[id]segment)
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
		s.start = int64(start)
		end, err = strconv.ParseFloat(vals[2], 64)
		if err != nil {
			panic(err)
		}
		s.end = int64(end)
		segments[id(vals[0])] = s
	}

	// Concatenate the start and end points
	var points int64arr
	for _, v := range segments {
		points = append(points, v.start, v.end)
	}
	// Sort them in ascending order
	sort.Sort(points)
	// Record the counts for each point
	var counts []int
	var count int
	for _, p := range points {
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
	for i := 0; i < len(points); i++ {
		fmt.Printf("[%d] %d\n", points[i], counts[i])
	}
}
