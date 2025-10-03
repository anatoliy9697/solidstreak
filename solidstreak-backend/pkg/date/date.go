package date

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
	"time"
)

type Date time.Time

func New(t time.Time) Date {
	return Date(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
}

func (d Date) String() string {
	return time.Time(d).Format(time.DateOnly)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format(time.DateOnly) + `"`), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return time.Time(d), nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date(time.Time{})
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*d = New(v)
		return nil
	case []byte:
		t, err := time.Parse(time.DateOnly, string(v))
		if err != nil {
			return err
		}
		*d = New(t)
		return nil
	case string:
		t, err := time.Parse(time.DateOnly, v)
		if err != nil {
			return err
		}
		*d = New(t)
		return nil
	default:
		return errors.New("cannot scan type into Date: " + reflect.TypeOf(value).String())
	}
}
