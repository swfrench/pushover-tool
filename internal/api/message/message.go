package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

type messageRequest struct {
	// Required fields:
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`
	// Optional fields:
	// TODO: Add more.
	Title *string `json:"title"`
}

type messageResponse struct {
	Status  int      `json:"status"`
	Request string   `json:"request"`
	Errors  []string `json:"errors"`
	// Note: There may be additional (ignored) fields in the response.
}

const messageStatusOK = 1

// Send uses the Message API to send the provided message to the associated user
// token. If non-nil, the message will use the provided title.
func (c *Client) Send(user, message string, title *string) error {
	bs, err := json.Marshal(&messageRequest{
		Token:   c.opts.Token,
		User:    user,
		Message: message,
		Title:   title,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal message API request body: %v", err)
	}
	buff := bytes.NewBuffer(bs)
	resp, err := c.client.Post("https://api.pushover.net/1/messages.json", "application/json", buff)
	if err != nil {
		return fmt.Errorf("failed to send message API request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read message API response body: %v", err)
	}
	mresp := new(messageResponse)
	if err := json.Unmarshal(body, mresp); err != nil {
		return fmt.Errorf("failed to unmarshal message API response bidy: %v", err)
	}
	if mresp.Status != messageStatusOK {
		return fmt.Errorf("message API returned non-OK status (%d) for request %q - errors: %s", mresp.Status, mresp.Request, strings.Join(mresp.Errors, ", "))
	}
	fmt.Println(mresp.Request)
	return nil
}