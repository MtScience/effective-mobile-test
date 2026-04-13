package types

import (
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01"`, string(b))
	if err != nil {
		return err
	}
	d.Time = date
	return
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Format(`"2006-01"`)), nil
}
