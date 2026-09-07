package value

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/token"
)

var ElkTokenClass *Class // Std::Elk::Token

// Elk runtime view of a lexer token.
type Token token.Token

func initElkToken() {
	ElkTokenClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ElkModule.AddConstantString("Token", Ref(ElkTokenClass))
	RegisterNativeClass("Std::Elk::Token", "value.ElkTokenClass")
}

func (*Token) Class() *Class {
	return ElkTokenClass
}

func (*Token) DirectClass() *Class {
	return ElkTokenClass
}

func (*Token) SingletonClass() *Class {
	return nil
}

func (t *Token) Copy() Reference {
	return t
}

func (t *Token) ToValue() Value {
	return Ref(t)
}

func (t *Token) InstanceVariables() *InstanceVariables {
	return nil
}

func (t *Token) Inspect() string {
	tok := (*token.Token)(t)
	var buff strings.Builder

	buff.WriteString("Std::Token{")
	if tok.Value != "" {
		fmt.Fprintf(&buff, "value: %s, ", String(tok.Value).Inspect())
	}

	fmt.Fprintf(
		&buff,
		"typ: %s, span: %s}",
		tok.Type.TypeName(),
		(*Span)(tok.Span()).Inspect(),
	)

	return buff.String()
}

func (t *Token) Error() string {
	return t.Inspect()
}
