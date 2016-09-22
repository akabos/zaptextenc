package zaptextenc

import (
	"io"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/uber-go/zap"
)

// Option is an Encoder option
type Option interface {
	Apply(*Encoder)
}

// TimeFormatter is a func used to render log entry time
type TimeFormatter func(*Buffer, time.Time)

// LevelFormatter is a func used to render log entry level
type LevelFormatter func(*Buffer, zap.Level)

// MessageFormatter is a func used to render log entry message
type MessageFormatter func(*Buffer, string)

// Encoder is zap.Encoder implementation that writes plain text messages
type Encoder struct {
	timeF    TimeFormatter
	levelF   LevelFormatter
	messageF MessageFormatter
	fields   *Buffer
	prefix   *Buffer
	output   *Buffer
}

var pool = sync.Pool{New: func() interface{} {
	return &Encoder{
		fields: NewBuffer(512),
		prefix: NewBuffer(512),
		output: NewBuffer(1024),
	}
}}

// New ...
func New(options ...Option) zap.Encoder {
	enc := pool.Get().(*Encoder)
	enc.reset()
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

func (enc *Encoder) reset() {
	enc.fields.Reset()
	enc.prefix.Reset()
	enc.output.Reset()
}

// Clone ...
func (enc *Encoder) Clone() zap.Encoder {
	clone := pool.Get().(*Encoder)

	clone.reset()
	clone.fields.Append(enc.fields.Bytes()...)
	clone.prefix.Append(enc.prefix.Bytes()...)
	clone.output.Append(enc.output.Bytes()...)

	clone.timeF = enc.timeF
	clone.levelF = enc.levelF
	clone.messageF = enc.messageF

	return clone
}

// Free ...
func (enc *Encoder) Free() {
	pool.Put(enc)
}

// WriteEntry ...
func (enc *Encoder) WriteEntry(w io.Writer, message string, level zap.Level, time time.Time) error {
	enc.timeF(enc.prefix, time)
	enc.levelF(enc.prefix, level)
	enc.messageF(enc.prefix, message)
	enc.output.Append(enc.prefix.Bytes()...)
	enc.output.Append(enc.fields.Bytes()...)
	enc.output.Append('\n')
	_, err := w.Write(enc.output.Bytes())
	return err
}

func (enc *Encoder) writeKey(key string) {
	if enc.fields.Len() > 0 {
		enc.fields.Append(' ')
	}
	enc.fields.AppendString(key)
	enc.fields.Append('=')
}

// AddString ...
func (enc *Encoder) AddString(key, value string) {
	enc.writeKey(key)
	enc.fields.AppendString(value)
}

// AddBool ...
func (enc *Encoder) AddBool(key string, value bool) {
	enc.writeKey(key)
	strconv.AppendBool(enc.fields.bytes, value)
}

// AddInt ...
func (enc *Encoder) AddInt(key string, value int) {
	enc.AddInt64(key, int64(value))
}

// AddInt64 ...
func (enc *Encoder) AddInt64(key string, value int64) {
	enc.writeKey(key)
	enc.fields.Set(strconv.AppendInt(enc.fields.Bytes(), value, 10))
}

// AddUint ...
func (enc *Encoder) AddUint(key string, value uint) {
	enc.AddUint64(key, uint64(value))
}

// AddUint64 ...
func (enc *Encoder) AddUint64(key string, value uint64) {
	enc.writeKey(key)
	enc.fields.Set(strconv.AppendUint(enc.fields.Bytes(), value, 10))
}

// AddFloat64 ...
func (enc *Encoder) AddFloat64(key string, value float64) {
	enc.writeKey(key)
	switch {
	case math.IsNaN(value):
		enc.fields.AppendString("NaN")
	case math.IsInf(value, 1):
		enc.fields.AppendString("+Inf")
	case math.IsInf(value, -1):
		enc.fields.AppendString("-Inf")
	default:
		enc.fields.Set(strconv.AppendFloat(enc.fields.bytes, value, 'f', -1, 64))
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
