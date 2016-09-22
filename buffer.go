package zaptextenc

// Buffer is a byte buffer
type Buffer struct {
	bytes []byte
}

// NewBuffer is Buffer constructor
func NewBuffer(size uint) *Buffer {
	return &Buffer{bytes: make([]byte, size)}
}

// Append ...
func (b *Buffer) Append(a ...byte) {
	b.bytes = append(b.bytes, a...)
}

// AppendString ...
func (b *Buffer) AppendString(a string) {
	b.bytes = append(b.bytes, a...)
}

// Bytes ...
func (b *Buffer) Bytes() []byte {
	return b.bytes
}

// Set ...
func (b *Buffer) Set(a []byte) {
	b.bytes = a
}

// Len ...
func (b *Buffer) Len() int {
	return len(b.bytes)
}

// Reset clears buffer contents without re-allocation
func (b *Buffer) Reset() {
	b.bytes = b.bytes[:0]
}
