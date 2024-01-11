package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Options represents required options used by the Message API Client.
type Options struct {
	Token string
}

// Client is a client for the Message API.
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
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Receipt string   `json:"receipt"`
	Errors  []string `json:"errors"`
	// Note: There may be additional (ignored) fields in the response.
}

const (
	messageStatusOK   = 1
	emergencyPriority = "2"
)

func stringSeconds(d time.Duration) string {
	return strconv.Itoa(int(d.Seconds()))
}

// Send uses the Message API to send the provided message to the specified user,
// using the token associated with this Client.
// If emergency is true, the message will be sent with emergency priority (2),
// with delivery of un-ACK'd messages retried with the specified retry period
// until expiration.
func (c *Client) Send(user, message, title string, emergency bool, retry, expire time.Duration) error {
	mreq := url.Values{}
	mreq.Set("token", c.opts.Token)
	mreq.Set("user", user)
	mreq.Set("message", message)
	mreq.Set("title", title)
	if emergency {
		mreq.Set("priority", emergencyPriority)
		mreq.Set("retry", stringSeconds(retry))
		mreq.Set("expire", stringSeconds(expire))
	}
	resp, err := c.client.Post("https://api.pushover.net/1/messages.json", "application/x-www-form-urlencoded", bytes.NewBufferString(mreq.Encode()))
	if err != nil {
		return fmt.Errorf("failed to send message API request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read message API response body: %v", err)
	}
	var mresp response
	if err := json.Unmarshal(body, &mresp); err != nil {
		return fmt.Errorf("failed to unmarshal message API response body: %v", err)
	}
	if mresp.Status != messageStatusOK {
		return fmt.Errorf("message API returned non-OK status (%d) for request %q - errors: %s", mresp.Status, mresp.Request, strings.Join(mresp.Errors, ", "))
	}
	if emergency {
		fmt.Println(mresp.Receipt)
	}
	return nil
}
