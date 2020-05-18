package pinboard

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type UpdateResponse struct {
	Time time.Time `xml:"time,attr" json:"update_time"`
}

type DatesResponse struct {
	Tag   string               `json:"tag"`
	User  string               `json:"user"`
	Dates map[Date]json.Number `json:"dates"`
}

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)
}

type DateCount struct {
	Count int  `json:"count"`
	Date  Date `json:"date"`
}

type PostsResponse struct {
	Date  time.Time `json:"date"`
	Tag   string    `json:"tag"`
	User  string    `json:"user"`
	Posts []*Post   `son:"posts"`
}

type Bool bool

type Tags []string

type Post struct {
	Href        string    `json:"href"`
	Description string    `json:"description"`
	Extended    string    `json:"extended"` // HTML (sometimes escaped) description, generally <blockquote>
	Hash        string    `json:"hash"`
	Tags        Tags      `json:"tags"`
	Time        time.Time `json:"time"`
	Meta        string    `json:"meta"`
	Shared      Bool      `json:"shared"`
	ToRead      Bool      `json:"toread"`
}

func (b *Bool) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "yes", "y", "true", "t", "1":
		*b = true
		return nil
	case "no", "n", "false", "f", "0":
		*b = false
		return nil
	}
	return fmt.Errorf("invalid bool %q", text)
}

func (t *Tags) UnmarshalText(text []byte) error {
	*t = strings.Split(string(text), " ")
	return nil
}

func (d *Date) UnmarshalText(text []byte) error {
	p := strings.Split(string(text), "-")
	if len(p) != 3 {
		return fmt.Errorf("invalid date %q", text)
	}
	var err error
	d.Year, err = strconv.Atoi(p[0])
	if err != nil {
		return fmt.Errorf("invalid year %q in date %q: %w", p[0], text, err)
	}
	d.Month, err = strconv.Atoi(strings.TrimLeft(p[1], "0"))
	if err != nil {
		return fmt.Errorf("invalid month %q in date %q: %w", p[1], text, err)
	}
	d.Day, err = strconv.Atoi(strings.TrimLeft(p[2], "0"))
	if err != nil {
		return fmt.Errorf("invalid day %q in date %q: %w", p[2], text, err)
	}
	return nil
}
