package value

import (
	"fmt"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/token"
)

var ElkLexerClass *Class // ::Std::Elk::Lexer

// Elk runtime view of a lexer.
type Lexer lexer.Lexer

func initElkLexer() {
	ElkLexerClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ElkModule.AddConstantString("Lexer", Ref(ElkLexerClass))
	RegisterNativeClass("Std::Elk::Lexer", "value.ElkLexerClass")
}

func (*Lexer) Class() *Class {
	return ElkLexerClass
}

func (*Lexer) DirectClass() *Class {
	return ElkLexerClass
}

func (*Lexer) SingletonClass() *Class {
	return nil
}

func (l *Lexer) Copy() Reference {
	return l
}

func (l *Lexer) ToValue() Value {
	return Ref(l)
}

func (*Lexer) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *Lexer) Inspect() string {
	lex := (*lexer.Lexer)(l)
	return fmt.Sprintf("Std::Elk::Lexer{&: %p, source_name: %s}", lex, String(lex.SourceName()).Inspect())
}

func (l *Lexer) Error() string {
	return l.Inspect()
}

func (l *Lexer) Next() *Token {
	return (*Token)((*lexer.Lexer)(l).Next())
}

// Lex the given string and return a list of Elk token values.
func LexValue(source string) *ArrayListOfValue {
	l := lexer.New(source)

	tokens := NewArrayListOfValue(10)
	for {
		tok := l.Next()
		if tok.Type == token.END_OF_FILE {
			break
		}
		tokens.Append(Ref((*Token)(tok)))
	}
	return tokens
}
