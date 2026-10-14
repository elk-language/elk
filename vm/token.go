package vm

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Std::Elk::Token
func initToken() {
	// Singleton methods
	c := value.ElkTokenClass.SingletonClass()
	Def(
		c,
		"type_name",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			typ := args[1].AsUInt16()
			result := value.String(token.Type(typ).TypeName())
			return value.Ref(result), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = value.ElkTokenClass
	Def(
		c,
		"#init",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			typ := args[1].MustUInt16()

			var val string
			if args[2].IsNotUndefined() {
				val = args[2].AsReference().(value.String).String()
			}

			var loc *position.Location
			if args[3].IsUndefined() {
				loc = position.ZeroLocation
			} else {
				loc = (*position.Location)(args[3].Pointer())
			}

			self := token.NewWithValue(
				loc,
				token.Type(typ),
				val,
			)
			return value.Ref((*value.Token)(self)), value.Undefined
		},
		DefWithParameters(3),
	)
	Def(
		c,
		"==",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			other := (*token.Token)(args[1].Pointer())
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"=~",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			other := (*token.Token)(args[1].Pointer())
			return value.BoolVal(self.LaxEqual(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"typ",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			return value.UInt16(self.Type).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"type_name",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			result := value.String(self.Type.TypeName())
			return value.Ref(result), value.Undefined
		},
	)
	Def(
		c,
		"value",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			return value.Ref(value.String(self.Value)), value.Undefined
		},
	)
	Def(
		c,
		"location",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			return value.Ref((*value.Location)(self.Location())), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Token)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

	typesToNames := NewHashMapOfValue(token.Length())
	namesToTypes := NewHashMapOfValue(token.Length())
	for tokenId, tokenName := range token.Types() {
		idVal := value.UInt16(tokenId).ToValue()
		nameVal := value.Ref(value.String(tokenName))

		value.ElkTokenClass.AddConstantString(tokenName, idVal)
		HashMapOfValueSet(nil, typesToNames, idVal, nameVal)
		HashMapOfValueSet(nil, namesToTypes, nameVal, idVal)
	}

	value.ElkTokenClass.AddConstantString("TYPES_TO_NAMES", value.Ref(typesToNames))
	value.ElkTokenClass.AddConstantString("NAMES_TO_TYPES", value.Ref(namesToTypes))
}
