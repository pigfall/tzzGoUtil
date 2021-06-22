package time

import (
	"time"
	stdtime "time"
)

const format_yyyy_mm_dd_hh_mm_ss = "2006-01-02 15:04:05"

func FormatToYYYY_MM_DD_HH_MM_SS(t stdtime.Time) string {
	return t.Format(format_yyyy_mm_dd_hh_mm_ss)
}

func ParseFromYYYY_MM_DD_HH_MM_SS(timeStr string) (time.Time, error) {
	return time.Parse(format_yyyy_mm_dd_hh_mm_ss, timeStr)
}
