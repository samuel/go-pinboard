package pinboard

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestXMLPosts(t *testing.T) {
	xmlInput := []byte(`<posts tag="tag" user="user">
        <post href="http://www.weather.com/" description="weather.com"
        hash="6cfedbe75f413c56b6ce79e6fa102aba" tag="weather reference"
        time="2005-11-29T20:30:47Z" />
        <post href="http://www.nytimes.com/"
        description="The New York Times - Breaking News, World News &amp; Multimedia"
        extended="requires login" hash="ca1e6357399774951eed4628d69eb84b"
        tag="news media" time="2005-11-29T20:30:05Z" />
    </posts>`)

	var posts PostsResponse
	if err := xml.Unmarshal(xmlInput, &posts); err != nil {
		t.Fatal(err)
	}
	if posts.Tag != "tag" {
		t.Fatalf("Expected 'tag', got '%s'", posts.Tag)
	}
	if posts.User != "user" {
		t.Fatalf("Expected 'user', got '%s'", posts.User)
	}
	if len(posts.Posts) != 2 {
		t.Fatalf("Expected 2, got %d", len(posts.Posts))
	}
	p := posts.Posts[0]
	if p.Time != time.Date(2005, 11, 29, 20, 30, 47, 0, time.UTC) {
		t.Fatalf("Wrong time: %s", p.Time.String())
	}
	if p.Href != "http://www.weather.com/" {
		t.Fatalf("Wrong href: %s", p.Href)
	}
	if p.Description != "weather.com" {
		t.Fatalf("Wrong description: %s", p.Description)
	}
}
