// Package bitfield implements useful bitfield
// structs which can be easily embedded
// in other structs to compress multiple bool flags.
package bitfield

import "iter"

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
func BitField8FromInt[T ~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64 | ~uint | ~int](i T) BitField8 {
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

// Turn on/off the given bit-flag.
func (b *BitField8) SetFlagValue(flag BitFlag8, val bool) {
	if val {
		b.SetFlag(flag)
	} else {
		b.UnsetFlag(flag)
	}
}

// Turn off the given bit-flag.
func (b *BitField8) UnsetFlag(flag BitFlag8) {
	b.bitfield = b.bitfield &^ flag
}

func (b BitField8) Byte() byte {
	return byte(b.bitfield)
}

func (b BitField8) AllSetFlags() iter.Seq2[int, BitFlag8] {
	return func(yield func(int, BitFlag8) bool) {
		for i := range 8 {
			flag := BitFlag8(1 << i)
			if b.HasFlag(flag) {
				if !yield(i, flag) {
					return
				}
			}
		}
	}
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

// Create a new 16-bit bit field from an int.
func BitField16FromInt[T ~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64 | ~uint | ~int](i T) BitField16 {
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

// Turn on/off the given bit-flag.
func (b *BitField16) SetFlagValue(flag BitFlag16, val bool) {
	if val {
		b.SetFlag(flag)
	} else {
		b.UnsetFlag(flag)
	}
}

func (b BitField16) ToBitFlag() BitFlag16 {
	return b.bitfield
}

func (b BitField16) Uint16() uint16 {
	return uint16(b.bitfield)
}

func (b BitField16) AllSetFlags() iter.Seq2[int, BitFlag16] {
	return func(yield func(int, BitFlag16) bool) {
		for i := range 16 {
			flag := BitFlag16(1 << i)
			if b.HasFlag(flag) {
				if !yield(i, flag) {
					return
				}
			}
		}
	}
}

// 32-bit bit-flag
type BitFlag32 uint32

// An 32-bit bit field.
// Zero value is ready to use.
type BitField32 struct {
	bitfield BitFlag32
}

// Create a new 16-bit bit field from an int.
func BitField32FromBitFlag(f BitFlag32) BitField32 {
	return BitField32{
		bitfield: f,
	}
}

// Create a new 32-bit bit field from an int.
func BitField32FromInt[T ~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64 | ~uint | ~int](i T) BitField32 {
	return BitField32{
		bitfield: BitFlag32(i),
	}
}

// Check whether any bit flag is on.
func (b BitField32) IsAnyFlagSet() bool {
	return b.bitfield != 0
}

// Check whether the given bit-flag is on.
func (b BitField32) HasFlag(flag BitFlag32) bool {
	return b.bitfield&flag != 0
}

// Turn on the given bit-flag.
func (b *BitField32) SetFlag(flag BitFlag32) {
	b.bitfield = b.bitfield | flag
}

// Turn off the given bit-flag.
func (b *BitField32) UnsetFlag(flag BitFlag32) {
	b.bitfield = b.bitfield &^ flag
}

// Turn on/off the given bit-flag.
func (b *BitField32) SetFlagValue(flag BitFlag32, val bool) {
	if val {
		b.SetFlag(flag)
	} else {
		b.UnsetFlag(flag)
	}
}

func (b BitField32) ToBitFlag() BitFlag32 {
	return b.bitfield
}

func (b BitField32) Uint32() uint32 {
	return uint32(b.bitfield)
}

func (b BitField32) AllSetFlags() iter.Seq2[int, BitFlag32] {
	return func(yield func(int, BitFlag32) bool) {
		for i := range 32 {
			flag := BitFlag32(1 << i)
			if b.HasFlag(flag) {
				if !yield(i, flag) {
					return
				}
			}
		}
	}
}
