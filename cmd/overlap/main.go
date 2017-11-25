package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kchristidis/overlap"
)

func main() {
	// Use the library
	results, err := overlap.Calculate("unix_tuples.txt")
	if err != nil {
		log.Fatal(err)
	}
	// Write results to file
	f, err := os.Create("results.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, v := range results {
		fmt.Fprintf(w, "%.2f\t%0.f\t%0.f\t%d\t%v\n",
			v.OverlapLength/(60*60*8760), // Convert to years
			v.OverlapStart, v.OverlapEnd,
			v.SegmentCount, v.SegmentList)
	}
	w.Flush()
}
