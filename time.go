package zaptextenc

import (
	"bytes"
	"time"
)

// LayoutTimeFormatterOption is an option for log entry time formatting according to specified layout
type LayoutTimeFormatterOption struct {
	layout string
}

// Apply option to encoder
func (f LayoutTimeFormatterOption) Apply(e *Encoder) {
	e.setTimeFormatter(f.formatter())
}

func (f LayoutTimeFormatterOption) formatter() TimeFormatter {
	return func(w *bytes.Buffer, t time.Time) {
		w.WriteString(t.Local().Format(f.layout))
		w.Write([]byte{' '})
	}
}

// ShortTime format
func ShortTime() Option {
	return &LayoutTimeFormatterOption{"15:04:05.000000"}
}

// RFC3339Time format
func RFC3339Time() Option {
	return &LayoutTimeFormatterOption{time.RFC3339}
}

// NoTimeFormatterOption defines an option which makes encoder to skip log level
type NoTimeFormatterOption struct{}

// Apply sets level formatter for an encoder
func (f NoTimeFormatterOption) Apply(e *Encoder) {
	e.setTimeFormatter(func(_ *bytes.Buffer, _ time.Time) {})
}

// NoTime skips log entry time
func NoTime() Option {
	return NoTimeFormatterOption{}
}
