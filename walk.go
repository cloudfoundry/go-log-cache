package logcache

import (
	"io/ioutil"
	"log"
	"time"

	"code.cloudfoundry.org/go-log-cache/rpc/logcache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

// Reader reads envelopes from LogCache. It will be invoked by Walker several
// time to traverse the length of the cache.
type Reader func(
	sourceID string,
	start time.Time,
	opts ...ReadOption,
) ([]*loggregator_v2.Envelope, error)

// Visitor is invoked for each envelope batch. If the function returns false,
// it doesn't make any more requests. Otherwise it reaches out for the next
// batch of envelopes.
type Visitor func([]*loggregator_v2.Envelope) bool

// Walk reads from the LogCache until the Visitor returns false.
func Walk(sourceID string, v Visitor, r Reader, opts ...WalkOption) {
	c := &walkConfig{
		log:     log.New(ioutil.Discard, "", 0),
		backoff: AlwaysDoneBackoff{},
	}

	for _, o := range opts {
		o.configure(c)
	}

	var readOpts []ReadOption
	if !c.end.IsZero() {
		readOpts = append(readOpts, WithEndTime(c.end))
	}

	if c.limit != nil {
		readOpts = append(readOpts, WithLimit(*c.limit))
	}

	if c.envelopeType != nil {
		readOpts = append(readOpts, WithEnvelopeType(*c.envelopeType))
	}

	for {
		es, err := r(sourceID, time.Unix(0, c.start), readOpts...)
		if err != nil {
			c.log.Print(err)
			if !c.backoff.OnErr(err) {
				return
			}
			continue
		}

		if len(es) == 0 {
			if !c.backoff.OnEmpty() {
				return
			}
			continue
		}

		c.backoff.Reset()

		// If visitor is done or the next timestamp would be outside of our
		// window (only when end is set), then be done.
		if !v(es) || (!c.end.IsZero() && es[len(es)-1].Timestamp+1 >= c.end.UnixNano()) {
			return
		}

		c.start = es[len(es)-1].Timestamp + 1
	}
}

// WalkOption overrides defaults for Walk.
type WalkOption interface {
	configure(*walkConfig)
}

// WithWalkLogger is used to set the logger for the Walk. It defaults to
// not logging.
func WithWalkLogger(l *log.Logger) WalkOption {
	return walkOptionFunc(func(c *walkConfig) {
		c.log = l
	})
}

// WithWalkStartTime sets the start time of the query.
func WithWalkStartTime(t time.Time) WalkOption {
	return walkOptionFunc(func(c *walkConfig) {
		c.start = t.UnixNano()
	})
}

// WithWalkEndTime sets the end time of the query. Once reached, Walk will
// exit.
func WithWalkEndTime(t time.Time) WalkOption {
	return walkOptionFunc(func(c *walkConfig) {
		c.end = t
	})
}

// WithWalkLimit sets the limit of the query.
func WithWalkLimit(limit int) WalkOption {
	return walkOptionFunc(func(c *walkConfig) {
		c.limit = &limit
	})
}

// WithWalkEnvelopeType sets the envelope_type of the query.
func WithWalkEnvelopeType(t logcache.EnvelopeTypes) WalkOption {
	return walkOptionFunc(func(c *walkConfig) {
		c.envelopeType = &t
	})
}

// WithWalkBackoff sets the backoff strategy for an empty batch or error. It
// defaults to stopping on an error or empty batch via AlwaysDoneBackoff.
func WithWalkBackoff(b Backoff) WalkOption {
	return walkOptionFunc(func(c *walkConfig) {
		c.backoff = b
	})
}

// Backoff is used to determine what to do if there is an empty batch or
// error. If there is an error, it will be passed to the method OnErr. If there is
// not an error and just an empty batch, the method OnEmpty will be invoked. If
// the either method returns false, then Walk exits. On a successful read that
// has envelopes, Reset will be invoked.
type Backoff interface {
	OnErr(error) bool
	OnEmpty() bool
	Reset()
}

// AlwaysDoneBackoff returns false for both OnErr and OnEmpty.
type AlwaysDoneBackoff struct{}

// OnErr implements Backoff.
func (b AlwaysDoneBackoff) OnErr(error) bool {
	return false
}

// OnEmpty implements Backoff.
func (b AlwaysDoneBackoff) OnEmpty() bool {
	return false
}

// Reset implements Backoff.
func (b AlwaysDoneBackoff) Reset() {}

// AlwaysRetryBackoff returns true for both OnErr and OnEmpty after sleeping
// the given interval.
type AlwaysRetryBackoff struct {
	interval time.Duration
}

// NewAlwaysRetryBackoff returns a new AlwaysRetryBackoff.
func NewAlwaysRetryBackoff(interval time.Duration) AlwaysRetryBackoff {
	return AlwaysRetryBackoff{
		interval: interval,
	}
}

// OnErr implements Backoff.
func (b AlwaysRetryBackoff) OnErr(error) bool {
	time.Sleep(b.interval)
	return true
}

// OnEmpty implements Backoff.
func (b AlwaysRetryBackoff) OnEmpty() bool {
	time.Sleep(b.interval)
	return true
}

// Reset implements Backoff.
func (b AlwaysRetryBackoff) Reset() {}

type walkOptionFunc func(*walkConfig)

func (f walkOptionFunc) configure(c *walkConfig) {
	f(c)
}

type walkConfig struct {
	log          *log.Logger
	backoff      Backoff
	start        int64
	end          time.Time
	limit        *int
	envelopeType *logcache.EnvelopeTypes
}
