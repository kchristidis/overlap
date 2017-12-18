package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"

	"github.com/kchristidis/overlap"
)

func main() {
	// Define flag
	headers := flag.Bool("headers", false, "Does the input file have headers?")
	// Parse the flag
	flag.Parse()
	// Parse the command-line arguments trailing the flag.
	args := flag.Args()
	inFile := args[0]
	outFile := "out_" + args[0] // allows us to keep the file extension
	if len(args) >= 2 {
		outFile = args[1]
	}
	// Use the library
	results, err := overlap.Calculate(inFile, *headers)
	if err != nil {
		log.Fatal(err)
	}
	// Write results to file
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	if err := w.WriteAll(results); err != nil { // calls Flush internally
		log.Fatal(err)
	}
	log.Println("Wrote: ", f.Name())
}
