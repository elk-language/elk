package value

import (
	"fmt"
)

var ResultClass *Class // ::Std::Result

type Result struct {
	value Value
	ok    bool
}

var _ ValueInterface = Result{}
var _ Reference = Result{}

// Create a new result
func MakeResult(value Value, ok bool) Result {
	return Result{
		value: value,
		ok:    ok,
	}
}

// Create a new successful result
func MakeOkResult(value Value) Result {
	return Result{
		ok:    true,
		value: value,
	}
}

// Create a new failed result
func MakeErrResult(err Value) Result {
	return Result{
		ok:    false,
		value: err,
	}
}

// Get the value and the error
func (r Result) Get() (value, err Value) {
	if r.ok {
		return r.value, Undefined
	}

	return Undefined, r.value
}

func (r Result) Ok() bool {
	return r.ok
}

func (r Result) Copy() Reference {
	return r
}

// Get the value
func (r Result) Value() Value {
	if !r.ok {
		return Nil
	}
	return r.value
}

// Get the error
func (r Result) Err() Value {
	if r.ok {
		return Nil
	}
	return r.value
}

func (r Result) ToValue() Value {
	var flag uint8
	if r.ok {
		flag = RESULT_OK_FLAG
	} else {
		flag = RESULT_ERR_FLAG
	}
	value := r.value

	// handle nested Results
	if r.value.IsInlineResult() {
		value = Ref(r.value.AsInlineResult())
	}

	return Value{
		data:        value.data,
		ptr:         value.ptr,
		result_flag: value.flag,
		flag:        flag,
	}
}

func (Result) Class() *Class {
	return ResultClass
}

func (Result) DirectClass() *Class {
	return ResultClass
}

func (Result) SingletonClass() *Class {
	return nil
}

func (r Result) Inspect() string {
	return fmt.Sprintf("Std::Result{value: %s, err: %s}", r.Value().Inspect(), r.Err().Inspect())
}

func (r Result) Error() string {
	return r.Inspect()
}

func (Result) InstanceVariables() *InstanceVariables {
	return nil
}

func initResult() {
	ResultClass = NewClass()
	StdModule.AddConstantString("Result", Ref(ResultClass))
	RegisterNativeClass("Std::Result", "value.ResultClass")
}
