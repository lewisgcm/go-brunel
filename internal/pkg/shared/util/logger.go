/*
 * Author: Lewis Maitland
 *
 * Copyright (c) 2019 Lewis Maitland
 */

package util

import (
	"strings"
)

// Logger should log the provided string, returning an error on failure
type Logger func(log string) error

// LoggerWriter writes lines of log messages to the Logger, it attempts to split lines and log only 1 line at a time
type LoggerWriter struct {
	Recorder Logger
	leftOver string
}

// Write will attempt to write any full lines using the Logger. It will try and mediate
// the line ending for cross platform support.
func (w *LoggerWriter) Write(p []byte) (int, error) {
	s := w.leftOver + strings.Replace(
		strings.Replace(string(p), "\r\n", "\n", -1),
		"\r",
		"\n",
		-1,
	)

	lines := strings.Split(s, "\n")
	for i := 0; i < len(lines)-1; i++ {
		err := w.Recorder(lines[i])
		if err != nil {
			return len(p), err
		}
	}
	w.leftOver = lines[len(lines)-1]

	return len(p), nil
}

// Close will flush any remaining data to the Logger
func (w *LoggerWriter) Close() error {
	if w.leftOver != "" {
		err := w.Recorder(w.leftOver)
		if err != nil {
			return err
		}
	}
	return nil
}
