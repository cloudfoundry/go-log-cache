package logcache_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"google.golang.org/grpc"
)

// Assert that logcache.Reader is fulfilled by Client.Read
var _ logcache.Reader = logcache.Reader(logcache.NewClient("").Read)

func TestClientRead(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	client := logcache.NewClient(logCache.addr())

	envelopes, err := client.Read(context.Background(), "some-id", time.Unix(0, 99))

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(envelopes) != 2 {
		t.Fatalf("expected to receive 2 envelopes, got %d", len(envelopes))
	}

	if envelopes[0].Timestamp != 99 || envelopes[1].Timestamp != 100 {
		t.Fatal("wrong envelopes")
	}

	if len(logCache.reqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.reqs))
	}

	if logCache.reqs[0].URL.Path != "/v1/read/some-id" {
		t.Fatalf("expected Path '/v1/read/some-id' but got '%s'", logCache.reqs[0].URL.Path)
	}

	assertQueryParam(t, logCache.reqs[0].URL, "start_time", "99")

	if len(logCache.reqs[0].URL.Query()) != 1 {
		t.Fatalf("expected only a single query parameter, but got %d", len(logCache.reqs[0].URL.Query()))
	}
}

func TestGrpcClientRead(t *testing.T) {
	t.Parallel()
	logCache := newStubGrpcLogCache()
	client := logcache.NewClient(logCache.addr(), logcache.WithViaGRPC(grpc.WithInsecure()))

	endTime := time.Now()

	envelopes, err := client.Read(context.Background(), "some-id", time.Unix(0, 99),
		logcache.WithLimit(10),
		logcache.WithEndTime(endTime),
		logcache.WithEnvelopeType(rpc.EnvelopeTypes_LOG),
	)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(envelopes) != 2 {
		t.Fatalf("expected to receive 2 envelopes, got %d", len(envelopes))
	}

	if envelopes[0].Timestamp != 99 || envelopes[1].Timestamp != 100 {
		t.Fatal("wrong envelopes")
	}

	if len(logCache.reqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.reqs))
	}

	if logCache.reqs[0].SourceId != "some-id" {
		t.Fatalf("expected SourceId (%s) to equal %s", logCache.reqs[0].SourceId, "some-id")
	}

	if logCache.reqs[0].StartTime != 99 {
		t.Fatalf("expected StartTime (%d) to equal %d", logCache.reqs[0].StartTime, 99)
	}

	if logCache.reqs[0].EndTime != endTime.UnixNano() {
		t.Fatalf("expected EndTime (%d) to equal %d", logCache.reqs[0].EndTime, endTime.UnixNano())
	}

	if logCache.reqs[0].Limit != 10 {
		t.Fatalf("expected Limit (%d) to equal %d", logCache.reqs[0].Limit, 10)
	}

	if logCache.reqs[0].EnvelopeType != rpc.EnvelopeTypes_LOG {
		t.Fatalf("expected EnvelopeType (%v) to equal %v", logCache.reqs[0].EnvelopeType, rpc.EnvelopeTypes_LOG)
	}
}

func TestClientReadWithOptions(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	client := logcache.NewClient(logCache.addr())

	_, err := client.Read(
		context.Background(),
		"some-id",
		time.Unix(0, 99),
		logcache.WithEndTime(time.Unix(0, 101)),
		logcache.WithLimit(103),
		logcache.WithEnvelopeType(rpc.EnvelopeTypes_LOG),
	)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(logCache.reqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.reqs))
	}

	if logCache.reqs[0].URL.Path != "/v1/read/some-id" {
		t.Fatalf("expected Path '/v1/read/some-id' but got '%s'", logCache.reqs[0].URL.Path)
	}

	assertQueryParam(t, logCache.reqs[0].URL, "start_time", "99")
	assertQueryParam(t, logCache.reqs[0].URL, "end_time", "101")
	assertQueryParam(t, logCache.reqs[0].URL, "limit", "103")
	assertQueryParam(t, logCache.reqs[0].URL, "envelope_type", "LOG")

	if len(logCache.reqs[0].URL.Query()) != 4 {
		t.Fatalf("expected only 4 query parameters, but got %d", len(logCache.reqs[0].URL.Query()))
	}
}

