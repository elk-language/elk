package vm

import (
	"reflect"
	"strings"

	"github.com/elk-language/elk/value"
)

type SelectResult struct {
	ReflectValue    reflect.Value
	Err             value.Value
	ChosenCaseIndex int
	ChannelIsOpen   bool
}

func (r *SelectResult) Value(channel value.AnyChannel) value.Value {
	if channel.IsTransformerChannel() {
		return channel.TransformAnyToValue(r.AnyValue())
	}

	switch channel.(type) {
	case *value.ChannelOfValue, *value.ReadChannelOfValue, *value.WriteChannelOfValue:
		return r.AnyValue().(value.Value)
	default:
		return r.AnyValue().(value.ValueInterface).ToValue()
	}
}

func (r *SelectResult) AnyValue() any {
	return r.ReflectValue.Interface()
}

func (r *SelectResult) Result(channel value.AnyChannel) value.Result {
	if !r.ChannelIsOpen {
		return value.MakeErrResult(value.ChannelClosedPopError.ToValue())
	}

	return value.MakeOkResult(r.Value(channel))
}

func DoSelect(cases []reflect.SelectCase) *SelectResult {
	var channelOpen bool
	var val reflect.Value
	var chosenCaseIndex int
	var elkErr value.Value
	func() {
		// catch panic: send on closed channel
		defer func() {
			panicVal := recover()
			if panicVal == nil {
				return
			}

			err, ok := panicVal.(error)
			if !ok {
				panic(panicVal)
			}
			if err.Error() != "send on closed channel" {
				panic(panicVal)
			}

			elkErr = value.ChannelClosedPushError.ToValue()
		}()

		chosenCaseIndex, val, channelOpen = reflect.Select(cases)
	}()

	if elkErr.IsNotUndefined() {
		return &SelectResult{
			Err: elkErr,
		}
	}
	if chosenCaseIndex == 0 {
		return &SelectResult{
			Err: value.ExecutionAbortedError.ToValue(),
		}
	}

	return &SelectResult{
		ChannelIsOpen:   channelOpen,
		ReflectValue:    val,
		ChosenCaseIndex: chosenCaseIndex - 1,
	}
}

// Wraps data for `select` expressions
type Select struct {
	Cases []SelectCase
	value.ValueBase
}

var _ value.Reference = &Select{}

func NewSelect(cases []SelectCase) *Select {
	return &Select{
		Cases: cases,
	}
}

func NewSelectVar(cases ...SelectCase) *Select {
	return &Select{
		Cases: cases,
	}
}

func (s *Select) ToValue() value.Value {
	return value.Ref(s)
}

func (s *Select) Copy() value.Reference {
	return &Select{
		Cases: s.Cases,
	}
}

func (s *Select) Inspect() string {
	var buff strings.Builder
	buff.WriteString("Select{")
	for i, selectCase := range s.Cases {
		if i != 0 {
			buff.WriteString(", ")
		}

		switch selectCase.Direction {
		case reflect.SelectSend:
			buff.WriteString("send")
		case reflect.SelectRecv:
			buff.WriteString("receive")
		case reflect.SelectDefault:
			buff.WriteString("else")
		default:
			buff.WriteString("unknown")
		}
	}
	buff.WriteString("}")

	return buff.String()
}

func (s *Select) String() string {
	return s.Inspect()
}

func (s *Select) Error() string {
	return s.Inspect()
}

type SelectCase struct {
	Direction reflect.SelectDir
}

func MakeSelectCase(dir reflect.SelectDir) SelectCase {
	return SelectCase{
		Direction: dir,
	}
}
