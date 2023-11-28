package timescanner

type Token uint8

// Name of the token.
func (t Token) String() string {
	if int(t) > len(tokenNames) {
		return "UNKNOWN"
	}

	return tokenNames[t]
}

const (
	ZERO_VALUE                             Token = iota // Zero value for Type
	INVALID_FORMAT_DIRECTIVE                            // Invalid format directive
	END_OF_FILE                                         // End Of File has been reached
	PERCENT                                             // Literal percent
	NEWLINE                                             // Literal newline character
	TAB                                                 // Literal tab character
	TEXT                                                // Literal text
	FULL_YEAR_WEEK_BASED                                // The ISO 8601 week-based year with century as a decimal number. The year corresponding to the ISO week number (see %V). This has the same format and value as %-Y, except that if the ISO week number belongs to the previous or next year, that year is used instead.
	FULL_YEAR_WEEK_BASED_SPACE_PADDED                   // (space-padded) The ISO 8601 week-based year with century as a decimal number. The year corresponding to the ISO week number (see %V). This has the same format and value as %_Y, except that if the ISO week number belongs to the previous or next year, that year is used instead.
	FULL_YEAR_WEEK_BASED_ZERO_PADDED                    // (zero-padded) The ISO 8601 week-based year with century as a decimal number. The 4-digit year corresponding to the ISO week number (see %V). This has the same format and value as %Y, except that if the ISO week number belongs to the previous or next year, that year is used instead.
	FULL_YEAR                                           // Year with century (can be negative)
	FULL_YEAR_ZERO_PADDED                               // Year with century (can be negative, 4 digits at least, zero padded)
	FULL_YEAR_SPACE_PADDED                              // Year with century (can be negative, space padded)
	CENTURY                                             // year / 100 (round down.  20 in 2009)
	CENTURY_SPACE_PADDED                                // year / 100 (round down.  20 in 2009), space-padded
	CENTURY_ZERO_PADDED                                 // year / 100 (round down.  20 in 2009), zero-padded
	YEAR_LAST_TWO_WEEK_BASED                            // The ISO 8601 week-based year % 100 (0..99)
	YEAR_LAST_TWO_WEEK_BASED_ZERO_PADDED                // The ISO 8601 week-based year % 100 (00..99), zero padded
	YEAR_LAST_TWO_WEEK_BASED_SPACE_PADDED               // The ISO 8601 week-based year % 100 ( 0..99), space padded
	YEAR_LAST_TWO                                       // year % 100 (0..99)
	YEAR_LAST_TWO_ZERO_PADDED                           // year % 100 (00..99), zero padded
	YEAR_LAST_TWO_SPACE_PADDED                          // year % 100 ( 0..99), space padded
	MONTH                                               // Month of the year (1..12)
	MONTH_ZERO_PADDED                                   // Month of the year, zero-padded (01..12)
	MONTH_SPACE_PADDED                                  // Month of the year, space-padded ( 1..12)
	MONTH_FULL_NAME                                     // The full month name "January"
	MONTH_FULL_NAME_UPPERCASE                           // The uppercase full month name "JANUARY"
	MONTH_ABBREVIATED_NAME                              // The abbreviated month name "Jan"
	MONTH_ABBREVIATED_NAME_UPPERCASE                    // The uppercase abbreviated month name "JAN"
	DAY_OF_MONTH                                        // Day of the month (1..31)
	DAY_OF_MONTH_SPACE_PADDED                           // Day of the month, space-padded ( 1..31)
	DAY_OF_MONTH_ZERO_PADDED                            // Day of the month, zero-padded (01..31)
	DAY_OF_YEAR                                         // Day of the year (1..366)
	DAY_OF_YEAR_SPACE_PADDED                            // Day of the year, space-padded (  1..366)
	DAY_OF_YEAR_ZERO_PADDED                             // Day of the year, zero-padded (001..366)
	HOUR_OF_DAY                                         // Hour of the day, 24-hour clock (0..23)
	HOUR_OF_DAY_SPACE_PADDED                            // Hour of the day, 24-hour clock, space-padded ( 0..23)
	HOUR_OF_DAY_ZERO_PADDED                             // Hour of the day, 24-hour clock, zero-padded (00..23)
	HOUR_OF_DAY12                                       // Hour of the day, 12-hour clock
	HOUR_OF_DAY12_SPACE_PADDED                          // Hour of the day, 12-hour clock, space-padded ( 1..12)
	HOUR_OF_DAY12_ZERO_PADDED                           // Hour of the day, 12-hour clock, zero-padded (01..12)
	MERIDIEM_INDICATOR_LOWERCASE                        // Meridiem indicator, lowercase (`am` or `pm`)
	MERIDIEM_INDICATOR_UPPERCASE                        // Meridiem indicator, uppercase (`AM` or `PM`)
	MINUTE_OF_HOUR                                      // Minute of the hour (0..59)
	MINUTE_OF_HOUR_SPACE_PADDED                         // Minute of the hour, space-padded ( 0..59)
	MINUTE_OF_HOUR_ZERO_PADDED                          // Minute of the hour, zero-padded (00..59)
	SECOND_OF_MINUTE                                    // Second of the minute (00..60)
	SECOND_OF_MINUTE_SPACE_PADDED                       // Second of the minute, space-padded ( 0..60)
	SECOND_OF_MINUTE_ZERO_PADDED                        // Second of the minute, zero-padded (00..60)
	MILLISECOND_OF_SECOND                               // Millisecond of the second (0..999)
	MILLISECOND_OF_SECOND_SPACE_PADDED                  // Millisecond of the second, space-padded (  0..999)
	MILLISECOND_OF_SECOND_ZERO_PADDED                   // Millisecond of the second, zero-padded (000..999)
	TIMEZONE_NAME                                       // Timezone name
	TIMEZONE_OFFSET                                     // Time zone as hour and minute offset from UTC (e.g. +0900)
	TIMEZONE_OFFSET_COLON                               // hour and minute offset from UTC with a colon (e.g. +09:00)
	DAY_OF_WEEK_FULL_NAME                               // The full weekday name "Sunday"
	DAY_OF_WEEK_FULL_NAME_UPPERCASE                     // The full weekday name "SUNDAY"
	DAY_OF_WEEK_ABBREVIATED_NAME                        // The abbreviated name "Sun"
	DAY_OF_WEEK_ABBREVIATED_NAME_UPPERCASE              // The abbreviated name "SUN"
	DAY_OF_WEEK_NUMBER                                  // The day of the week as a decimal, range 1 to 7, Monday being 1.
	DAY_OF_WEEK_NUMBER_ALT                              // The day of the week as a decimal, range 0 to 6, Sunday being 0
	UNIX_SECONDS                                        // Number of seconds since 1970-01-01 00:00:00 UTC.
	UNIX_MILLISECONDS                                   // Number of milliseconds since 1970-01-01 00:00:00 UTC.
	UNIX_MICROSECONDS                                   // Number of microseconds since 1970-01-01 00:00:00 UTC.
	UNIX_NANOSECONDS                                    // Number of nanoseconds since 1970-01-01 00:00:00 UTC.
	UNIX_PICOSECONDS                                    // Number of picoseconds since 1970-01-01 00:00:00 UTC.
	UNIX_FEMTOSECONDS                                   // Number of femtoseconds since 1970-01-01 00:00:00 UTC.
	UNIX_ATTOSECONDS                                    // Number of attoseconds since 1970-01-01 00:00:00 UTC.
	UNIX_ZEPTOSECONDS                                   // Number of zeptoseconds since 1970-01-01 00:00:00 UTC.
	UNIX_YOCTOSECONDS                                   // Number of yoctoseconds since 1970-01-01 00:00:00 UTC.
	WEEK_OF_WEEK_BASED_YEAR                             // Week number of the week-based year (1..53)
	WEEK_OF_WEEK_BASED_YEAR_SPACE_PADDED                // Week number of the week-based year, space-padded ( 1..53)
	WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED                 // Week number of the week-based year, zero-padded (01..53)
	WEEK_OF_YEAR                                        // Week number of the year. The week starts with Monday. (0..53) The week number of the current year as a decimal number, range 00 to 53, starting with the first Monday as the first day of week 01.
	WEEK_OF_YEAR_SPACE_PADDED                           // Week number of the year. The week starts with Monday. Space-padded ( 0..53). The week number of the current year as a decimal number, range 00 to 53, starting with the first Monday as the first day of week 01.
	WEEK_OF_YEAR_ZERO_PADDED                            // Week number of the year. The week starts with Monday. Zero-padded (00..53). The week number of the current year as a decimal number, range 00 to 53, starting with the first Monday as the first day of week 01.
	WEEK_OF_YEAR_ALT                                    // Week number of the year. The week starts with Sunday. (0..53)
	WEEK_OF_YEAR_ALT_SPACE_PADDED                       // Week number of the year. The week starts with Sunday. Space-padded ( 0..53)
	WEEK_OF_YEAR_ALT_ZERO_PADDED                        // Week number of the year. The week starts with Sunday. Zero-padded (00..53)
	DATE_AND_TIME                                       // date and time (%a %b %e %T %Y)
	DATE_AND_TIME_UPPERCASE                             // date and time (%a %b %e %T %Y)
	DATE                                                // Date (%m/%d/%y)
	ISO8601_DATE                                        // Equivalent to %Y-%m-%d (the ISO 8601 date format).
	TIME12                                              // 12-hour time (%I:%M:%S %p)
	TIME24                                              // 24-hour time (%H:%M)
	TIME24_SECONDS                                      // 24-hour time (%H:%M:%S)
	DATE1_FORMAT                                        // The date and time in date(1) format (%a %b %e %H:%M:%S %Z %Y)
	DATE1_FORMAT_UPPERCASE                              // The date and time in date(1) format (%a %b %e %H:%M:%S %Z %Y)
	MICROSECOND_OF_SECOND                               // Fractional seconds digits up to 6 digits (microsecond)
	MICROSECOND_OF_SECOND_SPACE_PADDED                  // Fractional seconds digits, up to 6 digits, space-padded (microsecond)
	MICROSECOND_OF_SECOND_ZERO_PADDED                   // Fractional seconds digits, 6 digits, zero-padded (microsecond)
	NANOSECOND_OF_SECOND                                // Fractional seconds digits up to 9 digits (nanosecond)
	NANOSECOND_OF_SECOND_SPACE_PADDED                   // Fractional seconds digits, up to 9 digits, space-padded (nanosecond)
	NANOSECOND_OF_SECOND_ZERO_PADDED                    // Fractional seconds digits, 9 digits, zero-padded (nanosecond)
	PICOSECOND_OF_SECOND                                // Fractional seconds digits up to 12 digits (picosecond)
	PICOSECOND_OF_SECOND_SPACE_PADDED                   // Fractional seconds digits, up to 12 digits, space-padded (picosecond)
	PICOSECOND_OF_SECOND_ZERO_PADDED                    // Fractional seconds digits, 12 digits, zero-padded (picosecond)
	FEMTOSECOND_OF_SECOND                               // Fractional seconds digits up to 15 digits (femtosecond)
	FEMTOSECOND_OF_SECOND_SPACE_PADDED                  // Fractional seconds digits, up to 15 digits, space-padded (femtosecond)
	FEMTOSECOND_OF_SECOND_ZERO_PADDED                   // Fractional seconds digits, 15 digits, zero-padded (femtosecond)
	ATTOSECOND_OF_SECOND                                // Fractional seconds digits up to 18 digits (attosecond)
	ATTOSECOND_OF_SECOND_SPACE_PADDED                   // Fractional seconds digits, up to 18 digits, space-padded (attosecond)
	ATTOSECOND_OF_SECOND_ZERO_PADDED                    // Fractional seconds digits, 18 digits, zero-padded (attosecond)
	ZEPTOSECOND_OF_SECOND                               // Fractional seconds digits up to 21 digits (zeptosecond)
	ZEPTOSECOND_OF_SECOND_SPACE_PADDED                  // Fractional seconds digits, up to 21 digits, space-padded (zeptosecond)
	ZEPTOSECOND_OF_SECOND_ZERO_PADDED                   // Fractional seconds digits, 21 digits, zero-padded (zeptosecond)
	YOCTOSECOND_OF_SECOND                               // Fractional seconds digits up to 24 digits (yoctosecond)
	YOCTOSECOND_OF_SECOND_SPACE_PADDED                  // Fractional seconds digits, up to 24 digits, space-padded (yoctosecond)
	YOCTOSECOND_OF_SECOND_ZERO_PADDED                   // Fractional seconds digits, 24 digits, zero-padded (yoctosecond)
)

