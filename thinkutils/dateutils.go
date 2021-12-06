package thinkutils

import "time"

type datetime struct {
}

func (this datetime) Timestamp() int64 {
	now := time.Now()
	return now.Unix()
}

func (this datetime) TimestampMs() int64 {
	now := time.Now()
	return now.UnixMilli()
}
