package value

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

func initSize() {
	SizeClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	RegisterNativeClass("Std::Size", "value.SizeClass")
}
