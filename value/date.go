package value

import (
	"fmt"
)

const dateYearBias = 1 << 22

const DateMaxYear = (1 << 22) - 1
const DateMinYear = -(1 << 22)

const DateMaxMonth = 12
const DateMinMonth = 1

const DateMaxDay = 31
const DateMinDay = 1

var DateClass *Class                  // ::Std::Date
var DateErrorClass *Class             // ::Std::Date::Error
var DateInvalidYearErrorClass *Class  // ::Std::Date::InvalidYearError
var DateInvalidMonthErrorClass *Class // ::Std::Date::InvalidMonthError
var DateInvalidDayErrorClass *Class   // ::Std::Date::InvalidDayError

// Represents a calendar date (year, month, day).
// It is an inline value on both 32bit and 64bit systems.
// The year range is from `-4_194_304` to `4_194_303`
type Date struct {
	bits uint32
}

// Make and validate a new date
func MakeValidatedDate(year, month, day int) (Date, Value) {
	if year > DateMaxYear || year < DateMinYear {
		return Date{}, Ref(
			Errorf(
				DateInvalidYearErrorClass,
				"year %d is out of range %d...%d",
				year,
				DateMinYear,
				DateMaxYear,
			),
		)
	}
	if month > DateMaxMonth || month < DateMinMonth {
		return Date{}, Ref(
			Errorf(
				DateInvalidMonthErrorClass,
				"month %d is out of range %d...%d",
				month,
				DateMinMonth,
				DateMaxMonth,
			),
		)
	}
	if day > DateMaxDay || day < DateMinDay {
		return Date{}, Ref(
			Errorf(
				DateInvalidDayErrorClass,
				"day %d is out of range %d...%d",
				day,
				DateMinDay,
				DateMaxDay,
			),
		)
	}

	return MakeDate(year, month, day), Undefined
}

// Make a new date
func MakeDate(year, month, day int) Date {
	y := uint32(year + dateYearBias)
	m := uint32(month)
	d := uint32(day)
	bits := (y << 9) | (m << 5) | d

	return Date{
		bits: bits,
	}
}

func (d Date) ToValue() Value {
	return Value{
		flag: DATE_FLAG,
		data: uintptr(d.bits),
	}
}

func (d Date) Copy() Reference {
	return d
}

func (Date) Class() *Class {
	return DateClass
}

func (Date) DirectClass() *Class {
	return DateClass
}

func (Date) SingletonClass() *Class {
	return nil
}

func (d Date) Inspect() string {
	return fmt.Sprintf("Std::Date(%d, %d, %d)", d.Year(), d.Month(), d.Day())
}

func (Date) InstanceVariables() *InstanceVariables {
	return nil
}

func (d Date) Error() string {
	return d.Inspect()
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
}

func (d Date) ToString() String {
	return String(d.String())
}

func (d Date) Year() int {
	y := int32(d.bits >> 9)      // extract top 23 bits
	return int(y - dateYearBias) // remove bias
}

func (d Date) Month() int {
	m := int32((d.bits >> 5) & 0b1111)
	return int(m)
}

func (d Date) Day() int {
	m := int32(d.bits & 0b11111)
	return int(m)
}

func initDate() {
	DateClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Date", Ref(DateClass))

	DateErrorClass = NewClassWithOptions(ClassWithSuperclass(ErrorClass))
	DateClass.AddConstantString("Error", Ref(DateErrorClass))

	DateInvalidYearErrorClass = NewClassWithOptions(ClassWithSuperclass(DateErrorClass))
	DateClass.AddConstantString("InvalidYearError", Ref(DateInvalidYearErrorClass))

	DateInvalidMonthErrorClass = NewClassWithOptions(ClassWithSuperclass(DateErrorClass))
	DateClass.AddConstantString("InvalidMonthError", Ref(DateInvalidMonthErrorClass))

	DateInvalidDayErrorClass = NewClassWithOptions(ClassWithSuperclass(DateErrorClass))
	DateClass.AddConstantString("InvalidDayError", Ref(DateInvalidDayErrorClass))
}
