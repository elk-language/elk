package value

import (
	"fmt"
	"time"
	_ "time/tzdata" // timezone database
)

var TimezoneClass *Class // ::Std::Timezone

// Elk's Timezone value
type Timezone time.Location

func (*Timezone) Class() *Class {
	return TimezoneClass
}

func (*Timezone) DirectClass() *Class {
	return TimezoneClass
}

func (*Timezone) SingletonClass() *Class {
	return nil
}

func (t *Timezone) Inspect() string {
	return fmt.Sprintf("Timezone('%s')", t.Name())
}

func (t *Timezone) InstanceVariables() SymbolMap {
	return nil
}

func (t *Timezone) ToGoLocation() *time.Location {
	return (*time.Location)(t)
}

func (t *Timezone) Name() string {
	return t.ToGoLocation().String()
}

// Create a new Timezone object.
func NewTimezone(loc *time.Location) *Timezone {
	return (*Timezone)(loc)
}

func MustLoadTimezone(name string) *Timezone {
	tz, err := LoadTimezone(name)
	if err != nil {
		panic(err)
	}

	return tz
}

// Load a timezone from the IANA database.
func LoadTimezone(name string) (*Timezone, *Error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, Errorf(
			InvalidTimezoneErrorClass,
			"invalid timezone: %s",
			name,
		)
	}

	return NewTimezone(loc), nil
}

func initTimezone() {
	TimezoneClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Timezone", TimezoneClass)
}
