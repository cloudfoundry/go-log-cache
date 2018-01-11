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
	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"google.golang.org/grpc"
)

type TC struct {
	*testing.T
	logCache *stubLogCache
	client   *logcache.Client
}

type TG struct {
	*testing.T
	logCache *stubGrpcLogCache
	client   *logcache.Client
}

func TestClientRead(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.Group("RESTful client", func() {
		o.BeforeEach(func(t *testing.T) TC {
			logCache := newStubLogCache()
			client := logcache.NewClient(logCache.addr())
			return TC{
				T:        t,
				logCache: logCache,
				client:   client,
			}
		})

		o.Spec("it reads from the client", func(t TC) {
			envelopes, err := t.client.Read(context.Background(), "some-id", time.Unix(0, 99))

			Expect(t, err).To(BeNil())
			Expect(t, envelopes).To(HaveLen(2))
			Expect(t, envelopes[0].Timestamp).To(Equal(int64(99)))
			Expect(t, envelopes[1].Timestamp).To(Equal(int64(100)))
			Expect(t, t.logCache.reqs).To(HaveLen(1))
			Expect(t, t.logCache.reqs[0].URL.Path).To(Equal("/v1/read/some-id"))

			assertQueryParam(t.T, t.logCache.reqs[0].URL, "start_time", "99")
			Expect(t, t.logCache.reqs[0].URL.Query()).To(HaveLen(1))
		})

		o.Spec("it reads with options", func(t TC) {
			_, err := t.client.Read(
				context.Background(),
				"some-id",
				time.Unix(0, 99),
				logcache.WithEndTime(time.Unix(0, 101)),
				logcache.WithLimit(103),
				logcache.WithEnvelopeType(rpc.EnvelopeTypes_LOG),
			)

			Expect(t, err).To(BeNil())
			Expect(t, t.logCache.reqs).To(HaveLen(1))
			Expect(t, t.logCache.reqs[0].URL.Path).To(Equal("/v1/read/some-id"))

			assertQueryParam(t.T, t.logCache.reqs[0].URL, "start_time", "99")
			assertQueryParam(t.T, t.logCache.reqs[0].URL, "end_time", "101")
			assertQueryParam(t.T, t.logCache.reqs[0].URL, "limit", "103")
			assertQueryParam(t.T, t.logCache.reqs[0].URL, "envelope_type", "LOG")
			Expect(t, t.logCache.reqs[0].URL.Query()).To(HaveLen(4))
		})

		o.Spec("it returns an error for a non 200 response", func(t TC) {
			t.logCache.statusCode = 500
			_, err := t.client.Read(context.Background(), "some-id", time.Unix(0, 99))

			Expect(t, err).To(Not(BeNil()))
		})

		o.Spec("it returns an error for an invalid response", func(t TC) {
			t.logCache.result = []byte("invalid")
			_, err := t.client.Read(context.Background(), "some-id", time.Unix(0, 99))

			Expect(t, err).To(Not(BeNil()))
		})

		o.Spec("it returns an error for an unknown URL", func(t TC) {
			client := logcache.NewClient("http://unknown.url")
			_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99))

			Expect(t, err).To(Not(BeNil()))
		})

		o.Spec("it returns an error for an invalid URL", func(t TC) {
			client := logcache.NewClient("-:-invalid")
			_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99))

			Expect(t, err).To(Not(BeNil()))
		})

		o.Spec("it returns an error when the context is cancelled", func(t TC) {
			t.logCache.block = true
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := t.client.Read(
				ctx,
				"some-id",
				time.Unix(0, 99),
				logcache.WithEndTime(time.Unix(0, 101)),
				logcache.WithLimit(103),
				logcache.WithEnvelopeType(rpc.EnvelopeTypes_LOG),
			)
			Expect(t, err).To(Not(BeNil()))
		})
	})

	o.Group("gRPC client", func() {
		o.BeforeEach(func(t *testing.T) TG {
			logCache := newStubGrpcLogCache()
			client := logcache.NewClient(
				logCache.addr(),
				logcache.WithViaGRPC(grpc.WithInsecure()),
			)
			return TG{
				T:        t,
				logCache: logCache,
				client:   client,
			}
		})

		o.Spec("it reads from the client", func(t TG) {
			endTime := time.Now()

			envelopes, err := t.client.Read(context.Background(), "some-id", time.Unix(0, 99),
				logcache.WithLimit(10),
				logcache.WithEndTime(endTime),
				logcache.WithEnvelopeType(rpc.EnvelopeTypes_LOG),
			)

			Expect(t, err).To(BeNil())
			Expect(t, envelopes).To(HaveLen(2))
			Expect(t, envelopes[0].Timestamp).To(Equal(int64(99)))
			Expect(t, envelopes[1].Timestamp).To(Equal(int64(100)))
			Expect(t, t.logCache.reqs).To(HaveLen(1))
			Expect(t, t.logCache.reqs[0].SourceId).To(Equal("some-id"))
			Expect(t, t.logCache.reqs[0].StartTime).To(Equal(int64(99)))
			Expect(t, t.logCache.reqs[0].EndTime).To(Equal(endTime.UnixNano()))
			Expect(t, t.logCache.reqs[0].Limit).To(Equal(int64(10)))
			Expect(t, t.logCache.reqs[0].EnvelopeType).To(Equal(rpc.EnvelopeTypes_LOG))
		})

		o.Spec("it returns an error when the context is cancelled", func(t TG) {
			t.logCache.block = true
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := t.client.Read(
				ctx,
				"some-id",
				time.Unix(0, 99),
				logcache.WithEndTime(time.Unix(0, 101)),
				logcache.WithLimit(103),
				logcache.WithEnvelopeType(rpc.EnvelopeTypes_LOG),
			)
			Expect(t, err).To(Not(BeNil()))
		})
	})
}

type stubLogCache struct {
	statusCode int
	server     *httptest.Server
	reqs       []*http.Request
	result     []byte
	block      bool
}

func newStubLogCache() *stubLogCache {
	s := &stubLogCache{
		statusCode: http.StatusOK,
		result: []byte(`{
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
	w.Write(s.result)
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
