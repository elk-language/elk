package value

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"time"
	_ "time/tzdata" // timezone database

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/ds"
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
	return t.ToGoLocation() == time.UTC
}

func (t *Timezone) IsLocal() bool {
	return t.ToGoLocation() == time.Local
}

// Returns the standard and DST offsets.
func (t *Timezone) Offsets() (TimeSpan, TimeSpan) {
	year := time.Now().Year()
	offsetSet := ds.MakeSet[TimeSpan]()

	for month := 1; month <= 12; month++ {
		dt := MakeDateTime(year, month, 15, 12, 0, 0, 0, 0, 0, t)
		offset := dt.ZoneOffset()
		offsetSet.Add(offset)
		if len(offsetSet) == 2 {
			break
		}
	}

	offsets := slices.Collect(maps.Keys(offsetSet))
	// No DST: only one offset
	if len(offsets) == 1 {
		return offsets[0], offsets[0]
	}

	// Two offsets: determine which one is standard and which is DST
	// DST offset is always the larger
	if offsets[0] < offsets[1] {
		return offsets[0], offsets[1]
	}
	return offsets[1], offsets[0]
}

func (t *Timezone) Equal(other Value) bool {
	otherTz, ok := other.SafeAsReference().(*Timezone)
	if !ok {
		return false
	}

	return t == otherTz ||
		t.Name() == otherTz.Name() &&
			t.StandardOffset() == otherTz.StandardOffset() &&
			t.DSTOffset() == otherTz.DSTOffset()
}

func (t *Timezone) StandardOffset() TimeSpan {
	offset, _ := t.Offsets()
	return offset
}

func (t *Timezone) DSTOffset() TimeSpan {
	_, offset := t.Offsets()
	return offset
}

// Create a new Timezone object.
func NewTimezone(loc *time.Location) *Timezone {
	return (*Timezone)(loc)
}

func NewFixedTimezone(hour, minute, second int) *Timezone {
	offset := TimeSpan(hour)*Hour + TimeSpan(minute)*Minute + TimeSpan(second)*Second
	if offset == 0 {
		return UTCTimezone
	}

	var sign rune
	if offset < 0 {
		sign = '-'
	} else {
		sign = '+'
	}

	name := fmt.Sprintf("UTC%c%02d:%02d:%02d", sign, hour, minute, second)
	return NewTimezone(time.FixedZone(name, int(offset/Second)))
}

func MustLoadTimezone(name string) *Timezone {
	tz, err := LoadTimezone(name)
	if !err.IsUndefined() {
		panic(err)
	}

	return tz
}

