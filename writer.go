package main

import "io"

type maxSizeWriter struct {
	Writer  io.Writer
	maxSize int
	written int
}

func (w *maxSizeWriter) Write(p []byte) (n int, err error) {
	if w.written+len(p) > w.maxSize {
		remaining := w.maxSize - w.written
		if remaining > 0 {
			n, err = w.Writer.Write(p[:remaining])
			w.written += n
		}
		return n, io.ErrShortWrite
	}
	n, err = w.Writer.Write(p)
	w.written += n
	return n, err
}
