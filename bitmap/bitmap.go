// Package bitmap implements useful bitmap
// structs which can be easily embedded
// in other structs to compress multiple bool flags.
package bitmap

// 8-bit bit-flag
type BitFlag8 uint8

// An 8-bit bitmap.
// Zero value is ready to use.
type Bitmap8 struct {
	bitmap BitFlag8
}

// Create a new 8-bit bitmap from an int.
func Bitmap8FromInt[T uint8 | int8 | uint16 | int16 | uint32 | int32 | uint64 | int64 | uint | int](i T) Bitmap8 {
	return Bitmap8{
		bitmap: BitFlag8(i),
	}
}

// Check whether the given bit-flag is on.
func (b *Bitmap8) HasFlag(flag BitFlag8) bool {
	return b.bitmap&flag != 0
}

// Turn on the given bit-flag.
func (b *Bitmap8) SetFlag(flag BitFlag8) {
	b.bitmap = b.bitmap | flag
}
