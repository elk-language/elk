// Package bitset implements useful bitset
// structs which can be easily embedded
// in other structs to compress multiple bool flags.
package bitset

// 8-bit bit-flag
type BitFlag8 uint8

// An 8-bit bitset.
// Zero value is ready to use.
type Bitset8 struct {
	bitset BitFlag8
}

// Create a new 8-bit bitset from an int.
func Bitset8FromInt[T uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | uint | int](i T) Bitset8 {
	return Bitset8{
		bitset: BitFlag8(i),
	}
}

// Check whether the given bit-flag is on.
func (b *Bitset8) HasFlag(flag BitFlag8) bool {
	return b.bitset&flag != 0
}

// Turn on the given bit-flag.
func (b *Bitset8) SetFlag(flag BitFlag8) {
	b.bitset = b.bitset | flag
}
