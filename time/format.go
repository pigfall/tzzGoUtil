package time

import (
	"time"
)

const format_yyyy_mm_dd_hh_mm_ss = "2006-01-02 15:04:05"
const format_yyyy_mm_dd = "2006-01-02"
const format_hh_mm_ss = "15:04:05"

func FormatToYYYY_MM_DD_HH_MM_SS(t time.Time) string {
	return t.Format(format_yyyy_mm_dd_hh_mm_ss)
}

func FormatToYYYY_MM_DD(t time.Time) string {
	return t.Format(format_yyyy_mm_dd )
}

// => eg:  "2006-01-02 15:04:05"
func FormatToHH_MM_SS(t time.Time)string{
	return t.Format(format_hh_mm_ss)
}


func ParseFromYYYY_MM_DD_HH_MM_SS(timeStr string) (time.Time, error) {
	return time.Parse(format_yyyy_mm_dd_hh_mm_ss, timeStr)
}

func ParseFromYYYY_MM_DD(timeStr string) (time.Time, error) {
	return time.Parse(format_yyyy_mm_dd, timeStr)
}

