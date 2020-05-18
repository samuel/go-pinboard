package pinboard

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type UpdateResponse struct {
	Time time.Time `xml:"time,attr" json:"time"`
}

type DatesResponse struct {
	Tag   string       `xml:"tag,attr" json:"tag"`
	User  string       `xml:"user,attr" json:"user"`
	Dates []*DateCount `xml:"date" json:"date"`
}

type Date struct {
	Year  int
	Month int
	Day   int
}

type DateCount struct {
	Count int  `xml:"count,attr" json:"count"`
	Date  Date `xml:"date,attr" json:"date"`
}

type PostsResponse struct {
	Date  time.Time `xml:"dt,atttr" json:"date"`
	Tag   string    `xml:"tag,attr" json:"tag"`
	User  string    `xml:"user,attr" json:"user"`
	Posts []*Post   `xml:"post" json:"posts"`
}

type Bool bool

type Tags []string

type Post struct {
	Href        string    `xml:"href,attr" json:"href"`
	Description string    `xml:"description,attr" json:"description"`
	Extended    string    `xml:"extended,attr" json:"extended"` // HTML (sometimes escaped) description, generally <blockquote>
	Hash        string    `xml:"hash,attr" json:"hash"`
	Tags        Tags      `xml:"tag,attr" json:"tag"`
	Time        time.Time `xml:"time,attr" json:"time"`
	Meta        string    `xml:"meta,attr" json:"meta"`
	Shared      Bool      `xml:"shared,attr" json:"shared"`
	ToRead      Bool      `xml:"toread,attr" json:"toread"`
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
