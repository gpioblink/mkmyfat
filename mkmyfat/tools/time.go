package tools

import "time"

func GetDateTimeForFAT(dateTime time.Time) (date uint16, time uint16, timeTenth uint16) {
	t := dateTime
	date = uint16(t.Year())<<9 | uint16(t.Month())<<5 | uint16(t.Day())
	time = uint16(t.Hour())<<11 | uint16(t.Minute())<<5 | uint16(t.Second()/2)
	timeTenth = uint16(t.Nanosecond() / 10000000)
	return date, time, timeTenth
}