var tokenNames = [...]string{
	INVALID_FORMAT_DIRECTIVE:               "INVALID_FORMAT_DIRECTIVE",
	END_OF_FILE:                            "END_OF_FILE",
	PERCENT:                                "%%",
	NEWLINE:                                "%n",
	TAB:                                    "%t",
	TEXT:                                   "TEXT",
	FULL_YEAR_WEEK_BASED:                   "%-G",
	FULL_YEAR_WEEK_BASED_SPACE_PADDED:      "%_G",
	FULL_YEAR_WEEK_BASED_ZERO_PADDED:       "%G",
	FULL_YEAR:                              "%-Y",
	FULL_YEAR_SPACE_PADDED:                 "%_Y",
	FULL_YEAR_ZERO_PADDED:                  "%Y",
	CENTURY:                                "%-C",
	CENTURY_SPACE_PADDED:                   "%_C",
	CENTURY_ZERO_PADDED:                    "%C",
	YEAR_LAST_TWO_WEEK_BASED:               "%-g",
	YEAR_LAST_TWO_WEEK_BASED_SPACE_PADDED:  "%_g",
	YEAR_LAST_TWO_WEEK_BASED_ZERO_PADDED:   "%g",
	YEAR_LAST_TWO:                          "%-y",
	YEAR_LAST_TWO_SPACE_PADDED:             "%_y",
	YEAR_LAST_TWO_ZERO_PADDED:              "%y",
	MONTH:                                  "%-m",
	MONTH_SPACE_PADDED:                     "%_m",
	MONTH_ZERO_PADDED:                      "%m",
	MONTH_FULL_NAME:                        "%B",
	MONTH_FULL_NAME_UPPERCASE:              "%^B",
	MONTH_ABBREVIATED_NAME:                 "%b",
	MONTH_ABBREVIATED_NAME_UPPERCASE:       "%^b",
	DAY_OF_MONTH:                           "%-d",
	DAY_OF_MONTH_SPACE_PADDED:              "%_d",
	DAY_OF_MONTH_ZERO_PADDED:               "%d",
	DAY_OF_YEAR:                            "%-j",
	DAY_OF_YEAR_SPACE_PADDED:               "%_j",
	DAY_OF_YEAR_ZERO_PADDED:                "%j",
	HOUR_OF_DAY:                            "%-H",
	HOUR_OF_DAY_SPACE_PADDED:               "%_H",
	HOUR_OF_DAY_ZERO_PADDED:                "%H",
	HOUR_OF_DAY12:                          "%-I",
	HOUR_OF_DAY12_SPACE_PADDED:             "%_I",
	HOUR_OF_DAY12_ZERO_PADDED:              "%I",
	MERIDIEM_INDICATOR_LOWERCASE:           "%P",
	MERIDIEM_INDICATOR_UPPERCASE:           "%p",
	MINUTE_OF_HOUR:                         "%-M",
	MINUTE_OF_HOUR_SPACE_PADDED:            "%_M",
	MINUTE_OF_HOUR_ZERO_PADDED:             "%M",
	SECOND_OF_MINUTE:                       "%-S",
	SECOND_OF_MINUTE_SPACE_PADDED:          "%_S",
	SECOND_OF_MINUTE_ZERO_PADDED:           "%S",
	MILLISECOND_OF_SECOND:                  "%-L",
	MILLISECOND_OF_SECOND_SPACE_PADDED:     "%_L",
	MILLISECOND_OF_SECOND_ZERO_PADDED:      "%L",
	TIMEZONE_NAME:                          "%Z",
	TIMEZONE_OFFSET:                        "%z",
	TIMEZONE_OFFSET_COLON:                  "%:z",
	DAY_OF_WEEK_FULL_NAME:                  "%A",
	DAY_OF_WEEK_FULL_NAME_UPPERCASE:        "%^A",
	DAY_OF_WEEK_ABBREVIATED_NAME:           "%a",
	DAY_OF_WEEK_ABBREVIATED_NAME_UPPERCASE: "%^a",
	DAY_OF_WEEK_NUMBER:                     "%u",
	DAY_OF_WEEK_NUMBER_ALT:                 "%w",
	UNIX_SECONDS:                           "%s",
	UNIX_MILLISECONDS:                      "%Q",
	UNIX_MICROSECONDS:                      "%6s",
	UNIX_NANOSECONDS:                       "%9s",
	UNIX_PICOSECONDS:                       "%12s",
	UNIX_FEMTOSECONDS:                      "%15s",
	UNIX_ATTOSECONDS:                       "%18s",
	UNIX_ZEPTOSECONDS:                      "%21s",
	UNIX_YOCTOSECONDS:                      "%24s",
	WEEK_OF_WEEK_BASED_YEAR:                "%-V",
	WEEK_OF_WEEK_BASED_YEAR_SPACE_PADDED:   "%_V",
	WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED:    "%V",
	WEEK_OF_YEAR:                           "%-W",
	WEEK_OF_YEAR_SPACE_PADDED:              "%_W",
	WEEK_OF_YEAR_ZERO_PADDED:               "%W",
	WEEK_OF_YEAR_ALT:                       "%-U",
	WEEK_OF_YEAR_ALT_SPACE_PADDED:          "%_U",
	WEEK_OF_YEAR_ALT_ZERO_PADDED:           "%U",
	DATE_AND_TIME:                          "%c",
	DATE_AND_TIME_UPPERCASE:                "%^c",
	DATE:                                   "%D",
	ISO8601_DATE:                           "%F",
	TIME12:                                 "%r",
	TIME24:                                 "%R",
	TIME24_SECONDS:                         "%T",
	DATE1_FORMAT:                           "%+",
	DATE1_FORMAT_UPPERCASE:                 "%^+",
	MICROSECOND_OF_SECOND:                  "%-6N",
	MICROSECOND_OF_SECOND_SPACE_PADDED:     "%_6N",
	MICROSECOND_OF_SECOND_ZERO_PADDED:      "%6N",
	NANOSECOND_OF_SECOND:                   "%-N",
	NANOSECOND_OF_SECOND_SPACE_PADDED:      "%_N",
	NANOSECOND_OF_SECOND_ZERO_PADDED:       "%N",
	PICOSECOND_OF_SECOND:                   "%-12N",
	PICOSECOND_OF_SECOND_SPACE_PADDED:      "%_12N",
	PICOSECOND_OF_SECOND_ZERO_PADDED:       "%12N",
	FEMTOSECOND_OF_SECOND:                  "%-15N",
	FEMTOSECOND_OF_SECOND_SPACE_PADDED:     "%_15N",
	FEMTOSECOND_OF_SECOND_ZERO_PADDED:      "%15N",
	ATTOSECOND_OF_SECOND:                   "%-18N",
	ATTOSECOND_OF_SECOND_SPACE_PADDED:      "%_18N",
	ATTOSECOND_OF_SECOND_ZERO_PADDED:       "%18N",
	ZEPTOSECOND_OF_SECOND:                  "%-21N",
	ZEPTOSECOND_OF_SECOND_SPACE_PADDED:     "%_21N",
	ZEPTOSECOND_OF_SECOND_ZERO_PADDED:      "%21N",
	YOCTOSECOND_OF_SECOND:                  "%-24N",
	YOCTOSECOND_OF_SECOND_SPACE_PADDED:     "%_24N",
	YOCTOSECOND_OF_SECOND_ZERO_PADDED:      "%24N",
}
