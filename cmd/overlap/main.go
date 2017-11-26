package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/kchristidis/overlap"
)

func main() {
	// Parse the command-line arguments
	args := os.Args
	inFile := args[1]
	outFile := args[1] + ".overlap"
	if len(args) >= 3 {
		outFile = args[2]
	}
	// Use the library
	results, err := overlap.Calculate(inFile)
	if err != nil {
		log.Fatal(err)
	}
	// Write results to file
	f, err := os.Create(outFile)
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
	log.Println("Wrote: ", f.Name())
}
