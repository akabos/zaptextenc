package zaptextenc

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/uber-go/zap"
)

var pool = sync.Pool{New: func() interface{} {
	return &Encoder{
		buffer: bytes.NewBuffer([]byte{}),
	}
}}

// Option is an Encoder option
type Option interface {
	Apply(*Encoder)
}

// TimeFormatter is a func used to render log entry time
type TimeFormatter func(*bytes.Buffer, time.Time)

// LevelFormatter is a func used to render log entry level
type LevelFormatter func(*bytes.Buffer, zap.Level)

// MessageFormatter is a func used to render log entry message
type MessageFormatter func(*bytes.Buffer, string)

// Encoder is zap.Encoder implementation that writes plain text messages
type Encoder struct {
	timeF    TimeFormatter
	levelF   LevelFormatter
	messageF MessageFormatter
	buffer   *bytes.Buffer
}

// New ...
func New(options ...Option) zap.Encoder {
	enc := pool.Get().(*Encoder)
	enc.buffer.Reset()
	ShortTime().Apply(enc)
	SimpleLevel().Apply(enc)
	SimpleMessage().Apply(enc)
	for _, opt := range options {
		opt.Apply(enc)
	}
	return enc
}

func (enc *Encoder) setTimeFormatter(f TimeFormatter) {
	enc.timeF = f
}

func (enc *Encoder) setLevelFormatter(f LevelFormatter) {
	enc.levelF = f
}

func (enc *Encoder) setMessageFormatter(f MessageFormatter) {
	enc.messageF = f
}

// Clone ...
func (enc *Encoder) Clone() zap.Encoder {
	clone := pool.Get().(*Encoder)
	clone.timeF = enc.timeF
	clone.levelF = enc.levelF
	clone.messageF = enc.messageF
	clone.buffer.Reset()
	enc.buffer.WriteTo(clone.buffer)
	return clone
}

// Free ...
func (enc *Encoder) Free() {
	pool.Put(enc)
}

// WriteEntry ...
func (enc *Encoder) WriteEntry(w io.Writer, message string, level zap.Level, time time.Time) error {
	buffer := bytes.NewBuffer([]byte{})
	enc.timeF(buffer, time)
	enc.levelF(buffer, level)
	enc.messageF(buffer, message)
	enc.buffer.WriteTo(buffer)
	buffer.Write([]byte{'\n'})
	buffer.WriteTo(w)
	return nil
}

func (enc *Encoder) writeKey(key string) {
	if enc.buffer.Len() > 0 {
		enc.buffer.WriteString(" ")
	}
	enc.buffer.WriteString(key)
	enc.buffer.WriteString("=")
}

// AddString ...
func (enc *Encoder) AddString(key, value string) {
	enc.writeKey(key)
	enc.buffer.WriteString(value)
}

// AddBool ...
func (enc *Encoder) AddBool(key string, value bool) {
	enc.writeKey(key)
	enc.buffer.WriteString(strconv.FormatBool(value))
}

// AddInt ...
func (enc *Encoder) AddInt(key string, value int) {
	enc.AddInt64(key, int64(value))
}

// AddInt64 ...
func (enc *Encoder) AddInt64(key string, value int64) {
	enc.writeKey(key)
	enc.buffer.WriteString(strconv.FormatInt(value, 10))
}

// AddUint ...
func (enc *Encoder) AddUint(key string, value uint) {
	enc.AddUint64(key, uint64(value))
}

// AddUint64 ...
func (enc *Encoder) AddUint64(key string, value uint64) {
	enc.writeKey(key)
	enc.buffer.WriteString(strconv.FormatUint(value, 10))
}

// AddFloat64 ...
func (enc *Encoder) AddFloat64(key string, value float64) {
	enc.writeKey(key)
	switch {
	case math.IsNaN(value):
		enc.buffer.WriteString("NaN")
	case math.IsInf(value, 1):
		enc.buffer.WriteString("+Inf")
	case math.IsInf(value, -1):
		enc.buffer.WriteString("-Inf")
	default:
		enc.buffer.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
	}
}

// AddMarshaler ...
func (enc *Encoder) AddMarshaler(key string, marshaler zap.LogMarshaler) error {
	panic("not implemented")
}

// AddObject ...
func (enc *Encoder) AddObject(key string, value interface{}) error {
	panic("not implemented")
}
