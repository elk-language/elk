package durationscanner

type Token uint8

// Name of the token.
func (t Token) String() string {
	if int(t) > len(tokenNames) {
		return "UNKNOWN"
	}

	return tokenNames[t]
}

const (
	ZERO_VALUE         Token = iota // Zero value for Type
	ERROR                           // Invalid format directive
	END_OF_FILE                     // End Of File has been reached
	YEARS_INT                       // Year component eg. "2Y"
	YEARS_FLOAT                     // Year component eg. "2.5Y"
	MONTHS_INT                      // Month component eg. "15M"
	MONTHS_FLOAT                    // Month component eg. "15.9M"
	DAYS_INT                        // Day component eg. "1D"
	DAYS_FLOAT                      // Day component eg. "1.8D"
	HOURS_INT                       // Hour component eg. "9h"
	HOURS_FLOAT                     // Hour component eg. "9.15h"
	MINUTES_INT                     // Minute component eg. "10m"
	MINUTES_FLOAT                   // Minute component eg. "10.76m"
	SECONDS_INT                     // Second component eg. "50s"
	SECONDS_FLOAT                   // Second component eg. "50.78s"
	MILLISECONDS_INT                // Millisecond component eg. "150ms"
	MILLISECONDS_FLOAT              // Millisecond component eg. "150.8ms"
	MICROSECONDS_INT                // Microsecond component eg. "600us", "600µs"
	MICROSECONDS_FLOAT              // Microsecond component eg. "600.256us", "600µs"
	NANOSECONDS_INT                 // Nanosecond component eg. "750ns"
	NANOSECONDS_FLOAT               // Nanosecond component eg. "750.72ns"
)

var tokenNames = [...]string{
	ERROR:              "ERROR",
	END_OF_FILE:        "END_OF_FILE",
	YEARS_INT:          "YEARS_INT",
	YEARS_FLOAT:        "YEARS_FLOAT",
	MONTHS_INT:         "MONTHS_INT",
	MONTHS_FLOAT:       "MONTHS_FLOAT",
	DAYS_INT:           "DAYS_INT",
	DAYS_FLOAT:         "DAYS_FLOAT",
	HOURS_INT:          "HOURS_INT",
	HOURS_FLOAT:        "HOURS_FLOAT",
	MINUTES_INT:        "MINUTES_INT",
	MINUTES_FLOAT:      "MINUTES_FLOAT",
	SECONDS_INT:        "SECONDS_INT",
	SECONDS_FLOAT:      "SECONDS_FLOAT",
	MILLISECONDS_INT:   "MILLISECONDS_INT",
	MILLISECONDS_FLOAT: "MILLISECONDS_FLOAT",
	MICROSECONDS_INT:   "MICROSECONDS_INT",
	MICROSECONDS_FLOAT: "MICROSECONDS_FLOAT",
	NANOSECONDS_INT:    "NANOSECONDS_INT",
	NANOSECONDS_FLOAT:  "NANOSECONDS_FLOAT",
}
