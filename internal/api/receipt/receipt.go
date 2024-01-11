package receipt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Options represents required options used by the Receipt API Client.
type Options struct {
	Token    string
	Interval time.Duration
}

// Client is a client for the Message API.
type Client struct {
	opts   *Options
	client http.Client
}

// minInterval is 2x the documented API limit (5s)
const minInterval = 10 * time.Second

// NewClient returns a new Client after verifying opts.
func NewClient(opts *Options) (*Client, error) {
	if opts.Token == "" {
		return nil, fmt.Errorf("invalid options: missing Token")
	}
	if opts.Interval < minInterval {
		return nil, fmt.Errorf("invalid options: interval (%v) is too small (minimum: %v)", opts.Interval, minInterval)
	}
	return &Client{opts: opts}, nil
}

type response struct {
	Status       int `json:"status"`
	Acknowledged int `json:"acknowledged"`
	// Note: There may be additional (ignored) fields in the response.
}

const (
	receiptStatusOK     = 1
	receiptAcknowledged = 1
)

var errNotAcknowledged = errors.New("not acknowledged")

// Wait uses the Receipt API to wait for acknowledgement of the provided
// receipt.
func (c *Client) Wait(receipt string) error {
	u := fmt.Sprintf("https://api.pushover.net/1/receipts/%s.json?%s", receipt, url.Values{"token": {c.opts.Token}}.Encode())
	check := func() error {
		log.Printf("Polling for ACK (receipt: %s)", receipt)
		resp, err := c.client.Get(u)
		if err != nil {
			return fmt.Errorf("failed to send receipt API request: %v", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read receipt API response body: %v", err)
		}
		var rresp response
		if err := json.Unmarshal(body, &rresp); err != nil {
			return fmt.Errorf("failed to unmarshal receipt API response body: %v", err)
		}
		if rresp.Status != receiptStatusOK {
			return fmt.Errorf("receipt API returned non-OK status (%d) for receipt %q", rresp.Status, receipt)
		}
		if rresp.Acknowledged != receiptAcknowledged {
			return errNotAcknowledged
		}
		return nil
	}
	tc := time.NewTicker(c.opts.Interval)
	defer tc.Stop()
	for {
		if err := check(); err == nil {
			return nil
		} else if !errors.Is(err, errNotAcknowledged) {
			return err
		}
		<-tc.C
	}
}
