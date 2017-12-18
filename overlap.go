package overlap

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type pointImpl struct {
	loc float64
	in  []string
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

type segmentImpl struct {
	id         string
	start, end float64
}

// Used for setting Calculate's return object.
const (
	overlapLength = iota
	overlapStart
	overlapEnd
	segmentCount
	segmentList
)

// Calculate reads a CSV file with segments and returns a slice identifying their overlaps.
// Each CSV record should identify a segment and consist of exactly three fields:
//
//		segment_id(string),segment_start(float64),segment_end(float64)
//
// The type of each field is included in parentheses.
func Calculate(filePath string, hasHeaders bool) ([][]string, error) {
	// Load a CSV file with tuples [id, start, end]
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s = %s", filePath, err)
	}

	// Read file into a slice of records
	r1 := bytes.NewReader(b)
	r2 := csv.NewReader(r1)
	records, err := r2.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read CSV file = %s", err)
	}

	// Sanitize the input
	if len(records) == 0 {
		return nil, fmt.Errorf("file is empty")
	}
	if hasHeaders {
		if len(records) == 1 {
			return nil, fmt.Errorf("file has no segments")
		}
		// Otherwise just omit the header line
		records = records[1:]
	}

	// Parse into segments
	var segment segmentImpl
	var segments []segmentImpl
	for i, record := range records {
		// Does the record have 3 fields?
		if len(record) != 3 {
			return nil, fmt.Errorf("expected 3 fields per line in row %d, got %d instead", i+1, len(record))
		}
		// The record has 3 fields as expected
		// Convert to the appropriate types
		segment.start, err = strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("could not convert element '%v' in column 2, row %d to a floating-point number = %s", record[1], i+1, err)
		}
		segment.end, err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("could not convert element '%v' in column 3, row %d to a floating-point number = %s", record[2], i+1, err)
		}
		segment.id = record[0]
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
	index := 0
	for k := range locMap {
		locList[index] = k
		index++
	}
	// Sort list in ascending order
	sort.Float64s(locList)

	// Record the segments that each point belongs to
	points := make([]*pointImpl, len(locList))
	for i, loc := range locList {
		p := &pointImpl{loc: loc}
		for _, s := range segments {
			if loc >= s.start && loc <= s.end {
				_ = p.addTo(s.id)
			}
		}
		// Sort the segments in increasing order
		sort.Strings(p.in)
		points[i] = p
	}

	// `locList` carries the monotonically increasing list of unique points
	// `points` carries the same list, plus the segments that each point corresponds to
	// `resMap` includes the data we're interested in
	resMap := make(map[float64]map[float64][]string)
	for i := 0; i < len(points)-1; i++ {
		// Initialize the inner map
		// See: https://stackoverflow.com/a/44305711/2363529
		resMap[points[i].loc] = make(map[float64][]string)
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
			if len(overlap) > 1 {
				resMap[points[i].loc][points[j].loc] = overlap
			}
			// Now let's examine the next segEnd (if any)
		}
	}

	// Return
	var results [][]string
	// Write the header
	header := []string{"overlap_length", "overlap_start", "overlap_end", "segment_count", "segment_list"}
	results = append(results, header)
	// Append the records
	result := make([]string, len(header))
	for k1, v1 := range resMap {
		for k2, v2 := range v1 {
			result[overlapLength] = strconv.FormatFloat(k2-k1, 'f', -1, 64)
			result[overlapStart] = strconv.FormatFloat(k1, 'f', -1, 64)
			result[overlapEnd] = strconv.FormatFloat(k2, 'f', -1, 64)
			result[segmentCount] = strconv.Itoa(len(v2))
			result[segmentList] = strings.Join(v2, ",")
			results = append(results, result)
		}
	}
	return results, nil
}
