package gothinkutils

import "time"

type Datetime struct {
}

func (this Datetime) Timestamp() int64 {
	now := time.Now()
	return now.Unix()
}

func (this Datetime) TimestampMs() int64 {
	now := time.Now()
	return now.UnixMilli()
}
