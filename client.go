package logcache

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/golang/protobuf/jsonpb"
)

// Client reads from LogCache via the RESTful API.
type Client struct {
	addr string
}

// NewIngressClient creates a Client.
func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

// Read queries the LogCache and returns the given envelopes. To override any
// query defaults (e.g., end time), use the according option.
func (c *Client) Read(sourceID string, start time.Time, opts ...ReadOption) ([]*loggregator_v2.Envelope, error) {
	u, err := url.Parse(c.addr)
	if err != nil {
		return nil, err
	}
	u.Path = "v1/read/" + sourceID
	q := u.Query()
	q.Set("start_time", strconv.FormatInt(start.UnixNano(), 10))

	// allow the given options to configure the URL.
	for _, o := range opts {
		o(u, q)
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var r logcache.ReadResponse
	if err := jsonpb.Unmarshal(resp.Body, &r); err != nil {
		return nil, err
	}

	return r.Envelopes.Batch, nil
}

// ReadOption configures the URL that is used to submit the query. The
// RawQuery is set to the decoded query parameters after each option is
// invoked.
type ReadOption func(u *url.URL, q url.Values)

// WithEndTime sets the 'end_time' query parameter to the given time. It
// defaults to empty, and therefore the end of the cache.
func WithEndTime(t time.Time) ReadOption {
	return func(u *url.URL, q url.Values) {
		q.Set("end_time", strconv.FormatInt(t.UnixNano(), 10))
	}
}

// WithLimit sets the 'limit' query parameter to the given value. It
// defaults to empty, and therefore 100 envelopes.
func WithLimit(limit int) ReadOption {
	return func(u *url.URL, q url.Values) {
		q.Set("limit", strconv.Itoa(limit))
	}
}

// WithEnvelopeType sets the 'envelope_type' query parameter to the given value. It
// defaults to empty, and therefore any envelope type.
func WithEnvelopeType(t logcache.EnvelopeTypes) ReadOption {
	return func(u *url.URL, q url.Values) {
		q.Set("envelope_type", t.String())
	}
}
