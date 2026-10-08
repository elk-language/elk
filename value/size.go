package value

import (
	"fmt"
	"strings"
)

type Size int64

var _ ValueInterface = Size(5)

const (
	Byte     Size = 1
	Kilobyte      = 1024 * Byte
	Megabyte      = 1024 * Kilobyte
	Gigabyte      = 1024 * Megabyte
	Terabyte      = 1024 * Gigabyte
	Petabyte      = 1024 * Terabyte
)

var SizeClass *Class

func (s Size) InBytes() Float {
	return Float(s)
}

func (s Size) Bytes() Value {
	return ToElkInt(int64(s))
}

func (s Size) InKilobytes() Float {
	return Float(s) / Float(Kilobyte)
}

func (s Size) Kilobytes() Value {
	return ToElkInt(int64(s / Kilobyte))
}

func (s Size) InMegabytes() Float {
	return Float(s) / Float(Megabyte)
}

func (s Size) Megabytes() Value {
	return ToElkInt(int64(s / Megabyte))
}

func (s Size) InGigabytes() Float {
	return Float(s) / Float(Gigabyte)
}

func (s Size) Gigabytes() Value {
	return ToElkInt(int64(s / Gigabyte))
}

func (s Size) InTerabytes() Float {
	return Float(s) / Float(Terabyte)
}

func (s Size) Terabytes() Value {
	return ToElkInt(int64(s / Terabyte))
}

func (s Size) InPetabyte() Float {
	return Float(s) / Float(Petabyte)
}

func (s Size) Petabytes() Value {
	return ToElkInt(int64(s / Petabyte))
}

func initSize() {
	SizeClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	RegisterNativeClass("Std::Size", "value.SizeClass")
}

func (Size) Class() *Class {
	return SizeClass
}

func (Size) DirectClass() *Class {
	return SizeClass
}

func (Size) SingletonClass() *Class {
	return nil
}

func (Size) InstanceVariables() *InstanceVariables {
	return nil
}

func (s Size) Inspect() string {
	return fmt.Sprintf("Std::Size.parse(%s)", s.ToString().Inspect())
}

func (s Size) Error() string {
	return s.Inspect()
}

func (s Size) ToString() String {
	return String(s.String())
}

func (s Size) String() string {
	var buff strings.Builder

	if s == 0 {
		return "0.kilobytes"
	}

	petabytes := s / Petabyte
	if petabytes != 0 {
		fmt.Fprintf(&buff, "%d", petabytes)
	}

	return buff.String()
}

// const (
// 	Byte     Size = 1
// 	Kilobyte      = 1024 * Byte
// 	Megabyte      = 1024 * Kilobyte
// 	Gigabyte      = 1024 * Megabyte
// 	Terabyte      = 1024 * Gigabyte
// 	Petabyte      = 1024 * Terabyte
// )

func (s Size) Copy() Reference {
	return s
}
