// +build !solution

package lab2

import (
	"io"
)

/*
Task 3: Rot 13

This task is taken from http://tour.golang.org.

A common pattern is an io.Reader that wraps another io.Reader, modifying the
stream in some way.

For example, the gzip.NewReader function takes an io.Reader (a stream of
compressed data) and returns a *gzip.Reader that also implements io.Reader (a
stream of the decompressed data).

Implement a rot13Reader that implements io.Reader and reads from an io.Reader,
modifying the stream by applying the rot13 substitution cipher to all
alphabetical characters.

The rot13Reader type is provided for you. Make it an io.Reader by implementing
its Read method.
*/

type rot13Reader struct {
	r io.Reader
}

func rot13Decrypt(b byte) byte {
	if !(b >= 65 && b <= 90) && !(b >= 97 && b <= 122) {
		return b
	}

	newPos := (b - 13)
	if b >= 65 && b <= 90 {
		if newPos < 65 {
			newPos = 90 - (64 - newPos)
		}
	} else {
		if newPos < 97 {
			newPos = 122 - (96 - newPos)
		}
	}

	return newPos
}

func (r rot13Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	for i := 0; i < n; i++ {
		p[i] = rot13Decrypt(p[i])
	}
	return
}
