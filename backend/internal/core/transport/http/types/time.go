package http_types

import "time"

const DateOnlyLayout = "2006-01-02"

// DateOnly wraps time.Time and marshals to JSON as "2006-01-02".
type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format(DateOnlyLayout) + `"`), nil
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return &time.ParseError{Layout: DateOnlyLayout, Value: s, Message: "missing quotes"}
	}
	s = s[1 : len(s)-1]

	t, err := time.Parse(DateOnlyLayout, s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) Time() time.Time {
	return time.Time(d)
}
