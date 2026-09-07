package value

import (
	"github.com/elk-language/elk/value/ivar"
)

// Maps instance variable names to their indices
type IvarIndices ivar.IvarIndices

func (n IvarIndices) SetIndex(name Symbol, i int) {
	(ivar.IvarIndices)(n).SetIndex(name.Id, i)
}

func (n IvarIndices) GetIndex(name Symbol) int {
	return (ivar.IvarIndices)(n).GetIndex(name.Id)
}

func (n IvarIndices) GetIndexOk(name Symbol) (int, bool) {
	return (ivar.IvarIndices)(n).GetIndexOk(name.Id)
}

func (n IvarIndices) GetName(index int) Symbol {
	return n.GetName(index)
}

func (n IvarIndices) GetNameOk(index int) (s Symbol, ok bool) {
	sym, ok := (ivar.IvarIndices)(n).GetNameOk(index)
	if ok {
		return S(sym), ok
	}

	return s, ok
}

func (n *IvarIndices) Copy() Reference {
	return n
}

func (i *IvarIndices) ToValue() Value {
	return Ref(i)
}

func (*IvarIndices) Class() *Class {
	return nil
}

func (*IvarIndices) DirectClass() *Class {
	return nil
}

func (*IvarIndices) SingletonClass() *Class {
	return nil
}

func (in *IvarIndices) Length() int {
	return (*ivar.IvarIndices)(in).Length()
}

const MAX_IVAR_INDICES_ELEMENTS_IN_INSPECT = 300

func (in *IvarIndices) Inspect() string {
	return (*ivar.IvarIndices)(in).Inspect()
}

func (n *IvarIndices) Error() string {
	return n.Inspect()
}

func (*IvarIndices) InstanceVariables() *InstanceVariables {
	return nil
}
