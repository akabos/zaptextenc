package zaptextenc

import (
	"io"
	"math"
	"strconv"
	"time"

	"github.com/uber-go/zap"
)

// Option is an Encoder option
type Option interface {
	Apply(*Encoder)
}

// TimeFormatter is a func used to render log entry time
type TimeFormatter func(io.Writer, time.Time)

// LevelFormatter is a func used to render log entry level
type LevelFormatter func(io.Writer, zap.Level)

// MessageFormatter is a func used to render log entry message
type MessageFormatter func(io.Writer, string)

// Encoder is zap.Encoder implementation that writes plain text messages
type Encoder struct {
	timeF    TimeFormatter
	levelF   LevelFormatter
	messageF MessageFormatter
	bytes    []byte
}

// New ...
func New(options ...Option) zap.Encoder {
	enc := &Encoder{}
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
	return &Encoder{
		timeF:    enc.timeF,
		levelF:   enc.levelF,
		messageF: enc.messageF,
	}
}

// Free ...
func (enc *Encoder) Free() {
	// noop
}

// WriteEntry ...
func (enc *Encoder) WriteEntry(sink io.Writer, message string, level zap.Level, time time.Time) error {
	enc.timeF(sink, time)
	enc.levelF(sink, level)
	enc.messageF(sink, message)
	sink.Write(enc.bytes)
	sink.Write([]byte{'\n'})
	return nil
}

func (enc *Encoder) appendKey(key string) {
	if len(enc.bytes) > 0 {
		enc.bytes = append(enc.bytes, ' ')
	}
	enc.bytes = append(enc.bytes, key...)
	enc.bytes = append(enc.bytes, '=')
}

// AddString ...
func (enc *Encoder) AddString(key, value string) {
	enc.appendKey(key)
	enc.bytes = append(enc.bytes, value...)
}

// AddBool ...
func (enc *Encoder) AddBool(key string, value bool) {
	enc.appendKey(key)
	switch value {
	case true:
		enc.bytes = append(enc.bytes, "true"...)
	case false:
		enc.bytes = append(enc.bytes, "false"...)
	default:
		panic("huh?")
	}
}

// AddInt ...
func (enc *Encoder) AddInt(key string, value int) {
	enc.AddInt64(key, int64(value))
}

// AddInt64 ...
func (enc *Encoder) AddInt64(key string, value int64) {
	enc.appendKey(key)
	enc.bytes = strconv.AppendInt(enc.bytes, value, 10)
}

// AddUint ...
func (enc *Encoder) AddUint(key string, value uint) {
	enc.AddUint64(key, uint64(value))
}

// AddUint64 ...
func (enc *Encoder) AddUint64(key string, value uint64) {
	enc.appendKey(key)
	enc.bytes = strconv.AppendUint(enc.bytes, value, 10)
}

// AddFloat64 ...
func (enc *Encoder) AddFloat64(key string, value float64) {
	enc.appendKey(key)
	switch {
	case math.IsNaN(value):
		enc.bytes = append(enc.bytes, `"NaN"`...)
	case math.IsInf(value, 1):
		enc.bytes = append(enc.bytes, `"+Inf"`...)
	case math.IsInf(value, -1):
		enc.bytes = append(enc.bytes, `"-Inf"`...)
	default:
		enc.bytes = strconv.AppendFloat(enc.bytes, value, 'f', -1, 64)
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
