package value

import (
	"fmt"
)

// Represents the difference between two dates
// as a 32 bit number of months and 32 bit number of days.
// Can store up to 2 billion days and 2 billion months (178 million years).
type DateSpan struct {
	months int32
	days   int32
}

func MakeDateSpan(years, months, days int) DateSpan {
	months += years * 12

	return DateSpan{
		months: int32(months),
		days:   int32(days),
	}
}

var DateSpanClass *Class // ::Std::Date::Span

func (d DateSpan) Copy() Reference {
	return d
}

func (DateSpan) Class() *Class {
	return DateSpanClass
}

func (DateSpan) DirectClass() *Class {
	return DateSpanClass
}

func (DateSpan) SingletonClass() *Class {
	return nil
}

func (d DateSpan) Inspect() string {
	return fmt.Sprintf("Std::Date::Span(%d, %d, %d)", d.Years(), d.MonthsMod(), d.Days())
}

func (DateSpan) InstanceVariables() *InstanceVariables {
	return nil
}

func (d DateSpan) Error() string {
	return d.Inspect()
}

func (d DateSpan) String() string {
	years := d.months / 12
	months := d.MonthsMod()
	days := d.days
	return fmt.Sprintf("%dY%dM%dD", years, months, days)
}

func (d DateSpan) ToString() String {
	return String(d.String())
}

func (d DateSpan) Days() int {
	return int(d.days)
}

func (d DateSpan) Months() int {
	return int(d.months)
}

func (d DateSpan) MonthsMod() int {
	return int(d.months % 12)
}

func (d DateSpan) Years() int {
	return int(d.months / 12)
}

func initDateSpan() {
	DateSpanClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	DateClass.AddConstantString("Span", Ref(DateSpanClass))
}
