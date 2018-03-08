package logcache_test

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"

	"code.cloudfoundry.org/go-log-cache"
	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

// Ensure logcache.Reader is fulfilled by Client.Read
var _ logcache.Reader = logcache.NewClient("").Read

func TestWalk(t *testing.T) {
	t.Parallel()

	r := &stubReader{}
	logcache.Walk(context.Background(), "some-id", func([]*loggregator_v2.Envelope) bool { return false }, r.read)

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

// If data comes in with a timestamp that is too new and other data is coming
// in with slightly older timestamps, we don't want to skip the data that came
// in a little later just because newer data arrived.
func TestWalkRejectsTooNewData(t *testing.T) {
	t.Parallel()

	r := &stubReader{
		envelopes: [][]*loggregator_v2.Envelope{
			{
				{Timestamp: 1},
				// Give future value. It should not take this value
				{Timestamp: time.Now().Add(time.Minute).UnixNano()},
			},
		},
		errs: []error{nil},
	}

	var called, es int
	logcache.Walk(context.Background(), "some-id", func(e []*loggregator_v2.Envelope) bool {
		defer func() { called++ }()
		es += len(e)
		return called == 0
	}, r.read)

	if len(r.starts) != 2 {
		t.Fatalf("expected starts to have 2 entries: %d", len(r.starts))
	}

	if r.starts[1] != 2 {
		t.Fatalf("expected to reject future/too new envelopes: %d", r.starts[1])
	}

	if es != 1 {
		t.Fatal("expected future/too new envelopes to be rejected")
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
	logcache.Walk(context.Background(), "some-id", func(b []*loggregator_v2.Envelope) bool {
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

	now := time.Now()
	times := []int64{
		now.Add(-3).UnixNano(),
		now.Add(-2).UnixNano(),
		now.Add(-1).UnixNano(),
		now.UnixNano(),
	}

	r := &stubReader{
		envelopes: [][]*loggregator_v2.Envelope{
			{
				{Timestamp: times[0]},
				{Timestamp: times[1]},
			},
			{
				{Timestamp: times[2]},
			},
			{
				{Timestamp: times[3]},
			},
		},
		errs: []error{nil, nil, nil},
	}

	var es []*loggregator_v2.Envelope
	logcache.Walk(context.Background(), "some-id", func(b []*loggregator_v2.Envelope) bool {
		es = append(es, b...)
		return true
	},
		r.read,
		logcache.WithWalkStartTime(time.Unix(0, times[0])),
		logcache.WithWalkEndTime(time.Unix(0, times[3])),
	)

	if len(r.sourceIDs) != 2 {
		t.Fatalf("expected read to be invoked 2 times: %d", len(r.sourceIDs))
	}

	if !reflect.DeepEqual(r.starts, []int64{times[0], times[2]}) {
		t.Fatalf("wrong starts: %v", r.starts)
	}

	if len(r.opts[0]) != 1 {
		t.Fatal("expected EndTime option to be set")
	}

	if len(es) != 3 {
		t.Fatalf("expected 3 envlopes: %d", len(es))
	}

	for i, x := range times[:len(times)-2] {
		if es[i].Timestamp != x {
			t.Fatalf("expected timestamp to equal %d: %d", x, es[i].Timestamp)
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
		context.Background(),
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

	if b.resetCalled != 1 {
		t.Fatalf("expected reset to be invoked 1 time: %d", b.resetCalled)
	}
}

func TestWalkCancels(t *testing.T) {
	t.Parallel()

	r := &stubReader{
		envelopes: [][]*loggregator_v2.Envelope{nil, {{Timestamp: 1}}},

		// Emulate the context being cancelled and the client returning errors
		// because of it.
		errs: []error{errors.New("some-error"), nil},
	}
	b := &stubBackoff{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var called int
	logcache.Walk(
		ctx,
		"some-id",
		func(b []*loggregator_v2.Envelope) bool {
			called++
			return false
		},
		r.read,
		logcache.WithWalkBackoff(b),
	)

	// No need to backoff because context is cancelled
	if len(b.errs) != 0 {
		t.Fatalf("expected backoff to be invoked 0 times: %d", len(b.errs))
	}
}

func TestWalkPassesOpts(t *testing.T) {
	t.Parallel()

	r := &stubReader{}
	logcache.Walk(
		context.Background(),
		"some-id",
		func(b []*loggregator_v2.Envelope) bool {
			return false
		},
		r.read,
		logcache.WithWalkLimit(99),
		logcache.WithWalkEnvelopeTypes(rpc.EnvelopeType_LOG, rpc.EnvelopeType_GAUGE),
	)

	u := &url.URL{}
	q := u.Query()
	for _, o := range r.opts[0] {
		o(u, q)
	}
	u.RawQuery = q.Encode()

	assertQueryParam(t, u, "limit", "99")
	assertQueryParam(t, u, "envelope_types", "LOG", "GAUGE")
}

type stubBackoff struct {
	errs          []error
	onErrReturn   bool
	onEmptyReturn bool
	resetCalled   int
}

func (s *stubBackoff) OnErr(err error) bool {
	s.errs = append(s.errs, err)
	return s.onErrReturn
}

func (s *stubBackoff) OnEmpty() bool {
	return s.onEmptyReturn
}

func (s *stubBackoff) Reset() {
	s.resetCalled++
}

type stubReader struct {
	sourceIDs []string
	starts    []int64
	opts      [][]logcache.ReadOption

	envelopes [][]*loggregator_v2.Envelope
	errs      []error
}

func newStubReader() *stubReader {
	return &stubReader{}
}

func (s *stubReader) read(ctx context.Context, sourceID string, start time.Time, opts ...logcache.ReadOption) ([]*loggregator_v2.Envelope, error) {
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
