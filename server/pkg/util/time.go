package util

import "time"

func CurrentTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func CurrentTimestamp() int {
	return int(time.Now().Unix())
}
