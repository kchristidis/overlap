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

type segmentImpl struct {
	id         string
	start, end float64
}

type pointImpl struct {
	location float64
	in       []string
}

type point interface {
	addTo(segmentID string) int
	belongsTo() []string
}

func (p *pointImpl) addTo(segmentID string) int {
	p.in = append(p.in, segmentID)
	return len(p.in)
}

func (p *pointImpl) belongsTo() []string {
	return p.in
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
	var segments []segmentImpl
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
		segments = append(segments, segmentImpl{id: vals[0], start: start, end: end})
	}
	// Debug
	fmt.Println("Number of segments:", len(segments))

	// Concatenate the start and end points
	// Use a map so as to remove duplicates
	locMap := make(map[float64]bool)
	for _, v := range segments {
		locMap[v.start] = true
		locMap[v.end] = true
	}
	// Debug
	fmt.Println("Number of unique points:", len(locMap))

	// Turn into a slice so that you can sort it in ascending order
	var locList []float64
	for k := range locMap {
		locList = append(locList, k)
	}
	// Sort them in ascending order
	sort.Float64s(locList)
	// Debug
	fmt.Println("Number of points in sorted list:", len(locList))

	// Record the segments that each point belongs to
	var points []*pointImpl
	for _, loc := range locList {
		p := &pointImpl{location: loc}
		for _, s := range segments {
			if loc >= s.start && loc <= s.end {
				_ = p.addTo(s.id)
			}
		}
		points = append(points, p)
	}
	// Debug
	for i := 0; i < len(points); i++ {
		fmt.Printf("[%.0f] %d\n", points[i].location, len(points[i].belongsTo()))
	}
}
