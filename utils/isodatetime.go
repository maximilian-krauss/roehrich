package utils

import "time"

type IsoDateTime struct {
	time.Time
}

func (t *IsoDateTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05Z0700"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}
