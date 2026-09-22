package types

import (
	"encoding/binary"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/ds"
)

// Intersection type represents a list of types.
// A value has to satisfy all of the types.
type Intersection struct {
	Elements []Ref[Type]
	id       ID
}

func NewIntersectionRef(elements ...Ref[Type]) *Intersection {
	return &Intersection{
		Elements: elements,
	}
}

func NewIntersection(elements ...Type) *Intersection {
	elementRefs := ds.MapSlice(elements, func(element Type) Ref[Type] {
		return ToRef(element)
	})

	return &Intersection{
		Elements: elementRefs,
	}
}

func (i *Intersection) ToRef() Ref[*Intersection] {
	return Ref[*Intersection](i.id)
}

func (i *Intersection) HashUint64() uint64 {
	d := xxhash.New()

	d.WriteString("inter:")
	b := make([]byte, 4)
	for _, elemRef := range i.Elements {
		binary.LittleEndian.PutUint32(b, uint32(elemRef))
		d.Write(b)
	}

	return d.Sum64()
}

func (i *Intersection) EqualAny(other any) bool {
	o, ok := other.(*Intersection)
	if !ok {
		return false
	}

	if i.id > 0 {
		return i.id == o.id
	}

	if len(i.Elements) != len(o.Elements) {
		return false
	}

	for j := range len(i.Elements) {
		iElem := i.Elements[j]
		oElem := o.Elements[j]
		if iElem != oElem {
			return false
		}
	}

	return true
}

func (s *Intersection) ID() ID {
	return s.id
}

func (s *Intersection) SetID(id ID) {
	s.id = id
}

func (i *Intersection) traverse(parent Type, enter func(node, parent Type) TraverseOption, leave func(node, parent Type) TraverseOption) TraverseOption {
	switch enter(i, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseContinue:
		return leave(i, parent)
	}

	for _, element := range i.Elements {
		if element.Get().traverse(i, enter, leave) == TraverseBreak {
			return TraverseBreak
		}
	}

	return leave(i, parent)
}

func (u *Intersection) ToNonLiteral() Type {
	return u
}

func (*Intersection) IsLiteral() bool {
	return false
}

func (u *Intersection) inspect() string {
	var buf strings.Builder
	for i, elementRef := range u.Elements {
		element := elementRef.Get()
		if i != 0 {
			buf.WriteString(" & ")
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

func (i *Intersection) Copy() *Intersection {
	return &Intersection{
		Elements: i.Elements,
	}
}

func (i *Intersection) CopyType() Type {
	return i.Copy()
}
