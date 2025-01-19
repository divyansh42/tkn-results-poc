package cli

import "io"

type Stream struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
