package pinboard

// https://pinboard.in/api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL   = "https://api.pinboard.in/v1/"
	userAgent = "go-pinboard/0.1"
)

type StatusCodeError int

func (e StatusCodeError) Error() string {
	return fmt.Sprintf("pinboard: bad status code %d", int(e))
}

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

// PostsUpdated returns the most recent time a bookmark was added, updated or deleted.
// Use this before calling posts/all to see if the data has changed since the last fetch.
func (c *Client) PostsUpdated(ctx context.Context) (time.Time, error) {
	var res UpdateResponse
	if err := c.get(ctx, "posts/update", nil, &res); err != nil {
		return time.Time{}, err
	}
	return res.Time, nil
}

// PostDates returns a list of dates with the number of posts at each date.
func (c *Client) PostDates(ctx context.Context, tag string) (map[Date]int, error) {
	params := url.Values{}
	if tag != "" {
		params.Set("tag", tag)
	}
	var res DatesResponse
	if err := c.get(ctx, "posts/dates", params, &res); err != nil {
		return nil, err
	}
	dates := make(map[Date]int, len(res.Dates))
	for d, n := range res.Dates {
		c, err := n.Int64()
		if err != nil {
			return nil, err
		}
		dates[d] = int(c)
	}
	return dates, nil
}

// RecentPosts returns a list of the user's most recent posts, filtered by tag.
func (c *Client) RecentPosts(ctx context.Context, tag string, count int) ([]*Post, error) {
	params := url.Values{}
	if tag != "" {
		params.Set("tag", tag)
	}
	if count != 0 {
		params.Set("count", strconv.Itoa(count))
	}
	var res PostsResponse
	if err := c.get(ctx, "posts/recent", params, &res); err != nil {
		return nil, err
	}
	return res.Posts, nil
}

// AllPosts returnss all bookmarks in the user's account.
func (c *Client) AllPosts(ctx context.Context, tag string, start, results int, fromDate, toDate time.Time, meta bool) ([]*Post, error) {
	params := url.Values{}
	if tag != "" {
		params.Set("tag", tag)
	}
	if start != 0 {
		params.Set("start", strconv.Itoa(start))
	}
	if results != 0 {
		params.Set("results", strconv.Itoa(results))
	}
	if !fromDate.IsZero() {
		params.Set("fromdt", fromDate.Format(time.RFC3339))
	}
	if !toDate.IsZero() {
		params.Set("todt", toDate.Format(time.RFC3339))
	}
	if meta {
		params.Set("meta", "1")
	}
	var posts []*Post
	if err := c.get(ctx, "posts/all", params, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// DeletePost deletes a bookmark.
func (c *Client) DeletePost(ctx context.Context, postURL string) error {
	return c.get(ctx, "posts/delete", url.Values{"url": []string{postURL}}, nil)
}

func (c *Client) get(ctx context.Context, path string, params url.Values, response interface{}) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("auth_token", c.token)
	params.Set("format", "json")
	u := baseURL + path + "?" + params.Encode()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", userAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return StatusCodeError(res.StatusCode)
	}
	if response == nil {
		return nil
	}
	dec := json.NewDecoder(res.Body)
	return dec.Decode(response)
}
