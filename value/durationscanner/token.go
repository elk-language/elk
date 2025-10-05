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
	ZERO_VALUE   Token = iota // Zero value for Type
	ERROR                     // Invalid format directive
	END_OF_FILE               // End Of File has been reached
	YEARS                     // Year component eg. "2Y"
	MONTHS                    // Month component eg. "15M"
	DAYS                      // Day component eg. "1.8D"
	HOURS                     // Hour component eg. "9h"
	MINUTES                   // Minute component eg. "10m"
	SECONDS                   // Second component eg. "50s"
	MILLISECONDS              // Millisecond component eg. "150ms"
	MICROSECONDS              // Microsecond component eg. "600us", "600µs"
	NANOSECONDS               // Nanosecond component eg. "750ns"
)

var tokenNames = [...]string{
	ERROR:        "ERROR",
	END_OF_FILE:  "END_OF_FILE",
	YEARS:        "YEARS",
	MONTHS:       "MONTHS",
	DAYS:         "DAYS",
	HOURS:        "HOURS",
	MINUTES:      "MINUTES",
	SECONDS:      "SECONDS",
	MILLISECONDS: "MILLISECONDS",
	MICROSECONDS: "MICROSECONDS",
	NANOSECONDS:  "NANOSECONDS",
}
