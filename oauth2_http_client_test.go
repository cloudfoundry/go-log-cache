package logcache_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"code.cloudfoundry.org/go-log-cache"
)

var _ logcache.HTTPClient = &logcache.Oauth2HTTPClient{}

func TestOauth2HTTPClient(t *testing.T) {
	t.Parallel()

	stubClient := newStubHTTPClient()

	c := logcache.NewOauth2HTTPClient(
		"http://oauth2.something.com",
		"my-user",
		"my-password",
		logcache.WithOauth2HTTPClient(stubClient),
	)
	req, err := http.NewRequest("GET", "http://some-target.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code to be 200: %d", resp.StatusCode)
	}

	var r oauth2Resp
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.TokenType != "bearer" {
		t.Fatalf("expected TokenType to equal bearer: %s", r.TokenType)
	}

	if r.AccessToken != "some-token" {
		t.Fatalf("expected AccessToken to equal some-token: %s", r.AccessToken)
	}

	if len(stubClient.reqs) != 2 {
		t.Fatalf("expected two requests: %d", len(stubClient.reqs))
	}

	if stubClient.reqs[0].Method != "POST" {
		t.Fatalf("expected method to equal POST: %s", stubClient.reqs[0].Method)
	}

	if stubClient.reqs[0].URL.Host != "oauth2.something.com" {
		t.Fatalf("expected Host to equal oauth2.something.com: %s", stubClient.reqs[0].URL.Host)
	}

	if stubClient.reqs[0].URL.Path != "/oauth/token" {
		t.Fatalf("expected Path to equal /oauth/token: %s", stubClient.reqs[0].URL.Path)
	}

	if stubClient.reqs[0].Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Fatalf("expected Header Content-Type to equal application/x-www-form-urlencoded: %s", stubClient.reqs[0].Header.Get("Content-Type"))
	}

	if stubClient.reqs[0].URL.Query().Get("client_id") != "my-user" {
		t.Fatalf("expected client_id to equal my-user: %s", stubClient.reqs[0].URL.Query().Get("client_id"))
	}

	if stubClient.reqs[0].URL.Query().Get("grant_type") != "client_credentials" {
		t.Fatalf("expected grant_type to equal client_credentials: %s", stubClient.reqs[0].URL.Query().Get("grant_type"))
	}

	if stubClient.reqs[0].URL.User == nil {
		t.Fatal("expected Username/Password to be set")
	}

	if stubClient.reqs[0].URL.User.Username() != "my-user" {
		t.Fatalf("expected Username to equal my-user: %s", stubClient.reqs[0].URL.User.Username())
	}

	if password, ok := stubClient.reqs[0].URL.User.Password(); !ok || password != "my-password" {
		t.Fatalf("expected Password to equal my-password: %s", password)
	}

	if stubClient.reqs[1].URL.String() != "http://some-target.com" {
		t.Fatalf("expected URL to equal http://some-target.com: %s", stubClient.reqs[1].URL.String())
	}

	if stubClient.reqs[1].Method != "GET" {
		t.Fatalf("expected method to equal GET: %s", stubClient.reqs[1].Method)
	}

	if stubClient.reqs[1].Header.Get("Authorization") != "bearer some-token" {
		t.Fatalf("expected Authorization header to equal bearer tome-token: %s", stubClient.reqs[1].Header.Get("Authorization"))
	}
}

func TestOauth2HTTPClientWithPasswordGrant(t *testing.T) {
	t.Parallel()

	stubClient := newStubHTTPClient()

	c := logcache.NewOauth2HTTPClient(
		"http://oauth2.something.com",
		"client",
		"client-secret",
		logcache.WithOauth2HTTPClient(stubClient),
		logcache.WithUser("user", "user-password"),
	)

	req, err := http.NewRequest("GET", "http://some-target.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code to be 200: %d", resp.StatusCode)
	}

	var r oauth2Resp
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if stubClient.reqs[0].Method != "POST" {
		t.Fatalf("expected method to equal POST: %s", stubClient.reqs[0].Method)
	}

	if stubClient.reqs[0].URL.Host != "oauth2.something.com" {
		t.Fatalf("expected Host to equal oauth2.something.com: %s", stubClient.reqs[0].URL.Host)
	}

	if stubClient.reqs[0].URL.Path != "/oauth/token" {
		t.Fatalf("expected Path to equal /oauth/token: %s", stubClient.reqs[0].URL.Path)
	}

	if stubClient.reqs[0].Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Fatalf("expected Header Content-Type to equal application/x-www-form-urlencoded: %s", stubClient.reqs[0].Header.Get("Content-Type"))
	}

	if stubClient.reqs[0].URL.Query().Get("client_id") != "client" {
		t.Fatalf("expected client_id to equal client: %s", stubClient.reqs[0].URL.Query().Get("client_id"))
	}

	if stubClient.reqs[0].URL.Query().Get("client_secret") != "client-secret" {
		t.Fatalf("expected client_secret to equal client-secret: %s", stubClient.reqs[0].URL.Query().Get("client_secret"))
	}

	if stubClient.reqs[0].URL.Query().Get("grant_type") != "password" {
		t.Fatalf("expected grant_type to equal password: %s", stubClient.reqs[0].URL.Query().Get("grant_type"))
	}

	if stubClient.reqs[0].URL.Query().Get("username") != "user" {
		t.Fatalf("expected username to equal user: %s", stubClient.reqs[0].URL.Query().Get("username"))
	}

	if stubClient.reqs[0].URL.Query().Get("password") != "user-password" {
		t.Fatalf("expected password to equal user-password: %s", stubClient.reqs[0].URL.Query().Get("password"))
	}
}

