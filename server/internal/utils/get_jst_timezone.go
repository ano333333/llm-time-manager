package utils

import "time"

func GetJSTTimezone() *time.Location {
	return time.FixedZone("JST", 9*60*60)
}