func TestClientReadNon200(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.statusCode = 500
	client := logcache.NewClient(logCache.addr())

	_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99))

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientReadInvalidResponse(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.result["GET/v1/read/some-id"] = []byte("invalid")
	client := logcache.NewClient(logCache.addr())

	_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99))

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientReadUnknownAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewClient("http://invalid.url")

	_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99))

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientReadInvalidAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewClient("-:-invalid")

	_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99))

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientReadCancelling(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.block = true
	client := logcache.NewClient(logCache.addr())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Read(
		ctx,
		"some-id",
		time.Unix(0, 99),
		logcache.WithEndTime(time.Unix(0, 101)),
		logcache.WithLimit(103),
		logcache.WithEnvelopeType(rpc.EnvelopeTypes_LOG),
	)

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestGrpcClientReadCancelling(t *testing.T) {
	t.Parallel()
	logCache := newStubGrpcLogCache()
	logCache.block = true
	client := logcache.NewClient(logCache.addr(), logcache.WithViaGRPC(grpc.WithInsecure()))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Read(
		ctx,
		"some-id",
		time.Unix(0, 99),
		logcache.WithEndTime(time.Unix(0, 101)),
		logcache.WithLimit(103),
		logcache.WithEnvelopeType(rpc.EnvelopeTypes_LOG),
	)

	if err == nil {
		t.Fatal("expected an error")
	}
}

type stubLogCache struct {
	statusCode int
	server     *httptest.Server
	reqs       []*http.Request
	result     map[string][]byte
	block      bool
}

func newStubLogCache() *stubLogCache {
	s := &stubLogCache{
		statusCode: http.StatusOK,
		result: map[string][]byte{
			"GET/v1/read/some-id": []byte(`{
		"envelopes": {
			"batch": [
			    {
					"timestamp": 99,
					"sourceId": "some-id"
				},
			    {
					"timestamp": 100,
					"sourceId": "some-id"
				}
			]
		}
	}`),
		},
	}
	s.server = httptest.NewServer(s)
	return s
}

func (s *stubLogCache) addr() string {
	return s.server.URL
}

func (s *stubLogCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.block {
		var block chan struct{}
		<-block
	}

	s.reqs = append(s.reqs, r)
	w.WriteHeader(s.statusCode)
	w.Write(s.result[r.Method+r.URL.Path])
}

func assertQueryParam(t *testing.T, u *url.URL, name, value string) {
	t.Helper()
	if u.Query().Get(name) == value {
		return
	}

	t.Fatalf("expected query parameter '%s' to equal '%s', but got '%s'", name, value, u.Query().Get(name))
}

type stubGrpcLogCache struct {
	mu    sync.Mutex
	reqs  []*rpc.ReadRequest
	lis   net.Listener
	block bool
}

func newStubGrpcLogCache() *stubGrpcLogCache {
	s := &stubGrpcLogCache{}
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	s.lis = lis
	srv := grpc.NewServer()
	rpc.RegisterEgressServer(srv, s)
	go srv.Serve(lis)

	return s
}

func (s *stubGrpcLogCache) addr() string {
	return s.lis.Addr().String()
}

func (s *stubGrpcLogCache) Read(c context.Context, r *rpc.ReadRequest) (*rpc.ReadResponse, error) {
	if s.block {
		var block chan struct{}
		<-block
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.reqs = append(s.reqs, r)

	return &rpc.ReadResponse{
		Envelopes: &loggregator_v2.EnvelopeBatch{
			Batch: []*loggregator_v2.Envelope{
				{Timestamp: 99, SourceId: "some-id"},
				{Timestamp: 100, SourceId: "some-id"},
			},
		},
	}, nil
}

func (s *stubGrpcLogCache) requests() []*rpc.ReadRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := make([]*rpc.ReadRequest, len(s.reqs))
	copy(r, s.reqs)
	return r
}
