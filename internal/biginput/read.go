package biginput

import (
	"errors"
	"io"
)

const (
	chBufSize = 4096 * 4096
)

func ReadChunks(r io.Reader, bufSize int) <-chan []byte {
	ch := make(chan []byte, chBufSize)

	go func() {
		defer close(ch)

		var remainder []byte
		var end int

		for {
			buf := make([]byte, bufSize)
			copy(buf, remainder)

			n, err := r.Read(buf[len(remainder):])
			if err != nil {
				// anything other than EOF is unexpected so bail
				if !errors.Is(err, io.EOF) {
					panic(err)
				}
			}

			end = len(remainder) + n

			// give up if we didn't actually read anything or have anything left
			if end < 1 {
				break
			}

			// find the last newline
			for i := end - 1; i >= 0; i-- {
				if buf[i] == '\n' {
					ch <- buf[:i+1]            // keep the newline
					remainder = buf[i+1 : end] // keep everything after the final newline
					break
				}
			}
		}
	}()

	return ch
}
