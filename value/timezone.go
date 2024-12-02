package value

import (
	"fmt"
	"time"
	_ "time/tzdata" // timezone database
)

var TimezoneClass *Class // ::Std::Timezone

// Elk's Timezone value
type Timezone time.Location

var LocalTimezone = NewTimezone(time.Local)
var UTCTimezone = NewTimezone(time.UTC)

func (*Timezone) Class() *Class {
	return TimezoneClass
}

func (*Timezone) DirectClass() *Class {
	return TimezoneClass
}

func (*Timezone) SingletonClass() *Class {
	return nil
}

func (t *Timezone) Copy() Value {
	return t
}

func (t *Timezone) Error() string {
	return t.Inspect()
}

func (t *Timezone) Inspect() string {
	return fmt.Sprintf("Std::Timezone['%s']", t.Name())
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

func (t *Timezone) IsUTC() bool {
	return t.Name() == "UTC"
}

func (t *Timezone) IsLocal() bool {
	return t.Name() == "Local"
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
func LoadTimezone(name string) (*Timezone, Value) {
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
	TimezoneClass = NewClass()
	StdModule.AddConstantString("Timezone", TimezoneClass)
	TimezoneClass.AddConstantString("UTC", UTCTimezone)
	TimezoneClass.AddConstantString("LOCAL", LocalTimezone)
}
