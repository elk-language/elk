package value

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var StackTraceClass *Class         // ::Std::StackTrace
var StackTraceIteratorClass *Class // ::Std::StackTrace::Iterator

type StackTrace []CallFrame

func (s *StackTrace) Copy() Reference {
	return s
}

func (*StackTrace) Class() *Class {
	return StackTraceClass
}

func (*StackTrace) DirectClass() *Class {
	return StackTraceClass
}

func (*StackTrace) SingletonClass() *Class {
	return nil
}

func (s *StackTrace) Inspect() string {
	return fmt.Sprintf("Std::StackTrace{&: %p}", s)
}

func (s *StackTrace) Error() string {
	return s.String()
}

func (*StackTrace) InstanceVariables() SymbolMap {
	return nil
}

func (s *StackTrace) Length() int {
	return len(*s)
}
func (s *StackTrace) Get(i int) (*CallFrame, Value) {
	l := len(*s)
	if i >= l || i < -l {
		return nil, Ref(NewIndexOutOfRangeError(fmt.Sprint(i), l))
	}

	if i < 0 {
		i = l + i
	}

	return &(*s)[i], Undefined
}

func (s *StackTrace) String() string {
	var buffer strings.Builder
	buffer.WriteString("Stack trace (the most recent call is last)\n")

	for i, callFrame := range *s {
		if callFrame.TailCallCounter > 0 {
			fmt.Fprintf(&buffer, " ... %d optimised tail call(s)\n", callFrame.TailCallCounter)
		}
		// "  %d: %s:%d, in `%s`\n"
		fmt.Fprint(&buffer, " ")
		color.New(color.FgHiBlue).Fprintf(&buffer, "%d", i)
		fmt.Fprintf(&buffer, ": %s:%d, in ", callFrame.FileName, callFrame.LineNumber)
		color.New(color.FgHiYellow).Fprintf(&buffer, "`%s`", callFrame.FuncName)
		fmt.Fprintln(&buffer)
	}
	// Stack trace (the most recent call is last):
	//   0: /tmp/test.elk:18, in `foo`
	//   1: /tmp/test.elk:11, in `bar`

	return buffer.String()
}

type StackTraceIterator struct {
	StackTrace *StackTrace
	Index      int
}

func NewStackTraceIterator(stackTrace *StackTrace) *StackTraceIterator {
	return &StackTraceIterator{
		StackTrace: stackTrace,
	}
}

func NewStackTraceIteratorWithIndex(stackTrace *StackTrace, index int) *StackTraceIterator {
	return &StackTraceIterator{
		StackTrace: stackTrace,
		Index:      index,
	}
}

func (*StackTraceIterator) Class() *Class {
	return StackTraceIteratorClass
}

func (*StackTraceIterator) DirectClass() *Class {
	return StackTraceIteratorClass
}

func (*StackTraceIterator) SingletonClass() *Class {
	return nil
}

func (s *StackTraceIterator) Copy() Reference {
	return &StackTraceIterator{
		StackTrace: s.StackTrace,
		Index:      s.Index,
	}
}

func (s *StackTraceIterator) Inspect() string {
	return fmt.Sprintf("Std::StackTrace::Iterator{&: %p, stack_trace: %s, index: %d}", s, s.StackTrace.Inspect(), s.Index)
}

func (s *StackTraceIterator) Error() string {
	return s.Inspect()
}

func (*StackTraceIterator) InstanceVariables() SymbolMap {
	return nil
}

func (s *StackTraceIterator) Next() (Value, Value) {
	if s.Index >= s.StackTrace.Length() {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next := &(*s.StackTrace)[s.Index]
	s.Index++
	return Ref(next), Undefined
}

func (s *StackTraceIterator) Reset() {
	s.Index = 0
}

func initStackTrace() {
	StackTraceClass = NewClassWithOptions()
	StdModule.AddConstantString("StackTrace", Ref(StackTraceClass))

	StackTraceIteratorClass = NewClassWithOptions()
	StackTraceClass.AddConstantString("Iterator", Ref(StackTraceIteratorClass))
}