func TestOauth2HTTPReturnsErrorForNon200(t *testing.T) {
	t.Parallel()

	stubClient := newStubHTTPClient()
	stubClient.resps = append(stubClient.resps, &http.Response{
		StatusCode: 500,
		Body:       ioutil.NopCloser(&bytes.Buffer{}),
	})
	stubClient.errs = append(stubClient.errs, nil)

	c := logcache.NewOauth2HTTPClient(
		"http://oauth2.something.com",
		"my-user",
		"my-password",
		logcache.WithOauth2HTTPClient(stubClient),
	)
	req, err := http.NewRequest("GET", "http://some-target.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Do(req)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestOauth2HTTPClientReusesToken(t *testing.T) {
	t.Parallel()

	stubClient := newStubHTTPClient()

	c := logcache.NewOauth2HTTPClient(
		"http://oauth2.something.com",
		"my-user",
		"my-password",
		logcache.WithOauth2HTTPClient(stubClient),
	)
	req, err := http.NewRequest("GET", "http://some-target.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		_, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	}

	if len(stubClient.reqs) != 4 {
		t.Fatalf("expected to only get the token once: %d", len(stubClient.reqs))
	}
}

func TestOauth2HTTPClientFetchesNewTokenOn401(t *testing.T) {
	t.Parallel()

	stubClient := newStubHTTPClient()

	stubClient.resps = append(stubClient.resps, tokenResp(), &http.Response{StatusCode: 401})
	stubClient.errs = append(stubClient.errs, nil, nil)

	c := logcache.NewOauth2HTTPClient(
		"http://oauth2.something.com",
		"my-user",
		"my-password",
		logcache.WithOauth2HTTPClient(stubClient),
	)
	req, err := http.NewRequest("GET", "http://some-target.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// Returns 401 and induces a new Oauth2 call
	_, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if len(stubClient.reqs) != 5 {
		t.Fatalf("expected to only get the token again: %d", len(stubClient.reqs))
	}
}

func TestOauth2HTTPClientDoesNothingIfHeaderExists(t *testing.T) {
	t.Parallel()

	stubClient := newStubHTTPClient()

	c := logcache.NewOauth2HTTPClient(
		"http://oauth2.something.com",
		"my-user",
		"my-password",
		logcache.WithOauth2HTTPClient(stubClient),
	)
	req, err := http.NewRequest("GET", "http://some-target.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "something")

	_, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if len(stubClient.reqs) != 1 {
		t.Fatalf("expected to not get the token: %d", len(stubClient.reqs))
	}
}

func TestOauth2HTTPClientSurvivesRaceDetector(t *testing.T) {
	t.Parallel()

	stubClient := newStubHTTPClient()

	c := logcache.NewOauth2HTTPClient(
		"http://oauth2.something.com",
		"my-user",
		"my-password",
		logcache.WithOauth2HTTPClient(stubClient),
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()

		req, err := http.NewRequest("GET", "http://some-target.com", nil)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 10; i++ {
			_, err = c.Do(req)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	wg.Wait()

	req, err := http.NewRequest("GET", "http://some-target.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		_, err = c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	}
}

type stubHTTPClient struct {
	mu    sync.Mutex
	reqs  []*http.Request
	resps []*http.Response
	errs  []error
}

func newStubHTTPClient() *stubHTTPClient {
	return &stubHTTPClient{}
}

func (s *stubHTTPClient) Do(r *http.Request) (*http.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reqs = append(s.reqs, r)

	if len(s.resps) != len(s.errs) {
		panic("resps and errs have to always have the same length")
	}

	if len(s.resps) == 0 {
		return tokenResp(), nil
	}

	resp := s.resps[0]
	err := s.errs[0]

	s.resps = s.resps[1:]
	s.errs = s.errs[1:]

	return resp, err
}

func tokenResp() *http.Response {
	data, err := json.Marshal(oauth2Resp{
		TokenType:   "bearer",
		AccessToken: "some-token",
	})
	if err != nil {
		panic(err)
	}

	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(data)),
	}
}

type oauth2Resp struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}
