package option

import (
	"io"
	"os"
)

// DebugOptions defines options for debug logging.
type DebugOptions struct {
	Enabled bool
	Time    bool
	Writer  io.Writer
}

// WithDebug enables debug logging with an optional name.
// This option applies only to sqlb builders who print built queries in debug mode.
func WithDebug(writer ...io.Writer) Option {
	return func(o *Options) {
		var w io.Writer
		switch len(writer) {
		case 0:
			w = os.Stdout
		case 1:
			w = writer[0]
		default:
			w = io.MultiWriter(writer...)
		}
		o.Debug.Enabled = true
		o.Debug.Writer = w
	}
}

// WithDebugTime enables debug logging with time measurement.
func WithDebugTime(writer ...io.Writer) Option {
	return func(o *Options) {
		WithDebug(writer...)(o)
		o.Debug.Time = true
	}
}
