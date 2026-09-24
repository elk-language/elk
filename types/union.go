package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/ds"
)

// Union type represents a list of types.
// A value has to satisfy at least one of the types.
type Union struct {
	Elements []Ref[Type]
	id       ID
}

func NewUnion(elements ...Type) *Union {
	elementRefs := ds.MapSlice(elements, func(element Type) Ref[Type] {
		return ToRef(element)
	})

	return &Union{
		Elements: elementRefs,
	}
}

func NewUnionRef(elements ...Ref[Type]) *Union {
	return &Union{
		Elements: elements,
	}
}

func (u *Union) ToRef() Ref[*Union] {
	return ToRef(u)
}

func (u *Union) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("inter:")
	b := make([]byte, 4)
	for _, elemRef := range u.Elements {
		binary.LittleEndian.PutUint32(b, uint32(elemRef))
		d.Write(b)
	}

	return d.Sum64()
}

func (u *Union) EqualAny(other any) bool {
	o, ok := other.(*Union)
	if !ok {
		return false
	}

	if u.id > 0 {
		return u.id == o.id
	}

	if len(u.Elements) != len(o.Elements) {
		return false
	}

	for j := range len(u.Elements) {
		iElem := u.Elements[j]
		oElem := o.Elements[j]
		if iElem != oElem {
			return false
		}
	}

	return true
}

func (u *Union) ID() ID {
	return u.id
}

func (u *Union) SetID(id ID) {
	u.id = id
}

func (u *Union) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(u, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(u, parent)
	}

	for _, element := range u.Elements {
		if element.Get().traverse(u, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(u, parent)
}

func (u *Union) ToNonLiteral() Type {
	return u
}

func (*Union) IsLiteral() bool {
	return false
}

func (u *Union) inspect() string {
	var buf strings.Builder
	for i, elementRef := range u.Elements {
		element := elementRef.Get()
		if i != 0 {
			buf.WriteString(" | ")
		}
		var addParens bool
		switch element.(type) {
		case *Union, *Intersection:
			addParens = true
		}

		if addParens {
			buf.WriteRune('(')
		}
		buf.WriteString(Inspect(element))
		if addParens {
			buf.WriteRune(')')
		}
	}
	return buf.String()
}

func (u *Union) Copy() *Union {
	return &Union{
		Elements: u.Elements,
	}
}

func (u *Union) CopyType() Type {
	return u.Copy()
}
