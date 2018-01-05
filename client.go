package logcache

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
)

// HTTPClient is an interface that represents a http.Client.
type HTTPClient interface {
	Get(string) (*http.Response, error)
}

// Client reads from LogCache via the RESTful API.
type Client struct {
	addr string

	httpClient HTTPClient
	grpcClient logcache.EgressClient
}

// NewIngressClient creates a Client.
func NewClient(addr string, opts ...ClientOption) *Client {
	c := &Client{
		addr: addr,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// ClientOption configures the LogCache client.
type ClientOption func(c *Client)

// WithHTTPClient sets the HTTP client. It defaults to a client that timesout
// after 5 seconds.
func WithHTTPClient(h HTTPClient) ClientOption {
	return func(c *Client) {
		c.httpClient = h
	}
}

// WithViaGRPC enables gRPC instead of HTTP/1 for reading from LogCache.
func WithViaGRPC(opts ...grpc.DialOption) ClientOption {
	return func(c *Client) {
		conn, err := grpc.Dial(c.addr, opts...)
		if err != nil {
			panic(fmt.Sprintf("failed to dial via gRPC: %s", err))
		}

		c.grpcClient = logcache.NewEgressClient(conn)
	}
}

// Read queries the LogCache and returns the given envelopes. To override any
// query defaults (e.g., end time), use the according option.
func (c *Client) Read(sourceID string, start time.Time, opts ...ReadOption) ([]*loggregator_v2.Envelope, error) {
	if c.grpcClient != nil {
		return c.grpcRead(sourceID, start, opts)
	}

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

	resp, err := c.httpClient.Get(u.String())
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

func (c *Client) grpcRead(sourceID string, start time.Time, opts []ReadOption) ([]*loggregator_v2.Envelope, error) {
	u := &url.URL{}
	q := u.Query()
	// allow the given options to configure the URL.
	for _, o := range opts {
		o(u, q)
	}

	req := &logcache.ReadRequest{
		SourceId:  sourceID,
		StartTime: start.UnixNano(),
	}

	if v, ok := q["limit"]; ok {
		req.Limit, _ = strconv.ParseInt(v[0], 10, 64)
	}

	if v, ok := q["end_time"]; ok {
		req.EndTime, _ = strconv.ParseInt(v[0], 10, 64)
	}

	if v, ok := q["envelope_type"]; ok {
		req.EnvelopeType = logcache.EnvelopeTypes(logcache.EnvelopeTypes_value[v[0]])
	}

	if v, ok := q["filter_template"]; ok {
		req.FilterTemplate = v[0]
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	resp, err := c.grpcClient.Read(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Envelopes.Batch, nil
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

// WithFilterTemplate sets the 'template_filter' query parameter to the given
// value. It defaults to empty.
func WithFilterTemplate(t string) ReadOption {
	return func(u *url.URL, q url.Values) {
		q.Set("filter_template", t)
	}
}
