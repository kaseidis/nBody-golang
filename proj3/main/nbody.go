package main

import (
	"encoding/json"
	"fmt"
	"os"
	"proj3/nbody"
	"strconv"
)

// Print usage and exit program
func PrintUsageAndExit() {
	const usage = "Usage: nbody numOfThreads\n" +
		"numOfThreads          = The thread count, 0 means run sequential version.\n"
	fmt.Fprintln(os.Stderr, usage)
	os.Exit(0)
}

// Check
func main() {
	// Get numOfThreads
	var numOfThreads int
	var err error
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid numOfThreads")
		PrintUsageAndExit()
	}
	if numOfThreads, err = strconv.Atoi(args[0]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		PrintUsageAndExit()
	}
	// Consturct json encoder/decoder from stdout/stdin
	input := json.NewDecoder(os.Stdin)
	output := json.NewEncoder(os.Stdout)
	// Call bsp and seqential version of bsp
	if numOfThreads == 0 {
		nbody.RunSeqential(input, output)
	} else {
		nbody.RunBsp(input, output, numOfThreads)
	}
}
