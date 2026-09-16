package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/lexer"
)

// Value used to decide what whether to skip the subtypes of a type,
// break the traversal or continue in the type Traverse method.
// The zero value continues the traversal.
type TraverseOption uint8

const (
	TraverseContinue TraverseOption = iota
	TraverseSkip
	TraverseBreak
)

// Type ID that indicates
// the absence of an ID
const ZERO_ID = 0

// Unique ID of the type
type ID uint32

// Reference to a type using an ID
type Ref[T Type] ID

func FetchID(t Type) ID {
	id := t.ID()
	if id != ZERO_ID {
		return id
	}

	return Env.RegisterType(t).ID()
}

func ToRef[T Type](t T) Ref[T] {
	return Ref[T](FetchID(t))
}

func CastRef[T Type](t Type) Ref[T] {
	return Ref[T](FetchID(t))
}

func (ref Ref[T]) ID() ID {
	return ID(ref)
}

func (ref Ref[T]) Cast[R Type]() Ref[R] {
	return Ref[R](ref)
}

func (ref Ref[T]) IsZero() bool {
	return ref == ZERO_ID
}

func (ref Ref[T]) IsPresent() bool {
	return ref != ZERO_ID
}

func (ref Ref[T]) Get() (result T) {
	if ref == 0 {
		return
	}

	return Env.GetType(ref.ID()).(T)
}

type Type interface {
	ds.Hashable
	ID() ID
	SetID(ID)
	ToNonLiteral() Type
	IsLiteral() bool
	inspect() string
	traverse(parent Type, enter func(typ, parent Type) TraverseOption, leave func(typ, parent Type) TraverseOption) TraverseOption
}

func noopTraverse(typ, parent Type) TraverseOption { return TraverseContinue }

func Traverse(typ Type, enter func(typ, parent Type) TraverseOption, leave func(typ, parent Type) TraverseOption) {
	if enter == nil {
		enter = noopTraverse
	}
	if leave == nil {
		leave = noopTraverse
	}
	typ.traverse(nil, enter, leave)
}

func IsPointerNil(val any) bool {
	if val == nil {
		return true
	}

	value := reflect.ValueOf(val)
	kind := value.Kind()
	return kind == reflect.Pointer && value.IsNil()
}

type ModifierSet struct {
	Abstract  bool
	Sealed    bool
	Primitive bool
	NoInit    bool
	Pure      bool
	Immutable bool
}

func InspectModifier(modifiers ModifierSet) string {
	str := ""
	if modifiers.Pure {
		str += " pure"
	}
	if modifiers.Abstract {
		str += " abstract"
	}
	if modifiers.Sealed {
		str += " sealed"
	}
	if modifiers.Primitive {
		str += " primitive"
	}
	if modifiers.NoInit {
		str += " noinit"
	}
	if modifiers.Immutable {
		str += " immutable"
	}
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return "default"
	}
	return str
}

func Inspect(typ Type) string {
	if typ == nil {
		return "void"
	}

	return typ.inspect()
}

func InspectInstanceVariable(name string) string {
	return fmt.Sprintf("@%s", name)
}

func InspectInstanceVariableWithColor(name string) string {
	return lexer.Colorize(InspectInstanceVariable(name))
}

func InspectInstanceVariableDeclaration(name string, typ Type) string {
	return fmt.Sprintf("var @%s: %s", name, Inspect(typ))
}

func InspectInstanceVariableDeclarationWithColor(name string, typ Type) string {
	return lexer.Colorize(InspectInstanceVariableDeclaration(name, typ))
}

func InspectInstanceValueDeclaration(name string, typ Type) string {
	return fmt.Sprintf("val @%s: %s", name, Inspect(typ))
}

func InspectInstanceValueDeclarationWithColor(name string, typ Type) string {
	return lexer.Colorize(InspectInstanceValueDeclaration(name, typ))
}

func InspectWithColor(typ Type) string {
	return lexer.Colorize(Inspect(typ))
}

func I(typ Type) string {
	return InspectWithColor(typ)
}

func GetMethod(typ Type, name string) *Method {
	typ = typ.ToNonLiteral()

	switch t := typ.(type) {
	case *Class:
		return t.MethodString(name)
	case *Module:
		return t.MethodString(name)
	}

	return nil
}
