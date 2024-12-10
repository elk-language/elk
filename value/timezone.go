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

func (t *Timezone) Copy() Reference {
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
	if !err.IsUndefined() {
		panic(err)
	}

	return tz
}

// Load a timezone from the IANA database.
func LoadTimezone(name string) (zone *Timezone, err Value) {
	loc, er := time.LoadLocation(name)
	if er != nil {
		return nil, Ref(Errorf(
			InvalidTimezoneErrorClass,
			"invalid timezone: %s",
			name,
		))
	}

	return NewTimezone(loc), err
}

func initTimezone() {
	TimezoneClass = NewClass()
	StdModule.AddConstantString("Timezone", Ref(TimezoneClass))
	TimezoneClass.AddConstantString("UTC", Ref(UTCTimezone))
	TimezoneClass.AddConstantString("LOCAL", Ref(LocalTimezone))
}
