# :zap: zap-text

Text encoder for blazing fast, structured, leveled logging in Go.

## Installation

`go get -u github.com/akabos/zaptextenc`

## Usage

```
import (
    github.com/akabos/zaptextenc
    github.com/uber-go/zap
)

logger := zap.New(zaptextenc.New())

logger.Info("Failed to fetch URL.",
  zap.String("url", url),
  zap.Int("attempt", tryNum),
  zap.Duration("backoff", sleepFor),
)
```
