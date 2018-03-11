package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	envstruct "code.cloudfoundry.org/go-envstruct"
	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

func main() {
	cfg := loadConfig()

	httpClient := newHTTPClient(cfg)

	client := logcache.NewShardGroupReaderClient(cfg.Addr, logcache.WithHTTPClient(httpClient))

	visitor := func(es []*loggregator_v2.Envelope) bool {
		for _, e := range es {
			fmt.Printf("%s\n", e.GetSourceId())
		}
		return true
	}

	logcache.Walk(
		context.Background(),
		cfg.GroupName,
		visitor,
		client.BuildReader(cfg.RequesterID),
		logcache.WithWalkBackoff(logcache.NewAlwaysRetryBackoff(time.Second)),
		logcache.WithWalkLogger(log.New(os.Stderr, "", 0)),
	)
}

type config struct {
	Addr        string `env:"ADDR, required"`
	AuthToken   string `env:"AUTH_TOKEN, required"`
	GroupName   string `env:"GROUP_NAME, required"`
	RequesterID uint64 `env:"REQUESTER_ID"`
}

func loadConfig() config {
	c := config{}

	if err := envstruct.Load(&c); err != nil {
		log.Fatal(err)
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
