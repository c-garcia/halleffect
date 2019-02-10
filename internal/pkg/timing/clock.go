package timing

//go:generate mockgen -source=clock.go -destination=mocks/clock.go -package=mocks

import "time"

type Clock interface {
	UnixTime() int64
}

type SystemClock struct{}

func NewSystemClock() *SystemClock {
	return nil
}

func (c *SystemClock) UnixTime() int64 {
	return time.Now().Unix()
}
