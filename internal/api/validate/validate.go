package validate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Options represents required options used by the Validate API Client.
type Options struct {
	Token string
}

// Client is a client for the Validate API.
type Client struct {
	opts   *Options
	client http.Client
}

// NewClient returns a new Client after verifying opts.
func NewClient(opts *Options) (*Client, error) {
	if opts.Token == "" {
		return nil, fmt.Errorf("invalid options: missing Token")
	}
	return &Client{opts: opts}, nil
}

type response struct {
	Status int `json:"status"`
	// Note: There may be additional (ignored) fields in the response.
}

// Status represents the validation status of the provided user/group key.
type Status int

const (
	// StatusUnknown indicates that validity of the provided user/group key is
	// unknown.
	StatusUnknown Status = iota
	// StatusOK indicates that the provided user/group key is valid.
	StatusOK Status = iota
	// StatusInvalid indicates that the provided user/group key is invalid.
	StatusInvalid Status = iota
)

// String returns a string representation of the Status.
func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusInvalid:
		return "INVALID"
	default:
		return "UNKNOWN"
	}
}

// Check uses the Validate API verify the user/group key provided, returning
// using the token associated with this Client.
func (c *Client) Check(ctx context.Context, user string) (Status, error) {
	vreq := url.Values{}
	vreq.Set("token", c.opts.Token)
	vreq.Set("user", user)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.pushover.net/1/users/validate.json", bytes.NewBufferString(vreq.Encode()))
	if err != nil {
		return StatusUnknown, fmt.Errorf("failed to initialize validate API request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return StatusUnknown, fmt.Errorf("failed to send validate API request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return StatusUnknown, fmt.Errorf("failed to read validate API response body: %v", err)
	}
	var vresp response
	if err := json.Unmarshal(body, &vresp); err != nil {
		return StatusUnknown, fmt.Errorf("failed to unmarshal validate API response body: %v", err)
	}
	if vresp.Status == 1 {
		return StatusOK, nil
	}
	return StatusInvalid, nil
}
