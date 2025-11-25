package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::String
func initString() {
	// Instance methods
	c := &value.StringClass.MethodContainer
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return value.RefErr(self.Concat(other))
		},
		DefWithParameters(1),
	)
	Alias(c, "concat", "+")

	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return value.RefErr(self.RemoveSuffix(other))
		},
		DefWithParameters(1),
	)
	Alias(c, "remove_suffix", "-")

	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return value.RefErr(self.Repeat(other))
		},
		DefWithParameters(1),
	)
	Alias(c, "repeat", "*")
	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.CompareVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.LessThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.LessThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.GreaterThanVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.GreaterThanEqualVal(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.EqualVal(other), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"rjust",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			len := args[1].AsInt()
			padding := args[2].AsChar()
			return value.Ref(self.RJust(len, padding)), value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"ljust",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			len := args[1].AsInt()
			padding := args[2].AsChar()
			return value.Ref(self.LJust(len, padding)), value.Undefined
		},
		DefWithParameters(2),
	)
	Def(
		c,
		"char_at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return value.ToValueErr(self.Subscript(other))
		},
		DefWithParameters(1),
	)
	Alias(c, "[]", "char_at")
	Def(
		c,
		"byte_at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return value.ToValueErr(self.ByteAt(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"grapheme_at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return value.RefErr(self.GraphemeAt(other))
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"uppercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.Ref(value.String(self.Uppercase())), value.Undefined
		},
	)
	Def(
		c,
		"lowercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.Ref(value.String(self.Lowercase())), value.Undefined
		},
	)

	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.SmallInt(self.CharCount()).ToValue(), value.Undefined
		},
	)
	Alias(c, "char_count", "length")
	Def(
		c,
		"byte_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.SmallInt(self.ByteCount()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"grapheme_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.SmallInt(self.GraphemeCount()).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.ToSymbol(self).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)

			var base int
			if !args[1].IsUndefined() {
				base = args[1].AsInt()
			}

			return self.ToInt(base)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)
	Def(
		c,
		"is_empty",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.ToElkBool(self.IsEmpty()), value.Undefined
		},
	)
	Def(
		c,
		"hash",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return self.Hash().ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			iterator := value.NewStringCharIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Alias(c, "char_iter", "iter")
	Def(
		c,
		"byte_iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			iterator := value.NewStringByteIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)
	Def(
		c,
		"grapheme_iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			iterator := value.NewStringGraphemeIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)

}

// ::Std::String::CharIterator
func initStringCharIterator() {
	// Instance methods
	c := &value.StringCharIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StringCharIterator)(args[0].Pointer())
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StringCharIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

// ::Std::String::ByteIterator
func initStringByteIterator() {
	// Instance methods
	c := &value.StringByteIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StringByteIterator)(args[0].Pointer())
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StringByteIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)

}

// ::Std::String::GraphemeIterator
func initStringGraphemeIterator() {
	// Instance methods
	c := &value.StringGraphemeIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StringGraphemeIterator)(args[0].Pointer())
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	Def(
		c,
		"reset",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.StringGraphemeIterator)(args[0].Pointer())
			self.Reset()
			return args[0], value.Undefined
		},
	)

}
