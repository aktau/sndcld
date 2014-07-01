package snd

import (
	"errors"
	"time"
)

type Time time.Time

const (
	SNDCLD_TIME_FORMAT = "2006/01/02 15:04:05 -0700"
)

func (j *Time) UnmarshalJSON(data []byte) (err error) {
	t, err := time.Parse(`"`+SNDCLD_TIME_FORMAT+`"`, string(data))
	*j = Time(t)
	return
}

func (j Time) MarshalJSON() ([]byte, error) {
	t := time.Time(j)
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(t.Format(`"` + SNDCLD_TIME_FORMAT + `"`)), nil
}
