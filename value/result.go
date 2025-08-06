package value

import (
	"fmt"
)

var ResultClass *Class // ::Std::Result

type Result struct {
	ok    bool
	value Value
}

// Create a new successful result
func NewOkResult(value Value) *Result {
	return &Result{
		ok:    true,
		value: value,
	}
}

// Create a new failed result
func NewErrResult(err Value) *Result {
	return &Result{
		ok:    false,
		value: err,
	}
}

// Get the value and the error
func (r *Result) Get() (value, err Value) {
	if r.ok {
		return r.value, Undefined
	}

	return Undefined, r.value
}

func (r *Result) Ok() bool {
	return r.ok
}

// Get the value
func (r *Result) Value() Value {
	if !r.ok {
		return Nil
	}
	return r.value
}

// Get the error
func (r *Result) Err() Value {
	if r.ok {
		return Nil
	}
	return r.value
}

func (r *Result) Copy() Reference {
	return &Result{
		ok:    r.ok,
		value: r.value,
	}
}

func (*Result) Class() *Class {
	return ResultClass
}

func (*Result) DirectClass() *Class {
	return ResultClass
}

func (*Result) SingletonClass() *Class {
	return nil
}

func (r *Result) Inspect() string {
	return fmt.Sprintf("Std::Result{value: %s, err: %s}", r.Value().Inspect(), r.Err().Inspect())
}

func (r *Result) Error() string {
	return r.Inspect()
}

func (*Result) InstanceVariables() *InstanceVariables {
	return nil
}

func initResult() {
	ResultClass = NewClass()
	StdModule.AddConstantString("Result", Ref(ResultClass))
}
