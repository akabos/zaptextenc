package zaptextenc

import (
	"bytes"
	"io"
)

//
// Simple
// ======

// SimpleMessageFormatterOption is an option for simple log message format
type SimpleMessageFormatterOption struct{}

// Apply option to encoder
func (f SimpleMessageFormatterOption) Apply(e *Encoder) {
	e.setMessageFormatter(f.formatter())
}

func (f SimpleMessageFormatterOption) formatter() MessageFormatter {
	return func(w io.Writer, message string) {
		io.WriteString(w, message)
		w.Write([]byte{' '})
	}
}

// SimpleMessage formats log message in most dull and unfancy way possible
func SimpleMessage() Option {
	return &SimpleMessageFormatterOption{}
}

//
// Fixed width
// ===========

// FixedWidthMessageFormatterOption is an option for fixed width log message format
type FixedWidthMessageFormatterOption struct {
	width int
}

// Apply option to encoder
func (f FixedWidthMessageFormatterOption) Apply(e *Encoder) {
	e.setMessageFormatter(f.formatter())
}

func (f FixedWidthMessageFormatterOption) formatter() MessageFormatter {
	return func(w io.Writer, message string) {
		result := make([]byte, f.width)
		if len(message) < f.width {
			extra := f.width - len(message)
			result = append(result, []byte(message)...)
			result = append(result, bytes.Repeat([]byte(" "), extra)...)
		} else {
			result = append(result, []byte(message[:f.width])...)
		}
		w.Write(result)
		w.Write([]byte{' '})
	}
}

// FixedWidthMessage returns option for trimming/expanding log message to specified width
func FixedWidthMessage(width int) Option {
	return &FixedWidthMessageFormatterOption{width}
}
