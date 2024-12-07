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
			return self.Compare(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.LessThanEqual(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.LessThan(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.GreaterThan(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.GreaterThanEqual(other)
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			other := args[1]
			return self.Equal(other), value.Nil
		},
		DefWithParameters(1),
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
			return value.Ref(value.String(self.Uppercase())), value.Nil
		},
	)
	Def(
		c,
		"lowercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.Ref(value.String(self.Lowercase())), value.Nil
		},
	)

	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.SmallInt(self.CharCount()).ToValue(), value.Nil
		},
	)
	Alias(c, "char_count", "length")
	Def(
		c,
		"byte_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.SmallInt(self.ByteCount()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"grapheme_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.SmallInt(self.GraphemeCount()).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.ToSymbol(self).ToValue(), value.Nil
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Nil
		},
	)
	Def(
		c,
		"to_int",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return self.ToInt()
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.Ref(value.String(self.Inspect())), value.Nil
		},
	)
	Def(
		c,
		"is_empty",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			return value.ToElkBool(self.IsEmpty()), value.Nil
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			iterator := value.NewStringCharIterator(self)
			return value.Ref(iterator), value.Nil
		},
	)
	Alias(c, "char_iter", "iter")
	Def(
		c,
		"byte_iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(value.String)
			iterator := value.NewStringByteIterator(self)
			return value.Ref(iterator), value.Nil
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
			self := args[0].MustReference().(*value.StringCharIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Nil
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
			self := args[0].MustReference().(*value.StringByteIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iter",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Nil
		},
	)

}
