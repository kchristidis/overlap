package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/kchristidis/overlap"
)

func main() {
	// Parse the command-line arguments
	args := os.Args
	inFile := args[1]
	outFile := "out_" + args[1] // allows us to keep the file extension
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
	w := csv.NewWriter(f)
	if err := w.WriteAll(results); err != nil { // calls Flush internally
		log.Fatal(err)
	}
	log.Println("Wrote: ", f.Name())
}
