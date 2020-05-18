package pinboard

import (
	"context"
	"os"
	"testing"
	"time"
)

func testClient(t *testing.T) *Client {
	token := os.Getenv("PINBOARD_TOKEN")
	if token == "" {
		t.Skip("Missing PINBOARD_TOKEN env")
	}
	return NewClient(token)
}

func TestPostsUpdated(t *testing.T) {
	c := testClient(t)
	tm, err := c.PostsUpdated()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tm.String())
}

func TestPostDates(t *testing.T) {
	c := testClient(t)
	dates, err := c.PostDates(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	for i, d := range dates {
		t.Logf("%+v", d)
		if i > 2 {
			break
		}
	}
}

func TestRecentPosts(t *testing.T) {
	c := testClient(t)
	posts, err := c.RecentPosts(context.Background(), "", 4)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range posts {
		t.Logf("%+v", p)
	}
}

func TestAllPosts(t *testing.T) {
	c := testClient(t)
	posts, err := c.AllPosts(context.Background(), "", 0, 4, time.Time{}, time.Time{}, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range posts {
		t.Logf("%+v", p)
	}
}
