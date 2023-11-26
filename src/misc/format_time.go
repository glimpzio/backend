package misc

import "time"

func FormatTime(timestamp time.Time) string {
	return timestamp.Format("2006-01-02 15:04:05")
}
