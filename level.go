package zaptextenc

import (
	"github.com/mgutz/ansi"
	"github.com/uber-go/zap"
)

// LevelMap defines a mapping between log level and it's representation
type LevelMap map[zap.Level][]byte

// LevelMapFormatterOption is an option for level formatting according to a mapping
type LevelMapFormatterOption struct {
	levelMap LevelMap
}

// Apply sets level formatter for an encoder
func (f LevelMapFormatterOption) Apply(e *Encoder) {
	e.setLevelFormatter(f.formatter())
}

func (f LevelMapFormatterOption) formatter() LevelFormatter {
	return func(b []byte, level zap.Level) {
		b = append(b, f.levelMap[level]...)
	}
}

var simpleLevelMap = LevelMap{
	zap.DebugLevel: []byte("DEBUG "),
	zap.InfoLevel:  []byte("INFO  "),
	zap.WarnLevel:  []byte("WARN  "),
	zap.ErrorLevel: []byte("ERROR "),
	zap.PanicLevel: []byte("PANIC "),
	zap.FatalLevel: []byte("FATAL "),
}

// SimpleLevel formats log level in most dull and unfancy way possible
func SimpleLevel() Option {
	return &LevelMapFormatterOption{simpleLevelMap}
}

var abbrLevelMap = LevelMap{
	zap.DebugLevel: []byte("DEB "),
	zap.InfoLevel:  []byte("INF "),
	zap.WarnLevel:  []byte("WRN "),
	zap.ErrorLevel: []byte("ERR "),
	zap.PanicLevel: []byte("PAN "),
	zap.FatalLevel: []byte("FAT "),
}

// AbbrLevel formats log level as 3-letter abbreviation
func AbbrLevel() Option {
	return &LevelMapFormatterOption{abbrLevelMap}
}

var (
	debugColor = ansi.ColorFunc("blue")
	infoColor  = ansi.ColorFunc("white")
	warnColor  = ansi.ColorFunc("yellow")
	errorColor = ansi.ColorFunc("red")
	panicColor = ansi.ColorFunc("magenta+b")
	fatalColor = ansi.ColorFunc("magenta+b")
)

var simpleColorLevelMap = LevelMap{
	zap.DebugLevel: []byte(debugColor("DEBUG ")),
	zap.InfoLevel:  []byte(infoColor("INFO  ")),
	zap.WarnLevel:  []byte(warnColor("WARN  ")),
	zap.ErrorLevel: []byte(errorColor("ERROR ")),
	zap.PanicLevel: []byte(panicColor("PANIC ")),
	zap.FatalLevel: []byte(fatalColor("FATAL ")),
}

// SimpleColorLevel formats log level as colored level names
func SimpleColorLevel() Option {
	return &LevelMapFormatterOption{simpleColorLevelMap}
}

// NoLevelFormatterOption defines an option which makes encoder to skip log level
type NoLevelFormatterOption struct{}

// Apply sets level formatter for an encoder
func (f NoLevelFormatterOption) Apply(e *Encoder) {
	e.setLevelFormatter(func(_ []byte, _ zap.Level) {})
}

// NoLevel skips log level
func NoLevel() Option {
	return NoLevelFormatterOption{}
}
