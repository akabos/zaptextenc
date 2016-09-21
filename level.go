package zaptextenc

import (
	"bytes"
	"fmt"

	"github.com/mgutz/ansi"
	"github.com/uber-go/zap"
)

// LevelMap defines a mapping between log level and it's representation
type LevelMap map[zap.Level]string

// LevelMapFormatterOption is an option for level formatting according to a mapping
type LevelMapFormatterOption struct {
	levelMap LevelMap
}

// Apply sets level formatter for an encoder
func (f LevelMapFormatterOption) Apply(e *Encoder) {
	e.setLevelFormatter(f.formatter())
}

func (f LevelMapFormatterOption) formatter() LevelFormatter {
	return func(w *bytes.Buffer, level zap.Level) {
		var (
			str   string
			found bool
		)
		str, found = f.levelMap[level]
		if !found {
			panic(fmt.Sprintf("unknown log level: %v", level))
		}
		w.WriteString(str)
	}
}

var simpleLevelMap = LevelMap{
	zap.DebugLevel: "DEBUG ",
	zap.InfoLevel:  "INFO  ",
	zap.WarnLevel:  "WARN  ",
	zap.ErrorLevel: "ERROR ",
	zap.PanicLevel: "PANIC ",
	zap.FatalLevel: "FATAL ",
}

// SimpleLevel formats log level in most dull and unfancy way possible
func SimpleLevel() Option {
	return &LevelMapFormatterOption{simpleLevelMap}
}

var abbrLevelMap = map[zap.Level]string{
	zap.DebugLevel: "DEB ",
	zap.InfoLevel:  "INF ",
	zap.WarnLevel:  "WRN ",
	zap.ErrorLevel: "ERR ",
	zap.PanicLevel: "PAN ",
	zap.FatalLevel: "FAT ",
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
	zap.DebugLevel: debugColor("DEBUG "),
	zap.InfoLevel:  infoColor("INFO  "),
	zap.WarnLevel:  warnColor("WARN  "),
	zap.ErrorLevel: errorColor("ERROR "),
	zap.PanicLevel: panicColor("PANIC "),
	zap.FatalLevel: fatalColor("FATAL "),
}

// SimpleColorLevel formats log level as colored level names
func SimpleColorLevel() Option {
	return &LevelMapFormatterOption{simpleColorLevelMap}
}

// NoLevelFormatterOption defines an option which makes encoder to skip log level
type NoLevelFormatterOption struct{}

// Apply sets level formatter for an encoder
func (f NoLevelFormatterOption) Apply(e *Encoder) {
	e.setLevelFormatter(func(_ *bytes.Buffer, _ zap.Level) {})
}

// NoLevel skips log level
func NoLevel() Option {
	return NoLevelFormatterOption{}
}
