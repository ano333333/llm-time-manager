package utils

import "time"

// 日本標準時のタイムゾーンを返す
func GetJSTTimezone() *time.Location {
	return time.FixedZone("JST", 9*60*60)
}
