package bitarray

import (
	"errors"
)

const byteSize = 8

// BitArray is a []Byte backed interface for setting individual bits in the result.
type BitArray struct {
	a []byte
}

// ErrInvalidSize is returned if the supplied size is invalid for the NewBitArray.
// this can be considered a compile time error (if you hard code the size, you can
// ignore the error return safely).
var ErrInvalidSize = errors.New("size must be positive integer and a multiple of 8")

// ErrInvalidIndex is returned when an invalid index is supplied for SetBit.
var ErrInvalidIndex = errors.New("index must be less than array size")

// NewBitArray returns a new BitArray with all values set to 0.
func NewBitArray(size int) (*BitArray, error) {
	if size < 0 || size%byteSize != 0 {
		return nil, ErrInvalidSize
	}

	return &BitArray{
		a: make([]byte, size/byteSize),
	}, nil
}

// SetBit sets a specific bit to on or off.
func (b *BitArray) SetBit(index int, on bool) error {
	if index < 0 || index > len(b.a)*byteSize {
		return ErrInvalidIndex
	}

	offset := index % byteSize
	pos := index / byteSize

	if !on {
		b.a[pos] &= ^(1 << offset)
	} else {
		b.a[pos] |= (1 << offset)
	}

	return nil
}

// ByteSlice return the []byte slice underneath the BitArray.
func (b *BitArray) ByteSlice() []byte {
	return b.a
}