var tzAbbrevOffsets = map[string]TimeSpan{
	"ACDT":  10*Hour + 30*Minute,
	"ACST":  9*Hour + 30*Minute,
	"ACWST": 8*Hour + 45*Minute,
	"ADT":   -3 * Hour,
	"AEDT":  11 * Hour,
	"AEST":  10 * Hour,
	"AET":   10 * Hour,
	"AFT":   4*Hour + 30*Minute,
	"AKDT":  -8 * Hour,
	"AKST":  -9 * Hour,
	"ALMT":  6 * Hour,
	"AMST":  -3 * Hour,
	"ANAT":  12 * Hour,
	"AQTT":  5 * Hour,
	"ART":   -3 * Hour,
	"AWST":  8 * Hour,
	"AZOST": 0,
	"AZOT":  -1 * Hour,
	"AZT":   4 * Hour,
	"BNT":   8 * Hour,
	"BIOT":  6 * Hour,
	"BIT":   -12 * Hour,
	"BOT":   -4 * Hour,
	"BRST":  -2 * Hour,
	"BRT":   -3 * Hour,
	"BTT":   6 * Hour,
	"CAT":   2 * Hour,
	"CCT":   6*Hour + 30*Minute,
	"CEST":  2 * Hour,
	"CET":   1 * Hour,
	"CHADT": 13*Hour + 45*Minute,
	"CHAST": 12*Hour + 45*Minute,
	"CHOT":  8 * Hour,
	"CHOST": 9 * Hour,
	"CHST":  10 * Hour,
	"CHUT":  10 * Hour,
	"CIST":  -8 * Hour,
	"CKT":   -10 * Hour,
	"CLST":  -3 * Hour,
	"CLT":   -4 * Hour,
	"COST":  -4 * Hour,
	"COT":   -5 * Hour,
	"CT":    -6 * Hour,
	"CVT":   -1 * Hour,
	"CWST":  8*Hour + 45*Minute,
	"CXT":   7 * Hour,
	"DAVT":  7 * Hour,
	"DDUT":  10 * Hour,
	"DFT":   1 * Hour,
	"EASST": -5 * Hour,
	"EAST":  -6 * Hour,
	"EAT":   3 * Hour,
	"EDT":   -4 * Hour,
	"EEST":  3 * Hour,
	"EET":   2 * Hour,
	"EGST":  0,
	"EGT":   -1 * Hour,
	"EST":   -5 * Hour,
	"ET":    -5 * Hour,
	"FET":   3 * Hour,
	"FJT":   12 * Hour,
	"FKST":  -3 * Hour,
	"FKT":   -4 * Hour,
	"FNT":   -2 * Hour,
	"GALT":  -6 * Hour,
	"GAMT":  -9 * Hour,
	"GET":   4 * Hour,
	"GFT":   -3 * Hour,
	"GILT":  12 * Hour,
	"GIT":   -9 * Hour,
	"GMT":   0,
	"GYT":   -4 * Hour,
	"HDT":   -9 * Hour,
	"HAEC":  2 * Hour,
	"HST":   -10 * Hour,
	"HKT":   8 * Hour,
	"HMT":   5 * Hour,
	"HOVST": 8 * Hour,
	"HOVT":  7 * Hour,
	"ICT":   7 * Hour,
	"IDLW":  -12 * Hour,
	"IDT":   3 * Hour,
	"IOT":   6 * Hour,
	"IRDT":  4*Hour + 30*Minute,
	"IRKT":  8 * Hour,
	"IRST":  3*Hour + 30*Minute,
	"JST":   9 * Hour,
	"KALT":  2 * Hour,
	"KGT":   6 * Hour,
	"KOST":  11 * Hour,
	"KRAT":  7 * Hour,
	"KST":   9 * Hour,
	"LINT":  14 * Hour,
	"MAGT":  12 * Hour,
	"MART":  -9*Hour - 30*Minute,
	"MAWT":  5 * Hour,
	"MDT":   -6 * Hour,
	"MET":   1 * Hour,
	"MEST":  2 * Hour,
	"MHT":   12 * Hour,
	"MIST":  11 * Hour,
	"MIT":   -9*Hour - 30*Minute,
	"MMT":   6*Hour + 30*Minute,
	"MSK":   3 * Hour,
	"MT":    -7 * Hour,
	"MUT":   4 * Hour,
	"MVT":   5 * Hour,
	"MYT":   8 * Hour,
	"NCT":   11 * Hour,
	"NDT":   -2*Hour - 30*Minute,
	"NFT":   11 * Hour,
	"NOVT":  7 * Hour,
	"NPT":   5*Hour + 45*Minute,
	"NST":   -3*Hour - 30*Minute,
	"NT":    -3*Hour - 30*Minute,
	"NUT":   -11 * Hour,
	"NZDT":  13 * Hour,
	"NZDST": 13 * Hour,
	"NZST":  12 * Hour,
	"OMST":  6 * Hour,
	"ORAT":  5 * Hour,
	"PDT":   -7 * Hour,
	"PET":   -5 * Hour,
	"PETT":  12 * Hour,
	"PGT":   10 * Hour,
	"PHOT":  13 * Hour,
	"PHT":   8 * Hour,
	"PHST":  8 * Hour,
	"PKT":   5 * Hour,
	"PMDT":  -2 * Hour,
	"PMST":  -3 * Hour,
	"PONT":  11 * Hour,
	"PST":   -8 * Hour,
	"PT":    -8 * Hour,
	"PWT":   9 * Hour,
	"PYST":  -3 * Hour,
	"PYT":   -4 * Hour,
	"RET":   4 * Hour,
	"ROTT":  -3 * Hour,
	"SAKT":  11 * Hour,
	"SAMT":  4 * Hour,
	"SAST":  2 * Hour,
	"SBT":   11 * Hour,
	"SCT":   4 * Hour,
	"SDT":   -10 * Hour,
	"SGT":   8 * Hour,
	"SLST":  5*Hour + 30*Minute,
	"SRET":  11 * Hour,
	"SRT":   -3 * Hour,
	"SST":   -11 * Hour,
	"SYOT":  3 * Hour,
	"TAHT":  -10 * Hour,
	"THA":   7 * Hour,
	"TFT":   5 * Hour,
	"TJT":   5 * Hour,
	"TKT":   13 * Hour,
	"TLT":   9 * Hour,
	"TMT":   5 * Hour,
	"TRT":   3 * Hour,
	"TOT":   13 * Hour,
	"TST":   8 * Hour,
	"TVT":   12 * Hour,
	"ULAST": 9 * Hour,
	"ULAT":  8 * Hour,
	"UTC":   0,
	"UYST":  -2 * Hour,
	"UYT":   -3 * Hour,
	"UZT":   5 * Hour,
	"VET":   -4 * Hour,
	"VLAT":  10 * Hour,
	"VOLT":  3 * Hour,
	"VOST":  6 * Hour,
	"VUT":   11 * Hour,
	"WAKT":  12 * Hour,
	"WAST":  2 * Hour,
	"WAT":   1 * Hour,
	"WEST":  1 * Hour,
	"WET":   0,
	"WIB":   7 * Hour,
	"WIT":   9 * Hour,
	"WITA":  8 * Hour,
	"WGST":  -2 * Hour,
	"WGT":   -3 * Hour,
	"WST":   8 * Hour,
	"YAKT":  9 * Hour,
	"YEKT":  5 * Hour,
}

