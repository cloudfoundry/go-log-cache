package logcache_test

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"

	"code.cloudfoundry.org/go-log-cache"
	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

func TestWalk(t *testing.T) {
	t.Parallel()

	r := &stubReader{}
	logcache.Walk("some-id", func([]*loggregator_v2.Envelope) bool { return false }, r.read)

	if len(r.sourceIDs) != 1 {
		t.Fatal("expected read to be invoked once")
	}

	if r.sourceIDs[0] != "some-id" {
		t.Fatalf("expected sourceID to equal 'some-id': %s", r.sourceIDs[0])
	}

	if r.starts[0] != 0 {
		t.Fatalf("expected start to equal '0': %d", r.starts[0])
	}
}

func TestWalkUsesEndTime(t *testing.T) {
	t.Parallel()

	r := &stubReader{
		envelopes: [][]*loggregator_v2.Envelope{
			{
				{Timestamp: 1},
				{Timestamp: 2},
			},
			{
				{Timestamp: 3},
				{Timestamp: 4},
				{Timestamp: 5},
			},
		},
		errs: []error{nil, nil},
	}
	expected := make([][]*loggregator_v2.Envelope, len(r.envelopes))
	copy(expected, r.envelopes)

	var es [][]*loggregator_v2.Envelope
	logcache.Walk("some-id", func(b []*loggregator_v2.Envelope) bool {
		es = append(es, b)
		return true
	},
		r.read)

	if len(r.sourceIDs) != 3 {
		t.Fatalf("expected read to be invoked 3 times: %d", len(r.sourceIDs))
	}

	if !reflect.DeepEqual(r.sourceIDs, []string{"some-id", "some-id", "some-id"}) {
		t.Fatalf("wrong sourceIDs': %v", r.sourceIDs)
	}

	if !reflect.DeepEqual(r.starts, []int64{0, 3, 6}) {
		t.Fatalf("wrong starts: %v", r.starts)
	}

	if !reflect.DeepEqual(expected, es) {
		t.Fatalf("wrong envelopes: %v || %v", es, expected)
	}
}

func TestWalkWithinWindow(t *testing.T) {
	t.Parallel()

	r := &stubReader{
		envelopes: [][]*loggregator_v2.Envelope{
			{
				{Timestamp: 1},
				{Timestamp: 2},
			},
			{
				{Timestamp: 3},
			},
			{
				{Timestamp: 4},
			},
		},
		errs: []error{nil, nil, nil},
	}

	var es []*loggregator_v2.Envelope
	logcache.Walk("some-id", func(b []*loggregator_v2.Envelope) bool {
		es = append(es, b...)
		return true
	},
		r.read,
		logcache.WithWalkStartTime(time.Unix(0, 1)),
		logcache.WithWalkEndTime(time.Unix(0, 4)),
	)

	if len(r.sourceIDs) != 2 {
		t.Fatalf("expected read to be invoked 2 times: %d", len(r.sourceIDs))
	}

	if !reflect.DeepEqual(r.starts, []int64{1, 3}) {
		t.Fatalf("wrong starts: %v", r.starts)
	}

	if len(r.opts[0]) != 1 {
		t.Fatal("expected EndTime option to be set")
	}

	if len(es) != 3 {
		t.Fatalf("expected 3 envlopes: %d", len(es))
	}

	for i := 1; i < 4; i++ {
		if es[i-1].Timestamp != int64(i) {
			t.Fatalf("expected timestamp to equal %d: %d", i, es[i-1].Timestamp)
		}
	}
}

func TestWalkRetriesOnError(t *testing.T) {
	t.Parallel()

	r := &stubReader{
		envelopes: [][]*loggregator_v2.Envelope{nil, {{Timestamp: 1}}},
		errs:      []error{errors.New("some-error"), nil},
	}
	b := &stubBackoff{
		onErrReturn: true,
	}

	var called int
	logcache.Walk(
		"some-id",
		func(b []*loggregator_v2.Envelope) bool {
			called++
			return false
		},
		r.read,
		logcache.WithWalkBackoff(b),
	)

	if len(r.sourceIDs) != 2 {
		t.Fatalf("expected read to be invoked 2 times: %d", len(r.sourceIDs))
	}

	if called != 1 {
		t.Fatalf("expected visit to be invoked 1 time: %d", called)
	}

	if len(b.errs) != 1 {
		t.Fatalf("expected backoff to be invoked 1 time: %d", len(b.errs))
	}
}

func TestWalkPassesOpts(t *testing.T) {
	t.Parallel()

	r := &stubReader{}
	logcache.Walk(
		"some-id",
		func(b []*loggregator_v2.Envelope) bool {
			return false
		},
		r.read,
		logcache.WithWalkLimit(99),
		logcache.WithWalkEnvelopeType(rpc.EnvelopeTypes_LOG),
	)

	u := &url.URL{}
	q := u.Query()
	for _, o := range r.opts[0] {
		o(u, q)
	}

	if q.Get("limit") != "99" {
		t.Fatal("expected 'limit' to be set")
	}

	if q.Get("envelope_type") != "LOG" {
		t.Fatal("expected 'envelope_type' to be set")
	}
}

type stubBackoff struct {
	errs          []error
	onErrReturn   bool
	onEmptyReturn bool
}

func (s *stubBackoff) OnErr(err error) bool {
	s.errs = append(s.errs, err)
	return s.onErrReturn
}

func (s *stubBackoff) OnEmpty() bool {
	return s.onEmptyReturn
}

func (s *stubBackoff) Reset() {
}

type stubReader struct {
	sourceIDs []string
	starts    []int64
	opts      [][]logcache.ReadOption

	envelopes [][]*loggregator_v2.Envelope
	errs      []error
}

func (s *stubReader) read(sourceID string, start time.Time, opts ...logcache.ReadOption) ([]*loggregator_v2.Envelope, error) {
	s.sourceIDs = append(s.sourceIDs, sourceID)
	s.starts = append(s.starts, start.UnixNano())
	s.opts = append(s.opts, opts)

	if len(s.envelopes) != len(s.errs) {
		panic("envelopes and errs should have same len")
	}

	if len(s.envelopes) == 0 {
		return nil, nil
	}

	defer func() {
		s.envelopes = s.envelopes[1:]
		s.errs = s.errs[1:]
	}()

	return s.envelopes[0], s.errs[0]
}
