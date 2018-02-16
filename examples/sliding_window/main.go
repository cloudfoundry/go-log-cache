package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

func main() {
	cfg := loadConfig()

	httpClient := newHTTPClient(cfg)

	client := logcache.NewClient(cfg.Addr, logcache.WithHTTPClient(httpClient))

	visitor := func(es []*loggregator_v2.Envelope) bool {
		fmt.Println("*********************Start Window********************")
		defer fmt.Println("**********************End Window*********************")
		for _, e := range es {
			if cfg.PrintTimestamps {
				fmt.Printf("%d\n", time.Unix(0, e.GetTimestamp()).Unix())
				continue
			}

			fmt.Printf("%+v\n", e)
		}
		return true
	}

	walker := logcache.BuildWalker(cfg.SourceID, client.Read)
	logcache.Window(
		context.Background(),
		visitor,
		walker,
		logcache.WithWindowWidth(cfg.WindowWidth),
		logcache.WithWindowInterval(cfg.WindowInterval),
		logcache.WithWindowStartTime(time.Unix(0, cfg.StartTime)),
	)
}

type config struct {
	Addr            string        `env:"ADDR, required"`
	AuthToken       string        `env:"AUTH_TOKEN, required"`
	SourceID        string        `env:"SOURCE_ID, required"`
	WindowInterval  time.Duration `env:"WINDOW_INTERVAL"`
	WindowWidth     time.Duration `env:"WINDOW_WIDTH"`
	StartTime       int64         `env:"START_TIME"`
	PrintTimestamps bool          `env:"PRINT_TIMESTAMP"`
}

func loadConfig() config {
	c := config{
		WindowWidth:    time.Hour,
		WindowInterval: time.Minute,
	}

	if err := envstruct.Load(&c); err != nil {
		log.Fatal(err)
	}

	if c.StartTime == 0 {
		c.StartTime = time.Now().Add(-c.WindowWidth).UnixNano()
	}

	return c
}

type HTTPClient struct {
	cfg    config
	client *http.Client
}

func newHTTPClient(c config) *HTTPClient {
	return &HTTPClient{cfg: c, client: http.DefaultClient}
}

func (h *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", h.cfg.AuthToken)
	return h.client.Do(req)
}
