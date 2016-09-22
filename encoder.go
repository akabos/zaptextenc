package zaptextenc

import (
	"io"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/uber-go/zap"
)

var pool = sync.Pool{New: func() interface{} {
	return &Encoder{
		// Pre-allocate some capacity to avoid allocations
		fields: make([]byte, 512),
		prefix: make([]byte, 512),
		output: make([]byte, 1024),
	}
}}

// Option is an Encoder option
type Option interface {
	Apply(*Encoder)
}

// TimeFormatter is a func used to render log entry time
type TimeFormatter func([]byte, time.Time)

// LevelFormatter is a func used to render log entry level
type LevelFormatter func([]byte, zap.Level)

// MessageFormatter is a func used to render log entry message
type MessageFormatter func([]byte, string)

// Encoder is zap.Encoder implementation that writes plain text messages
type Encoder struct {
	timeF    TimeFormatter
	levelF   LevelFormatter
	messageF MessageFormatter
	fields   []byte
	prefix   []byte
	output   []byte
}

// New ...
func New(options ...Option) zap.Encoder {
	enc := pool.Get().(*Encoder)
	enc.truncate()
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

func (enc *Encoder) truncate() {
	enc.fields = enc.fields[:0]
	enc.prefix = enc.prefix[:0]
	enc.output = enc.output[:0]
}

// Clone ...
func (enc *Encoder) Clone() zap.Encoder {
	clone := pool.Get().(*Encoder)
	clone.truncate()
	clone.fields = append(clone.fields, enc.fields...)
	clone.prefix = append(clone.prefix, enc.prefix...)
	clone.output = append(clone.output, enc.output...)

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
	enc.output = append(enc.output, enc.prefix...)
	enc.output = append(enc.output, enc.fields...)
	enc.output = append(enc.output, '\n')
	w.Write(enc.output)
	return nil
}

func (enc *Encoder) writeKey(key string) {
	if len(enc.fields) > 0 {
		enc.fields = append(enc.fields, ' ')
	}
	enc.fields = append(enc.fields, []byte(key)...)
	enc.fields = append(enc.fields, '=')
}

// AddString ...
func (enc *Encoder) AddString(key, value string) {
	enc.writeKey(key)
	enc.fields = append(enc.fields, value...)
}

// AddBool ...
func (enc *Encoder) AddBool(key string, value bool) {
	enc.writeKey(key)
	strconv.AppendBool(enc.fields, value)
}

// AddInt ...
func (enc *Encoder) AddInt(key string, value int) {
	enc.AddInt64(key, int64(value))
}

// AddInt64 ...
func (enc *Encoder) AddInt64(key string, value int64) {
	enc.writeKey(key)
	strconv.AppendInt(enc.fields, value, 10)
}

// AddUint ...
func (enc *Encoder) AddUint(key string, value uint) {
	enc.AddUint64(key, uint64(value))
}

// AddUint64 ...
func (enc *Encoder) AddUint64(key string, value uint64) {
	enc.writeKey(key)
	strconv.AppendUint(enc.fields, value, 10)
}

// AddFloat64 ...
func (enc *Encoder) AddFloat64(key string, value float64) {
	enc.writeKey(key)
	switch {
	case math.IsNaN(value):
		enc.fields = append(enc.fields, []byte("NaN")...)
	case math.IsInf(value, 1):
		enc.fields = append(enc.fields, []byte("+Inf")...)
	case math.IsInf(value, -1):
		enc.fields = append(enc.fields, []byte("-Inf")...)
	default:
		strconv.AppendFloat(enc.fields, value, 'f', -1, 64)
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
