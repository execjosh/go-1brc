package main

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/execjosh/go-1brc/internal/biginput"
)

const bufSize = 1024 * 1024 * 2

func main() {
	cpuProfFile, err := os.Create("cpu.prof")
	if err != nil {
		panic(fmt.Errorf("creating cpu.prof: %w", err))
	}
	defer cpuProfFile.Close()
	if err := pprof.StartCPUProfile(cpuProfFile); err != nil {
		panic(fmt.Errorf("starting CPU profiling: %w", err))
	}
	defer pprof.StopCPUProfile()

	infile, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	chunks := biginput.ReadChunks(infile, bufSize)
	for range chunks {
	}
}
