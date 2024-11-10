// Package bitfield implements useful bitfield
// structs which can be easily embedded
// in other structs to compress multiple bool flags.
package bitfield

// 8-bit bit-flag
type BitFlag8 uint8

// An 8-bit bit field.
// Zero value is ready to use.
type BitField8 struct {
	bitfield BitFlag8
}

// Create a new 8-bit bit field from an int.
func BitField8FromBitFlag(f BitFlag8) BitField8 {
	return BitField8{
		bitfield: f,
	}
}

// Create a new 8-bit bit field from an int.
func BitField8FromInt[T uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | uint | int](i T) BitField8 {
	return BitField8{
		bitfield: BitFlag8(i),
	}
}

// Check whether any bit flag is on.
func (b BitField8) IsAnyFlagSet() bool {
	return b.bitfield != 0
}

// Check whether the given bit-flag is on.
func (b BitField8) HasFlag(flag BitFlag8) bool {
	return b.bitfield&flag != 0
}

// Turn on the given bit-flag.
func (b *BitField8) SetFlag(flag BitFlag8) {
	b.bitfield = b.bitfield | flag
}

// Turn off the given bit-flag.
func (b *BitField8) UnsetFlag(flag BitFlag8) {
	b.bitfield = b.bitfield &^ flag
}

func (b BitField8) Byte() byte {
	return byte(b.bitfield)
}

// 16-bit bit-flag
type BitFlag16 uint16

// An 16-bit bit field.
// Zero value is ready to use.
type BitField16 struct {
	bitfield BitFlag16
}

// Create a new 16-bit bit field from an int.
func BitField16FromBitFlag(f BitFlag16) BitField16 {
	return BitField16{
		bitfield: f,
	}
}

// Create a new 8-bit bit field from an int.
func BitField16FromInt[T uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | uint | int](i T) BitField16 {
	return BitField16{
		bitfield: BitFlag16(i),
	}
}

// Check whether any bit flag is on.
func (b BitField16) IsAnyFlagSet() bool {
	return b.bitfield != 0
}

// Check whether the given bit-flag is on.
func (b BitField16) HasFlag(flag BitFlag16) bool {
	return b.bitfield&flag != 0
}

// Turn on the given bit-flag.
func (b *BitField16) SetFlag(flag BitFlag16) {
	b.bitfield = b.bitfield | flag
}

// Turn off the given bit-flag.
func (b *BitField16) UnsetFlag(flag BitFlag16) {
	b.bitfield = b.bitfield &^ flag
}

func (b BitField16) Byte() byte {
	return byte(b.bitfield)
}
