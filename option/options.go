package option

import (
	"os"
)

// Options defines options for scanning.
type Options struct {
	Debug  DebugOptions
	Select SelectOptions
}

// Option defines a function type for setting Options.
type Option func(*Options)

// New creates a new Options with the given Option functions applied.
func New(opts ...Option) *Options {
	options := &Options{
		Debug: DebugOptions{
			Writer: os.Stdout,
		},
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
