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
	// fmt.Println("Number of segments:", len(segments))

	// Concatenate the start and end points
	// Use a map so as to remove duplicates
	locMap := make(map[float64]bool)
	for _, v := range segments {
		locMap[v.start] = true
		locMap[v.end] = true
	}
	// Debug
	// fmt.Println("Number of unique points:", len(locMap))

	// Turn into a slice so that you can sort it in ascending order
	var locList []float64
	for k := range locMap {
		locList = append(locList, k)
	}
	// Sort them in ascending order
	sort.Float64s(locList)
	// Debug
	// fmt.Println("Number of points in sorted list:", len(locList))

	// Record the segments that each point belongs to
	var points []*pointImpl
	for _, loc := range locList {
		p := &pointImpl{location: loc}
		for _, s := range segments {
			if loc >= s.start && loc <= s.end {
				_ = p.addTo(s.id)
			}
		}
		// Sort the segments alphabetically
		sort.Strings(p.in)
		points = append(points, p)
	}
	// Debug
	/* for i := 0; i < len(points); i++ {
		fmt.Printf("%0.f\t%d\n", points[i].location, len(points[i].belongsTo()))
	} */

	// `locList` carries the monotonically increasing list of points
	// `points` carries the same list, plus the segments that each point corresponds to

	// result includes the data we're interested in
	result := make(map[float64]map[float64][]string)
	for i := 0; i < len(points)-1; i++ {
		// Initialize the inner map
		// See: https://stackoverflow.com/a/44305711/2363529
		result[points[i].location] = make(map[float64][]string)
		// What are the segments that go over this point?
		segStart := points[i].belongsTo()
		for j := len(points) - 1; j > i; j-- {
			var overlap []string
			// What are the segments that go over this point?
			segEnd := points[j].belongsTo()
			// What is the overlap between segStart and segEnd?
			for _, seg1 := range segStart {
				for _, seg2 := range segEnd {
					if strings.Compare(seg1, seg2) == 0 {
						overlap = append(overlap, seg1)
						break // Let's move on to the next item in segStart
					}
				}
			}
			if len(overlap) > 0 {
				result[points[i].location][points[j].location] = overlap
			}
			// Now let's examine the next segEnd...
		}
	}
	// Debug
	for k1, v1 := range result {
		for k2, v2 := range v1 {
			fmt.Printf("Overlap starting from %0.f and ending at %0.f (length: %0.f) shared by %d segments: %v\n",
				k1, k2, k2-k1, len(v2), v2)
		}
	}
}
