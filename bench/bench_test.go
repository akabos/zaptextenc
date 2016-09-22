package bench

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/akabos/zaptextenc"
	"github.com/uber-go/zap"
)

var errExample = errors.New("fail")

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func runSimpleWithEncoder(b *testing.B, encoder *zap.Encoder) {
	logger := zap.New(*encoder, zap.DiscardOutput)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.")
		}
	})
}

func run10FieldsWithEncoder(b *testing.B, encoder *zap.Encoder) {
	logger := zap.New(*encoder, zap.DiscardOutput)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.",
				zap.Int("int", 1),
				zap.Int64("int", int64(1)),
				zap.Float64("float64", float64(3.0)),
				zap.String("string", "four!"),
				zap.Bool("true", true),
				zap.Bool("false", false),
				zap.Time("time", time.Unix(0, 0)),
				zap.Duration("duration", time.Second),
				zap.Error(errExample),
				zap.String("another string", "done!"),
			)
		}
	})
}

func BenchmarkJSONSimple(b *testing.B) {
	encoder := zap.NewJSONEncoder(zap.NoTime())
	runSimpleWithEncoder(b, &encoder)
}

func BenchmarkJSON10Fields(b *testing.B) {
	encoder := zap.NewJSONEncoder(zap.NoTime())
	run10FieldsWithEncoder(b, &encoder)
}

func BenchmarkTextSimple(b *testing.B) {
	encoder := zaptextenc.New(zaptextenc.NoTime())
	runSimpleWithEncoder(b, &encoder)
}

func BenchmarkText10Fields(b *testing.B) {
	encoder := zaptextenc.New(zaptextenc.NoTime())
	run10FieldsWithEncoder(b, &encoder)
}
