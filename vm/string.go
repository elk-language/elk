package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	// Instance methods
	c := &value.StringClass.MethodContainer
	Def(
		c,
		"+",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Concat(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Alias(c, "concat", "+")

	Def(
		c,
		"-",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.RemoveSuffix(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Alias(c, "remove_suffix", "-")

	Def(
		c,
		"*",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Repeat(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Alias(c, "repeat", "*")

	Def(
		c,
		"<=>",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Compare(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.LessThanEqual(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"<",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.LessThan(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.GreaterThan(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		">=",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.GreaterThanEqual(other))
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return self.Equal(other), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"===",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return self.StrictEqual(other), nil
		},
		DefWithParameters("other"),
		DefWithSealed(),
	)
	Def(
		c,
		"char_at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.Subscript(other))
		},
		DefWithParameters("index"),
		DefWithSealed(),
	)
	Alias(c, "[]", "char_at")
	Def(
		c,
		"byte_at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.ByteAt(other))
		},
		DefWithParameters("index"),
		DefWithSealed(),
	)
	Def(
		c,
		"grapheme_at",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			other := args[1]
			return value.ToValueErr(self.GraphemeAt(other))
		},
		DefWithParameters("index"),
		DefWithSealed(),
	)

	Def(
		c,
		"uppercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return self.Uppercase(), nil
		},
	)
	Def(
		c,
		"lowercase",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return self.Lowercase(), nil
		},
	)

	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.CharCount()), nil
		},
	)
	Alias(c, "char_count", "length")

	Def(
		c,
		"byte_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.ByteCount()), nil
		},
	)
	Def(
		c,
		"grapheme_count",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.SmallInt(self.GraphemeCount()), nil
		},
	)
	Def(
		c,
		"to_symbol",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.ToSymbol(self), nil
		},
	)
	Def(
		c,
		"to_string",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.String(self.Inspect()), nil
		},
	)
	Def(
		c,
		"is_empty",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			return value.ToElkBool(self.IsEmpty()), nil
		},
	)
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			iterator := value.NewStringCharIterator(self)
			return iterator, nil
		},
	)
	Alias(c, "char_iterator", "iterator")
	Def(
		c,
		"byte_iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.String)
			iterator := value.NewStringByteIterator(self)
			return iterator, nil
		},
	)

}

// ::Std::String::CharIterator
func init() {
	// Instance methods
	c := &value.StringCharIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.StringCharIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)

}

// ::Std::String::ByteIterator
func init() {
	// Instance methods
	c := &value.StringByteIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.StringByteIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)

}
