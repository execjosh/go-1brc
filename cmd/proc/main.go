package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"slices"
	"sync"

	"github.com/execjosh/go-1brc/internal/biginput"

	"golang.org/x/exp/maps"
)

const bufSize = 1024 * 1024 * 4
const chBufSize = 4096 * 4096

type stats struct {
	name  []byte
	min   float64
	max   float64
	count int
	sum   float64
}

type resultMap map[string]*stats

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
	results := mapChunksToResults(chunks)
	stations := reduceResults(results)

	keys := maps.Keys(stations)
	slices.Sort(keys)

	var b bytes.Buffer
	b.Grow(bufSize)
	b.WriteByte('{')
	name := keys[0]
	b.WriteString(name)
	b.WriteByte('=')
	stn := stations[name]
	fmt.Fprintf(&b, "%2.1f", stn.min)
	b.WriteByte('/')
	fmt.Fprintf(&b, "%2.1f", stn.sum/float64(stn.count))
	b.WriteByte('/')
	fmt.Fprintf(&b, "%2.1f", stn.max)
	for _, name = range keys[1:] {
		stn = stations[name]
		b.WriteByte(',')
		b.WriteByte(' ')
		b.WriteString(name)
		b.WriteByte('=')
		fmt.Fprintf(&b, "%2.1f", stn.min)
		b.WriteByte('/')
		fmt.Fprintf(&b, "%2.1f", stn.sum/float64(stn.count))
		b.WriteByte('/')
		fmt.Fprintf(&b, "%2.1f", stn.max)
	}
	b.WriteByte('}')
	b.WriteByte('\n')

	os.Stdout.Write(b.Bytes())
}

func mapChunksToResults(chunks <-chan []byte) <-chan resultMap {
	var wg sync.WaitGroup
	wg.Add(1)

	results := make(chan resultMap, chBufSize)

	go func() {
		defer close(results)
		wg.Wait()
	}()

	go func() {
		defer wg.Done()

		for chunk := range chunks {
			chunk := chunk
			wg.Add(1)
			go func() {
				defer wg.Done()

				stations := resultMap{}

				var name []byte
				for {
					for i, b := range chunk {
						if b != ';' {
							continue
						}
						name = chunk[:i]
						chunk = chunk[i+1:]
						break
					}

					var i int
					var neg bool
					if chunk[0] == '-' {
						neg = true
						chunk = chunk[1:]
						i = 1
					}
					var temp float64
					for ; i < len(chunk); i++ {
						c := chunk[i]
						if c == '.' {
							i++
							// if i >= len(chunk) {
							// 	panic(fmt.Errorf("%q: end after dot", name))
							// }
							c = chunk[i]
							i++
							// if i >= len(chunk) {
							// 	panic(fmt.Errorf("%q: end after tenths", name))
							// }
							chunk = chunk[i:]
							switch c {
							case '1':
								temp += 0.1
							case '2':
								temp += 0.2
							case '3':
								temp += 0.3
							case '4':
								temp += 0.4
							case '5':
								temp += 0.5
							case '6':
								temp += 0.6
							case '7':
								temp += 0.7
							case '8':
								temp += 0.8
							case '9':
								temp += 0.9
							}
							break
						}
						temp = temp*10 + float64(c-'0')
					}
					if neg {
						temp *= -1
					}

					stn, ok := stations[string(name)]
					if !ok {
						stn = &stats{
							name: name,
						}
						stations[string(name)] = stn
					}
					stn.count++
					stn.min = math.Min(stn.min, temp)
					stn.max = math.Max(stn.max, temp)
					stn.sum += temp

					if chunk[0] != '\n' {
						panic(fmt.Errorf("expected newline but got %02x for %s", chunk[0], name))
					}
					chunk = chunk[1:]
					if len(chunk) < 1 {
						break
					}
				}

				results <- stations
			}()
		}
	}()

	return results
}

func reduceResults(results <-chan resultMap) resultMap {
	stations := resultMap{}

	ch := make(chan *stats, chBufSize)
	go func() {
		for s := range ch {
			stn, ok := stations[string(s.name)]
			if !ok {
				stations[string(s.name)] = s
				continue
			}

			stn.count += s.count
			stn.min = math.Min(stn.min, s.min)
			stn.max = math.Max(stn.max, s.max)
			stn.sum += s.sum
		}
	}()
	defer close(ch)

	for result := range results {
		for _, s := range result {
			ch <- s
		}
	}

	return stations
}
