package prefixwriter

import (
	"bytes"
	"io"
	"sync"
)

type prefixWriter struct {
	prefix string
	w      io.Writer

	mu      sync.Mutex
	atStart bool // we are at the start of a line
}

// New creates a writer that prepends a prefix to every line it writes.
// It is safe to use the writer from multiple goroutines.
func New(prefix string, w io.Writer) io.Writer {
	return &prefixWriter{
		w:       w,
		atStart: true,
		prefix:  prefix,
	}
}

func (w *prefixWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	segments := bytes.Split(p, []byte("\n"))
	for i, s := range segments {
		if len(s) > 0 {
			if w.atStart {
				// write the prefix if at start of a line
				_, err = w.w.Write([]byte(w.prefix))
				if err != nil {
					return
				}
			}
			_, err = w.w.Write(s)
			if err != nil {
				return
			}
			w.atStart = false
		} else {
			// If segment is empty, we're at start of a line
			w.atStart = true
		}

		if i < (len(segments) - 1) {
			// If not at the end of the segments, write a newline
			_, err = w.w.Write([]byte("\n"))
			if err != nil {
				return
			}
			w.atStart = true
		}
	}
	n = len(p)
	return
}
