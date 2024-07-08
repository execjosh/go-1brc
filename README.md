# go-1brc

This is my attempt at the [1BRC] to see what I can come up with, but in Go.

## time go run ./cmd/gen

This command generates a `measurements.txt` file.

## time go run ./cmd/justread

This command just reads `measurements.txt`.  I use this to get a baseline of how
quickly the read should be.

### TODO

- [ ] Make chunk size a parameter

## time go run ./cmd/proc

This is the main implementation.  It processes `measurements.txt`.

### TODO

- [ ] Make chunk size a parameter


[1BRC]: https://github.com/gunnarmorling/1brc
