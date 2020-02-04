package common

import "time"

const (
	LayoutYYYYdMMdDD = "2006-01-02" // time parse layout for YYYY-MM-DD
	LayoutYYYYMMDD   = "20060102"   // time parse layout for YYYYMMDD
)

func ParseYYYYdMMdDD(str string) (parsed time.Time, err error) {
	return time.Parse(LayoutYYYYdMMdDD, str)
}

func ParseYYYYMMDD(str string) (parsed time.Time, err error) {
	return time.Parse(LayoutYYYYMMDD, str)
}
