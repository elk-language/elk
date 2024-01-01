// Package bitfield implements useful bitfield
// structs which can be easily embedded
// in other structs to compress multiple bool flags.
package bitfield

// 8-bit bit-flag
type BitFlag8 uint8

// An 8-bit bitfield.
// Zero value is ready to use.
type Bitfield8 struct {
	bitfield BitFlag8
}

// Create a new 8-bit bitfield from an int.
func Bitfield8FromInt[T uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | uint | int](i T) Bitfield8 {
	return Bitfield8{
		bitfield: BitFlag8(i),
	}
}

// Check whether the given bit-flag is on.
func (b *Bitfield8) HasFlag(flag BitFlag8) bool {
	return b.bitfield&flag != 0
}

// Turn on the given bit-flag.
func (b *Bitfield8) SetFlag(flag BitFlag8) {
	b.bitfield = b.bitfield | flag
}

func (b *Bitfield8) Byte() byte {
	return byte(b.bitfield)
}
