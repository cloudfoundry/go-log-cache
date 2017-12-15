package logcache_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache"
)

func TestClientRead(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	client := logcache.NewClient(logCache.addr())

	envelopes, err := client.Read("some-id", time.Unix(0, 99))

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(envelopes) != 2 {
		t.Fatalf("expected to receive 2 envlopes, got %d", len(envelopes))
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

	if logCache.reqs[0].URL.Query().Get("start_time") != "99" {
		t.Fatalf("expected query parameter 'start_time' to equal '99', but got '%s'", logCache.reqs[0].URL.Query().Get("start_time"))
	}

	if len(logCache.reqs[0].URL.Query()) != 1 {
		t.Fatalf("expected only a single query parameter, but got %d", len(logCache.reqs[0].URL.Query()))
	}
}

func TestClientReadWithOptions(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	client := logcache.NewClient(logCache.addr())

	_, err := client.Read(
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

	if logCache.reqs[0].URL.Query().Get("start_time") != "99" {
		t.Fatalf("expected query parameter 'start_time' to equal '99', but got '%s'", logCache.reqs[0].URL.Query().Get("start_time"))
	}

	if logCache.reqs[0].URL.Query().Get("end_time") != "101" {
		t.Fatalf("expected query parameter 'end_time' to equal '101', but got '%s'", logCache.reqs[0].URL.Query().Get("end_time"))
	}

	if logCache.reqs[0].URL.Query().Get("limit") != "103" {
		t.Fatalf("expected query parameter 'limit' to equal '103', but got '%s'", logCache.reqs[0].URL.Query().Get("limit"))
	}

	if logCache.reqs[0].URL.Query().Get("envelope_type") != "LOG" {
		t.Fatalf("expected query parameter 'envelope_type' to equal 'LOG', but got '%s'", logCache.reqs[0].URL.Query().Get("envelope_type"))
	}

	if len(logCache.reqs[0].URL.Query()) != 4 {
		t.Fatalf("expected only 4 query parameters, but got %d", len(logCache.reqs[0].URL.Query()))
	}
}

func TestClientReadNon200(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.statusCode = 500
	client := logcache.NewClient(logCache.addr())

	_, err := client.Read("some-id", time.Unix(0, 99))

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientReadInvalidResponse(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.result = []byte("invalid")
	client := logcache.NewClient(logCache.addr())

	_, err := client.Read("some-id", time.Unix(0, 99))

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientReadUnknownAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewClient("http://invalid.url")

	_, err := client.Read("some-id", time.Unix(0, 99))

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientReadInvalidAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewClient("-:-invalid")

	_, err := client.Read("some-id", time.Unix(0, 99))

	if err == nil {
		t.Fatal("expected an error")
	}
}

type stubLogCache struct {
	statusCode int
	server     *httptest.Server
	reqs       []*http.Request
	result     []byte
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
	s.reqs = append(s.reqs, r)
	w.WriteHeader(s.statusCode)
	w.Write(s.result)
}
