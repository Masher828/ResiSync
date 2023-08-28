package shared_utils

import "time"

func NowInUTC() time.Time {
	return time.Now().UTC()
}
