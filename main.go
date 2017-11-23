package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type segmentImpl struct {
	id         int
	start, end float64
}

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

func main() {
	// Create a logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// Load a file with tuples [id, start, end]
	b, err := ioutil.ReadFile("unix_tuples.txt")
	if err != nil {
		panic(err)
	}

	// Read file into slice
	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	var line string
	var segment segmentImpl
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
		segment.id, err = strconv.Atoi(vals[0])
		if err != nil {
			panic(err)
		}
		segment.start, err = strconv.ParseFloat(vals[1], 64)
		if err != nil {
			panic(err)
		}
		segment.end, err = strconv.ParseFloat(vals[2], 64)
		if err != nil {
			panic(err)
		}
		segments = append(segments, segment)
	}
	sugar.Debug("Number of segments in file: ", len(segments))

	// Create list of unique points
	// Use a map so as to remove duplicates
	locMap := make(map[float64]bool)
	for _, v := range segments {
		locMap[v.start] = true
		locMap[v.end] = true
	}
	sugar.Debug("Number of unique points: ", len(locMap))

	// Turn into a slice so that you can sort it in ascending order
	var locList []float64
	for k := range locMap {
		locList = append(locList, k)
	}
	// Sort list in ascending order
	sort.Float64s(locList)

	// For each point, record the segments that it belongs to
	var points []*pointImpl
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
	// `result` includes the data we're interested in
	result := make(map[float64]map[float64][]int)
	for i := 0; i < len(points)-1; i++ {
		// Initialize the inner map
		// See: https://stackoverflow.com/a/44305711/2363529
		result[points[i].loc] = make(map[float64][]int)
		// What are the segments that go over this point?
		segStart := points[i].belongsTo()
		sugar.Debugf("Examining point %0.f which belongs to segments: %v\n", points[i].loc, segStart)
		for j := len(points) - 1; j > i; j-- {
			var overlap []int
			// What are the segments that go over this point?
			segEnd := points[j].belongsTo()
			sugar.Debugf("\tComparing with point %0.f which belongs to segments: %v\n", points[j].loc, segEnd)
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
				result[points[i].loc][points[j].loc] = overlap
				sugar.Debugf("\t\tThe two points overlap in: %v\n", overlap)
			} else {
				sugar.Debugf("\t\tThe two points do not overlap\n")
			}
			// Now let's examine the next segEnd (if any)
		}
	}

	for k1, v1 := range result {
		for k2, v2 := range v1 {
			sugar.Debugf("Overlap starting from %0.f (%s) and ending at %0.f (%s) (length: %.1fd) shared by %d segments: %v\n",
				k1, time.Unix(int64(k1), 0), k2, time.Unix(int64(k2), 0), (k2-k1)/60/60/24, len(v2), v2)
		}
	}
}
