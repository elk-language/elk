package vm

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Std::Elk::Token
func initToken() {
	// Singleton methods
	c := &value.ElkTokenClass.SingletonClass().MethodContainer
	Def(
		c,
		"type_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			typ := args[1].AsUInt16()
			result := value.String(token.Type(typ).TypeName())
			return value.Ref(result), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.ElkTokenClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			typ := args[1].MustUInt16()
			span := (*value.Span)(args[2].Pointer())
			val := args[3].SafeAsReference().(value.String)

			self := token.NewWithValue(
				(*position.Span)(span),
				token.Type(typ),
				string(val),
			)
			return value.Ref(self), value.Undefined
		},
		DefWithParameters(3),
	)
	Def(
		c,
		"==",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			other := (*token.Token)(args[1].Pointer())
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"typ",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			return value.UInt16(self.Type).ToValue(), value.Undefined
		},
	)
	Def(
		c,
		"type_name",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			result := value.String(self.Type.TypeName())
			return value.Ref(result), value.Undefined
		},
	)
	Def(
		c,
		"value",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			return value.Ref(value.String(self.Value)), value.Undefined
		},
	)
	Def(
		c,
		"span",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			return value.Ref((*value.Span)(self.Span())), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*token.Token)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

	typesToNames := value.NewHashMap(token.Length())
	namesToTypes := value.NewHashMap(token.Length())
	for tokenId, tokenName := range token.Types() {
		idVal := value.UInt16(tokenId).ToValue()
		nameVal := value.Ref(value.String(tokenName))

		value.ElkTokenClass.AddConstantString(tokenName, idVal)
		HashMapSet(nil, typesToNames, idVal, nameVal)
		HashMapSet(nil, namesToTypes, nameVal, idVal)
	}

	value.ElkTokenClass.AddConstantString("TYPES_TO_NAMES", value.Ref(typesToNames))
	value.ElkTokenClass.AddConstantString("NAMES_TO_TYPES", value.Ref(namesToTypes))
}
