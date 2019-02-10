package doubles

type StoppedClock struct {
	at int64
}

func (s *StoppedClock) UnixTime() int64 {
	return s.at
}

func NewStoppedClock(at int64) *StoppedClock {
	return &StoppedClock{
		at: at,
	}
}
