package overlap

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type pointImpl struct {
	loc float64
	in  []int
}

type point interface {
	addTo(segmentID int) int
	belongsTo() []int
}

func (p *pointImpl) addTo(segmentID int) int {
	p.in = append(p.in, segmentID)
	return len(p.in)
}

func (p *pointImpl) belongsTo() []int {
	return p.in
}

type segmentImpl struct {
	id         int
	start, end float64
}

// Result describes an overlap. It carries its length, its start
// and end points, the number of segments that go over it, as well
// as their IDs.
type Result struct {
	OverlapLength            float64
	OverlapStart, OverlapEnd float64
	SegmentCount             int
	SegmentList              []int
}

// Calculate reads a file with segments and returns a slice of their overlaps.
// Each line (segment) in the input file should follow the format:
//		segment_id(int) segment_start(float64) segment_end(float64)
// The fields are tab-separated. Their types are listed in parentheses.
func Calculate(filePath string) ([]Result, error) {
	// Load a file with tuples [id, start, end]
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %s", filePath, err)
	}

	// Read file into slice
	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	var count int
	var line string
	var segment segmentImpl
	var segments []segmentImpl
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		line = scanner.Text()
		// We now have a line
		count++
		// Does it have 3 fields?
		vals := strings.Split(line, "\t")
		if len(vals) != 3 {
			return nil, fmt.Errorf("expected 3 fields per line, got %d instead", len(vals))
		}
		// The line has 3 fields as expected
		// Convert to the appropriate types
		segment.id, err = strconv.Atoi(vals[0])
		if err != nil {
			return nil, fmt.Errorf("could not convert element %v in column 1, row %d to an integer: %s", vals[0], count, err)
		}
		segment.start, err = strconv.ParseFloat(vals[1], 64)
		if err != nil {
			return nil, fmt.Errorf("could not convert element %v in column 2, row %d to a floating-point number: %s", vals[1], count, err)
		}
		segment.end, err = strconv.ParseFloat(vals[2], 64)
		if err != nil {
			return nil, fmt.Errorf("could not convert element %v in column 3, row %d to a floating-point number: %s", vals[2], count, err)
		}
		segments = append(segments, segment)
	}

	// Create list of unique points
	// Use a map to remove duplicates
	locMap := make(map[float64]bool)
	for _, v := range segments {
		locMap[v.start] = true
		locMap[v.end] = true
	}

	// Turn into slice so that you can sort in ascending order
	locList := make([]float64, len(locMap))
	for k := range locMap {
		locList = append(locList, k)
	}
	// Sort list in ascending order
	sort.Float64s(locList)

	// Record the segments that each point belongs to
	points := make([]*pointImpl, len(locList))
	for _, loc := range locList {
		p := &pointImpl{loc: loc}
		for _, s := range segments {
			if loc >= s.start && loc <= s.end {
				_ = p.addTo(s.id)
			}
		}
		// Sort the segments in ascending order
		sort.Ints(p.in)
		points = append(points, p)
	}

	// `locList` carries the monotonically increasing list of unique points
	// `points` carries the same list, plus the segments that each point corresponds to
	// `resMap` includes the data we're interested in
	resMap := make(map[float64]map[float64][]int)
	for i := 0; i < len(points)-1; i++ {
		// Initialize the inner map
		// See: https://stackoverflow.com/a/44305711/2363529
		resMap[points[i].loc] = make(map[float64][]int)
		// What are the segments that go over this point?
		segStart := points[i].belongsTo()
		for j := len(points) - 1; j > i; j-- {
			var overlap []int
			// What are the segments that go over this point?
			segEnd := points[j].belongsTo()
			// What is the overlap between segStart and segEnd?
			for _, seg1 := range segStart {
				for _, seg2 := range segEnd {
					if seg1 == seg2 {
						overlap = append(overlap, seg1)
						break // Let's move on to the next item in segStart
					}
				}
			}
			if len(overlap) > 1 {
				resMap[points[i].loc][points[j].loc] = overlap
			}
			// Now let's examine the next segEnd (if any)
		}
	}

	// Return
	var results []Result
	for k1, v1 := range resMap {
		for k2, v2 := range v1 {
			// Line format:
			// overlap length (years) - overlap start - overlap end - segment count - segment list
			r := Result{
				OverlapLength: k2 - k1,
				OverlapStart:  k1,
				OverlapEnd:    k2,
				SegmentCount:  len(v2),
				SegmentList:   v2,
			}
			results = append(results, r)
		}
	}

	return results, nil
}
