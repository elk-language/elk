package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/value/symbol"
)

type Variance uint8

func (v Variance) String() string {
	return varianceString[v]
}

func (v Variance) Name() string {
	switch v {
	case COVARIANT:
		return "covariant"
	case CONTRAVARIANT:
		return "contravariant"
	case BIVARIANT:
		return "bivariant"
	case INVARIANT:
		return "invariant"
	default:
		return "<unknown>"
	}
}

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
	BIVARIANT
)

var varianceString = []string{
	INVARIANT:     "INVARIANT",
	COVARIANT:     "COVARIANT",
	CONTRAVARIANT: "CONTRAVARIANT",
	BIVARIANT:     "BIVARIANT",
}

type TypeParameter struct {
	Name       symbol.Symbol
	Namespace  Ref[Namespace]
	LowerBound Ref[Type]
	UpperBound Ref[Type]
	Default    Ref[Type]
	Variance   Variance
	id         ID
}

func (n *TypeParameter) ToRef() Ref[*TypeParameter] {
	return Ref[*TypeParameter](n.id)
}

func (t *TypeParameter) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("typeparam:")
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(t.Namespace))
	d.Write(b)
	d.WriteString(t.Name.String())

	return d.Sum64()
}

func (t *TypeParameter) EqualAny(other any) bool {
	o, ok := other.(*TypeParameter)
	if !ok {
		return false
	}

	if t.id > 0 {
		return t.id == o.id
	}

	return t.Namespace == o.Namespace &&
		t.Name == o.Name &&
		t.LowerBound == o.LowerBound &&
		t.UpperBound == o.UpperBound &&
		t.Default == o.Default &&
		t.Variance == o.Variance
}

func (t *TypeParameter) ID() ID {
	return t.id
}

func (t *TypeParameter) SetID(id ID) {
	t.id = id
}

// Returns the default inferred type for the type parameter
func (t *TypeParameter) InferredType() Type {
	if !t.Default.IsZero() {
		return t.Default.Get()
	}
	if !IsNeverRef(t.LowerBound) {
		return t.LowerBound.Get()
	}

	return t.UpperBound.Get()
}

func (t *TypeParameter) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(t, parent) {
	case TraverseBreak:
		return TraverseBreak
	default:
		return leave(t, parent)
	}
}

func (t *TypeParameter) Copy() *TypeParameter {
	return &TypeParameter{
		Name:       t.Name,
		Namespace:  t.Namespace,
		LowerBound: t.LowerBound,
		UpperBound: t.UpperBound,
		Variance:   t.Variance,
		id:         t.id,
	}
}

func NewTypeParameter(name symbol.Symbol, namespace Namespace, lowerBound, upperBound, def Type, variance Variance) *TypeParameter {
	t := &TypeParameter{
		Name:       name,
		Namespace:  ToRef(namespace),
		LowerBound: ToRef(lowerBound),
		UpperBound: ToRef(upperBound),
		Default:    ToRef(def),
		Variance:   variance,
	}
	return t
}

func (t *TypeParameter) ToNonLiteral() Type {
	return t
}

func (*TypeParameter) IsLiteral() bool {
	return false
}

func (t *TypeParameter) inspect() string {
	return t.Name.String()
}

func (t *TypeParameter) InspectSignature() string {
	buffer := new(strings.Builder)
	switch t.Variance {
	case COVARIANT:
		buffer.WriteRune('+')
	case CONTRAVARIANT:
		buffer.WriteRune('-')
	case BIVARIANT:
		buffer.WriteString("+-")
	}
	buffer.WriteString(t.Name.String())

	if !IsNeverRef(t.LowerBound) {
		buffer.WriteString(" > ")
		buffer.WriteString(Inspect(t.LowerBound.Get()))
	}

	if !IsAnyRef(t.UpperBound) {
		buffer.WriteString(" < ")
		buffer.WriteString(Inspect(t.UpperBound.Get()))
	}

	return buffer.String()
}

func (t *TypeParameter) InspectSignatureWithColor() string {
	return lexer.Colorize(t.InspectSignature())
}

func RequiredTypeParameters(typeParams []Ref[*TypeParameter]) int {
	var counter int
	for _, typeParam := range typeParams {
		if !typeParam.Get().Default.IsZero() {
			return counter
		}
		counter++
	}

	return counter
}
