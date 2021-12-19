package util

import "time"

func ToFormattedDate(t time.Time) string {
	return t.Format("02-Jan-2006")
}
