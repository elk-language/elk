package value

import (
	"fmt"
	"strconv"
	"time"
	_ "time/tzdata" // timezone database

	"github.com/elk-language/elk/bitfield"
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

func (t *Timezone) InstanceVariables() *InstanceVariables {
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

func NewFixedTimezone(hour, minute, second int) *Timezone {
	offset := hour*tzHour + minute*tzMinute + second*tzSecond

	var sign rune
	if offset < 0 {
		sign = '-'
	} else {
		sign = '+'
	}

	name := fmt.Sprintf("UTC%c%02d:%02d:%02d", sign, hour, minute, second)
	return NewTimezone(time.FixedZone(name, offset))
}

func MustLoadTimezone(name string) *Timezone {
	tz, err := LoadTimezone(name)
	if !err.IsUndefined() {
		panic(err)
	}

	return tz
}

const (
	tzSecond = 1
	tzMinute = tzSecond * 60
	tzHour   = tzMinute * 60
)

var tzAbbrevOffsets = map[string]int{
	"ACDT":  10*tzHour + 30*tzMinute,
	"ACST":  9*tzHour + 30*tzMinute,
	"ACWST": 8*tzHour + 45*tzMinute,
	"ADT":   -3 * tzHour,
	"AEDT":  11 * tzHour,
	"AEST":  10 * tzHour,
	"AET":   10 * tzHour,
	"AFT":   4*tzHour + 30*tzMinute,
	"AKDT":  -8 * tzHour,
	"AKST":  -9 * tzHour,
	"ALMT":  6 * tzHour,
	"AMST":  -3 * tzHour,
	"ANAT":  12 * tzHour,
	"AQTT":  5 * tzHour,
	"ART":   -3 * tzHour,
	"AWST":  8 * tzHour,
	"AZOST": 0,
	"AZOT":  -1 * tzHour,
	"AZT":   4 * tzHour,
	"BNT":   8 * tzHour,
	"BIOT":  6 * tzHour,
	"BIT":   -12 * tzHour,
	"BOT":   -4 * tzHour,
	"BRST":  -2 * tzHour,
	"BRT":   -3 * tzHour,
	"BTT":   6 * tzHour,
	"CAT":   2 * tzHour,
	"CCT":   6*tzHour + 30*tzMinute,
	"CEST":  2 * tzHour,
	"CET":   1 * tzHour,
	"CHADT": 13*tzHour + 45*tzMinute,
	"CHAST": 12*tzHour + 45*tzMinute,
	"CHOT":  8 * tzHour,
	"CHOST": 9 * tzHour,
	"CHST":  10 * tzHour,
	"CHUT":  10 * tzHour,
	"CIST":  -8 * tzHour,
	"CKT":   -10 * tzHour,
	"CLST":  -3 * tzHour,
	"CLT":   -4 * tzHour,
	"COST":  -4 * tzHour,
	"COT":   -5 * tzHour,
	"CT":    -6 * tzHour,
	"CVT":   -1 * tzHour,
	"CWST":  8*tzHour + 45*tzMinute,
	"CXT":   7 * tzHour,
	"DAVT":  7 * tzHour,
	"DDUT":  10 * tzHour,
	"DFT":   1 * tzHour,
	"EASST": -5 * tzHour,
	"EAST":  -6 * tzHour,
	"EAT":   3 * tzHour,
	"EDT":   -4 * tzHour,
	"EEST":  3 * tzHour,
	"EET":   2 * tzHour,
	"EGST":  0,
	"EGT":   -1 * tzHour,
	"EST":   -5 * tzHour,
	"ET":    -5 * tzHour,
	"FET":   3 * tzHour,
	"FJT":   12 * tzHour,
	"FKST":  -3 * tzHour,
	"FKT":   -4 * tzHour,
	"FNT":   -2 * tzHour,
	"GALT":  -6 * tzHour,
	"GAMT":  -9 * tzHour,
	"GET":   4 * tzHour,
	"GFT":   -3 * tzHour,
	"GILT":  12 * tzHour,
	"GIT":   -9 * tzHour,
	"GMT":   0,
	"GYT":   -4 * tzHour,
	"HDT":   -9 * tzHour,
	"HAEC":  2 * tzHour,
	"HST":   -10 * tzHour,
	"HKT":   8 * tzHour,
	"HMT":   5 * tzHour,
	"HOVST": 8 * tzHour,
	"HOVT":  7 * tzHour,
	"ICT":   7 * tzHour,
	"IDLW":  -12 * tzHour,
	"IDT":   3 * tzHour,
	"IOT":   6 * tzHour,
	"IRDT":  4*tzHour + 30*tzMinute,
	"IRKT":  8 * tzHour,
	"IRST":  3*tzHour + 30*tzMinute,
	"JST":   9 * tzHour,
	"KALT":  2 * tzHour,
	"KGT":   6 * tzHour,
	"KOST":  11 * tzHour,
	"KRAT":  7 * tzHour,
	"KST":   9 * tzHour,
	"LINT":  14 * tzHour,
	"MAGT":  12 * tzHour,
	"MART":  -9*tzHour - 30*tzMinute,
	"MAWT":  5 * tzHour,
	"MDT":   -6 * tzHour,
	"MET":   1 * tzHour,
	"MEST":  2 * tzHour,
	"MHT":   12 * tzHour,
	"MIST":  11 * tzHour,
	"MIT":   -9*tzHour - 30*tzMinute,
	"MMT":   6*tzHour + 30*tzMinute,
	"MSK":   3 * tzHour,
	"MT":    -7 * tzHour,
	"MUT":   4 * tzHour,
	"MVT":   5 * tzHour,
	"MYT":   8 * tzHour,
	"NCT":   11 * tzHour,
	"NDT":   -2*tzHour - 30*tzMinute,
	"NFT":   11 * tzHour,
	"NOVT":  7 * tzHour,
	"NPT":   5*tzHour + 45*tzMinute,
	"NST":   -3*tzHour - 30*tzMinute,
	"NT":    -3*tzHour - 30*tzMinute,
	"NUT":   -11 * tzHour,
	"NZDT":  13 * tzHour,
	"NZDST": 13 * tzHour,
	"NZST":  12 * tzHour,
	"OMST":  6 * tzHour,
	"ORAT":  5 * tzHour,
	"PDT":   -7 * tzHour,
	"PET":   -5 * tzHour,
	"PETT":  12 * tzHour,
	"PGT":   10 * tzHour,
	"PHOT":  13 * tzHour,
	"PHT":   8 * tzHour,
	"PHST":  8 * tzHour,
	"PKT":   5 * tzHour,
	"PMDT":  -2 * tzHour,
	"PMST":  -3 * tzHour,
	"PONT":  11 * tzHour,
	"PST":   -8 * tzHour,
	"PT":    -8 * tzHour,
	"PWT":   9 * tzHour,
	"PYST":  -3 * tzHour,
	"PYT":   -4 * tzHour,
	"RET":   4 * tzHour,
	"ROTT":  -3 * tzHour,
	"SAKT":  11 * tzHour,
	"SAMT":  4 * tzHour,
	"SAST":  2 * tzHour,
	"SBT":   11 * tzHour,
	"SCT":   4 * tzHour,
	"SDT":   -10 * tzHour,
	"SGT":   8 * tzHour,
	"SLST":  5*tzHour + 30*tzMinute,
	"SRET":  11 * tzHour,
	"SRT":   -3 * tzHour,
	"SST":   -11 * tzHour,
	"SYOT":  3 * tzHour,
	"TAHT":  -10 * tzHour,
	"THA":   7 * tzHour,
	"TFT":   5 * tzHour,
	"TJT":   5 * tzHour,
	"TKT":   13 * tzHour,
	"TLT":   9 * tzHour,
	"TMT":   5 * tzHour,
	"TRT":   3 * tzHour,
	"TOT":   13 * tzHour,
	"TST":   8 * tzHour,
	"TVT":   12 * tzHour,
	"ULAST": 9 * tzHour,
	"ULAT":  8 * tzHour,
	"UTC":   0,
	"UYST":  -2 * tzHour,
	"UYT":   -3 * tzHour,
	"UZT":   5 * tzHour,
	"VET":   -4 * tzHour,
	"VLAT":  10 * tzHour,
	"VOLT":  3 * tzHour,
	"VOST":  6 * tzHour,
	"VUT":   11 * tzHour,
	"WAKT":  12 * tzHour,
	"WAST":  2 * tzHour,
	"WAT":   1 * tzHour,
	"WEST":  1 * tzHour,
	"WET":   0,
	"WIB":   7 * tzHour,
	"WIT":   9 * tzHour,
	"WITA":  8 * tzHour,
	"WGST":  -2 * tzHour,
	"WGT":   -3 * tzHour,
	"WST":   8 * tzHour,
	"YAKT":  9 * tzHour,
	"YEKT":  5 * tzHour,
}

var utcOffsetTzRegex = MustCompileRegex("^UTC(\\+|\\-)(\\d{2}):(\\d{2})(:(\\d{2}))?$", bitfield.BitField8{})

func loadOffsetTimezone(name string, matches [][]byte) (zone *Timezone, err error) {
	sign := matches[1]
	var positive int
	switch sign[0] {
	case '+':
		positive = 1
	default:
		positive = -1
	}

	hourStr := matches[2]
	hour, err := strconv.Atoi(string(hourStr))
	if err != nil {
		return nil, err
	}

	minuteStr := matches[3]
	minute, err := strconv.Atoi(string(minuteStr))
	if err != nil {
		return nil, err
	}

	var second int
	if len(matches) > 5 {
		secondStr := matches[5]
		if len(secondStr) > 0 {
			second, err = strconv.Atoi(string(secondStr))
			if err != nil {
				return nil, err
			}
		}
	}

	offset := positive * (hour*tzHour + minute*tzMinute + second*tzSecond)
	loc := time.FixedZone(name, offset)
	return NewTimezone(loc), err
}

// Create a new timezone based on a fixed offset.
func NewTimezoneFromOffset(offset int) *Timezone {
	var sign rune
	var rest int
	if offset < 0 {
		sign = '-'
		rest = -offset
	} else {
		sign = '+'
		rest = offset
	}

	hours := rest / tzHour
	rest = rest % tzHour

	minutes := rest / tzMinute
	rest = rest % tzMinute

	seconds := rest / tzSecond
	rest = rest % tzSecond

	name := fmt.Sprintf("UTC%c%02d:%02d:%02d", sign, hours, minutes, seconds)
	return NewTimezone(time.FixedZone(name, offset))
}

// Load a timezone from the IANA database.
func LoadTimezone(name string) (zone *Timezone, err Value) {
	matches := utcOffsetTzRegex.Re.FindSubmatch([]byte(name))
	if matches != nil {
		zone, er := loadOffsetTimezone(name, matches)
		if er != nil {
			return nil, Ref(Errorf(
				InvalidTimezoneErrorClass,
				"invalid timezone: %s",
				name,
			))
		}
		return zone, Undefined
	}
	loc, er := time.LoadLocation(name)
	if er != nil {
		return nil, Ref(Errorf(
			InvalidTimezoneErrorClass,
			"invalid timezone: %s",
			name,
		))
	}

	return NewTimezone(loc), Undefined
}

func initTimezone() {
	TimezoneClass = NewClass()
	StdModule.AddConstantString("Timezone", Ref(TimezoneClass))
	TimezoneClass.AddConstantString("UTC", Ref(UTCTimezone))
	TimezoneClass.AddConstantString("LOCAL", Ref(LocalTimezone))
}
