package value

import (
	"fmt"
	"time"
	_ "time/tzdata" // timezone database
)

var TimeClass *Class // ::Std::Time

// Elk's Time value
type Time struct {
	time.Time
}

func (Time) Class() *Class {
	return TimeClass
}

func (Time) DirectClass() *Class {
	return TimeClass
}

func (Time) SingletonClass() *Class {
	return nil
}

func (t Time) Inspect() string {
	return fmt.Sprintf("Time('%s')", t.Format(time.RFC3339Nano))
}

func (t Time) InstanceVariables() SymbolMap {
	return nil
}

func TimeNow() Time {
	return Time{Time: time.Now()}
}

func initTime() {
	TimeClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Time", TimeClass)
}