var utcOffsetTzRegex = MustCompileRegex("^UTC(\\+|\\-)(\\d{2}):(\\d{2})(:(\\d{2}))?$", bitfield.BitField8{})

func loadOffsetTimezone(name string, matches [][]byte) (zone *Timezone, err error) {
	sign := matches[1]
	var positive TimeSpan
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

	offset := positive * (TimeSpan(hour)*Hour + TimeSpan(minute)*Minute + TimeSpan(second)*Second)
	if offset == 0 {
		return UTCTimezone, Undefined
	}
	if offset >= 24*Hour || offset <= -24*Hour {
		return nil, errors.New("offset out of range")
	}
	loc := time.FixedZone(name, int(offset/Second))
	return NewTimezone(loc), err
}

// Create a new timezone based on a fixed offset and return an error for invalid values
func NewTimezoneFromOffsetErr(offset TimeSpan) (*Timezone, Value) {
	if offset >= 24*Hour || offset <= -24*Hour {
		return nil, Ref(
			Errorf(
				OutOfRangeErrorClass,
				"invalid timezone offset: `%s`",
				offset.Inspect(),
			),
		)
	}

	return NewTimezoneFromOffset(offset), Undefined
}

// Create a new timezone based on a fixed offset.
func NewTimezoneFromOffset(offset TimeSpan) *Timezone {
	if offset == 0 {
		return UTCTimezone
	}

	var sign rune
	var rest TimeSpan
	if offset < 0 {
		sign = '-'
		rest = -offset
	} else {
		sign = '+'
		rest = offset
	}

	hours := rest / Hour
	rest = rest % Hour

	minutes := rest / Minute
	rest = rest % Minute

	seconds := rest / Second
	rest = rest % Second

	name := fmt.Sprintf("UTC%c%02d:%02d:%02d", sign, hours, minutes, seconds)
	return NewTimezone(time.FixedZone(name, int(offset/Second)))
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
