package logcache_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	rpc "code.cloudfoundry.org/go-log-cache/rpc/logcache_v1"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
)

// Assert that logcache.Reader is fulfilled by GroupReaderClient.BuildReader
var _ logcache.Reader = logcache.Reader(logcache.NewShardGroupReaderClient("").BuildReader(999))

func TestClientGroupRead(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.result["GET/v1/shard_group/some-name"] = []byte(`{
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
	}`)
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	reader := client.BuildReader(999)

	envelopes, err := reader(context.Background(), "some-name", time.Unix(0, 99))

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

	if logCache.reqs[0].URL.Path != "/v1/shard_group/some-name" {
		t.Fatalf("expected Path '/v1/shard_group/some-name' but got '%s'", logCache.reqs[0].URL.Path)
	}

	assertQueryParam(t, logCache.reqs[0].URL, "start_time", "99")
	assertQueryParam(t, logCache.reqs[0].URL, "requester_id", "999")

	if len(logCache.reqs[0].URL.Query()) != 2 {
		t.Fatalf("expected only two query parameters, but got %d", len(logCache.reqs[0].URL.Query()))
	}
}

func TestClientGroupReadWithOptions(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.result["GET/v1/shard_group/some-name"] = []byte(`{
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
	}`)
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	_, err := client.Read(
		context.Background(),
		"some-name",
		time.Unix(0, 99),
		999,
		logcache.WithEndTime(time.Unix(0, 101)),
		logcache.WithLimit(103),
		logcache.WithEnvelopeTypes(rpc.EnvelopeType_LOG),
	)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(logCache.reqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.reqs))
	}

	if logCache.reqs[0].URL.Path != "/v1/shard_group/some-name" {
		t.Fatalf("expected Path '/v1/shard_group/some-name' but got '%s'", logCache.reqs[0].URL.Path)
	}

	assertQueryParam(t, logCache.reqs[0].URL, "start_time", "99")
	assertQueryParam(t, logCache.reqs[0].URL, "end_time", "101")
	assertQueryParam(t, logCache.reqs[0].URL, "limit", "103")
	assertQueryParam(t, logCache.reqs[0].URL, "envelope_types", "LOG")
	assertQueryParam(t, logCache.reqs[0].URL, "requester_id", "999")

	if len(logCache.reqs[0].URL.Query()) != 5 {
		t.Fatalf("expected only 5 query parameters, but got %d", len(logCache.reqs[0].URL.Query()))
	}
}

func TestClientGroupReadNon200(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.statusCode = 500
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99), 999)

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientGroupReadInvalidResponse(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.result["GET/v1/group/some-name"] = []byte("invalid")
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	_, err := client.Read(context.Background(), "some-name", time.Unix(0, 99), 999)

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientGroupReadUnknownAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewShardGroupReaderClient("http://invalid.url")

	_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99), 999)

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientGroupReadInvalidAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewShardGroupReaderClient("-:-invalid")

	_, err := client.Read(context.Background(), "some-id", time.Unix(0, 99), 999)

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientGroupReadCancelling(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.block = true
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Read(
		ctx,
		"some-id",
		time.Unix(0, 99),
		999,
		logcache.WithEndTime(time.Unix(0, 101)),
		logcache.WithLimit(103),
		logcache.WithEnvelopeTypes(rpc.EnvelopeType_LOG),
	)

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestGrpcClientGroupRead(t *testing.T) {
	t.Parallel()
	logCache := newStubGrpcGroupReader()
	client := logcache.NewShardGroupReaderClient(logCache.addr(), logcache.WithViaGRPC(grpc.WithInsecure()))

	endTime := time.Now()

	envelopes, err := client.Read(context.Background(), "some-id", time.Unix(0, 99), 999,
		logcache.WithLimit(10),
		logcache.WithEndTime(endTime),
		logcache.WithEnvelopeTypes(rpc.EnvelopeType_LOG),
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

	if len(logCache.readReqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.readReqs))
	}

	if logCache.readReqs[0].RequesterId != 999 {
		t.Fatalf("expected RequesterId (%d) to equal %d", logCache.readReqs[0].RequesterId, 999)
	}

	if logCache.readReqs[0].StartTime != 99 {
		t.Fatalf("expected StartTime (%d) to equal %d", logCache.readReqs[0].StartTime, 99)
	}

	if logCache.readReqs[0].EndTime != endTime.UnixNano() {
		t.Fatalf("expected EndTime (%d) to equal %d", logCache.readReqs[0].EndTime, endTime.UnixNano())
	}

	if logCache.readReqs[0].Limit != 10 {
		t.Fatalf("expected Limit (%d) to equal %d", logCache.readReqs[0].Limit, 10)
	}

	if len(logCache.readReqs[0].EnvelopeTypes) == 0 {
		t.Fatalf("expected EnvelopeTypes to not be empty")
	}

	if logCache.readReqs[0].EnvelopeTypes[0] != rpc.EnvelopeType_LOG {
		t.Fatalf("expected EnvelopeTypes (%v) to equal %v", logCache.readReqs[0].EnvelopeTypes, rpc.EnvelopeType_LOG)
	}
}

func TestGrpcClientGroupReadCancelling(t *testing.T) {
	t.Parallel()
	logCache := newStubGrpcGroupReader()
	logCache.block = true
	client := logcache.NewShardGroupReaderClient(logCache.addr(), logcache.WithViaGRPC(grpc.WithInsecure()))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Read(
		ctx,
		"some-id",
		time.Unix(0, 99),
		999,
		logcache.WithEndTime(time.Unix(0, 101)),
		logcache.WithLimit(103),
		logcache.WithEnvelopeTypes(rpc.EnvelopeType_LOG),
	)

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientSetGroup(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.result["PUT/v1/shard_group/some-name"] = []byte("{}")
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	err := client.SetShardGroup(
		context.Background(),
		"some-name",
		"some-id-1",
		"some-id-2",
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(logCache.reqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.reqs))
	}

	if logCache.reqs[0].URL.Path != "/v1/shard_group/some-name" {
		t.Fatalf("expected Path '/v1/shard_group/some-name' but got '%s'", logCache.reqs[0].URL.Path)
	}

	if logCache.reqs[0].Method != "PUT" {
		t.Fatalf("expected Method to be PUT: %s", logCache.reqs[0].Method)
	}

	var g rpc.GroupedSourceIds
	r := bytes.NewReader(logCache.bodies[0])
	if err := jsonpb.Unmarshal(r, &g); err != nil {
		t.Fatalf("unable to unmarshal body: %s", err)
	}

	if !reflect.DeepEqual(g.SourceIds, []string{"some-id-1", "some-id-2"}) {
		t.Fatalf("expected some-id-1 and some-id-2: %v", g.SourceIds)
	}
}

func TestClientSetGroupUnknownAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewShardGroupReaderClient("http://invalid.url")

	err := client.SetShardGroup(context.Background(), "some-name", "some-id")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientSetGroupInvalidAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewShardGroupReaderClient("-:-invalid")

	err := client.SetShardGroup(context.Background(), "some-name", "some-id")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientSetGroupNon200(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.statusCode = 500
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	err := client.SetShardGroup(context.Background(), "some-name", "some-id")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientSetGroupCancelling(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.block = true
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.SetShardGroup(ctx, "some-name", "some-id")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestGrpcClientSetGroup(t *testing.T) {
	t.Parallel()
	logCache := newStubGrpcGroupReader()
	client := logcache.NewShardGroupReaderClient(logCache.addr(), logcache.WithViaGRPC(grpc.WithInsecure()))

	err := client.SetShardGroup(context.Background(), "some-name", "some-id")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(logCache.setReqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.setReqs))
	}

	if logCache.setReqs[0].Name != "some-name" {
		t.Fatalf("expected Name 'some-name' but got '%s'", logCache.setReqs[0].Name)
	}

	if !reflect.DeepEqual(logCache.setReqs[0].GetSubGroup().GetSourceIds(), []string{"some-id"}) {
		t.Fatalf("expected SourceId 'some-id' but got '%v'", logCache.setReqs[0].GetSubGroup().GetSourceIds())
	}

	logCache.addErr = errors.New("some-error")
	err = client.SetShardGroup(context.Background(), "some-name", "some-id")
	if err == nil {
		t.Fatal("expected err")
	}
}

func TestClientGroupMeta(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()

	expectedResp := &rpc.ShardGroupResponse{
		SubGroups: []*rpc.GroupedSourceIds{
			{
				SourceIds: []string{"a", "b"},
			},
		},
		RequesterIds: []uint64{1, 2},
	}

	data, err := json.Marshal(expectedResp)
	if err != nil {
		t.Fatal(err)
	}

	logCache.result["GET/v1/shard_group/some-name/meta"] = data
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	resp, err := client.ShardGroup(context.Background(), "some-name")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(logCache.reqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.reqs))
	}

	if logCache.reqs[0].URL.Path != "/v1/shard_group/some-name/meta" {
		t.Fatalf("expected Path '/v1/shard_group/some-name/meta' but got '%s'", logCache.reqs[0].URL.Path)
	}

	if logCache.reqs[0].Method != "GET" {
		t.Fatalf("expected Method to be GET: %s", logCache.reqs[0].Method)
	}

	if len(resp.SubGroups) != 1 {
		t.Fatalf(`expected to have a SubGroup: %d`, len(resp.SubGroups))
	}

	if !reflect.DeepEqual(resp.SubGroups[0].SourceIDs, []string{"a", "b"}) {
		t.Fatalf(`expected SourceIds to equal: ["a", "b"]: %s`, resp.SubGroups[0].SourceIDs)
	}

	if !reflect.DeepEqual(resp.RequesterIDs, []uint64{1, 2}) {
		t.Fatalf(`expected RequesterIds to equal: [1, 2]: %s`, resp.RequesterIDs)
	}
}

func TestClientGroupsUnknownAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewShardGroupReaderClient("http://invalid.url")

	_, err := client.ShardGroup(context.Background(), "some-name")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientGroupInvalidAddr(t *testing.T) {
	t.Parallel()
	client := logcache.NewShardGroupReaderClient("-:-invalid")

	_, err := client.ShardGroup(context.Background(), "some-name")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientGroupNon200(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.statusCode = 500
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	_, err := client.ShardGroup(context.Background(), "some-name")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientGroupInvalidResponse(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.result["GET/v1/shard_group/some-name/meta"] = []byte("invalid")
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	_, err := client.ShardGroup(context.Background(), "some-name")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestClientGroupCancelling(t *testing.T) {
	t.Parallel()
	logCache := newStubLogCache()
	logCache.block = true
	client := logcache.NewShardGroupReaderClient(logCache.addr())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.ShardGroup(ctx, "some-name")

	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestGrpcClientGroup(t *testing.T) {
	t.Parallel()
	logCache := newStubGrpcGroupReader()
	client := logcache.NewShardGroupReaderClient(logCache.addr(), logcache.WithViaGRPC(grpc.WithInsecure()))

	resp, err := client.ShardGroup(context.Background(), "some-name")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(logCache.groupReqs) != 1 {
		t.Fatalf("expected have 1 request, have %d", len(logCache.groupReqs))
	}

	if logCache.groupReqs[0].Name != "some-name" {
		t.Fatalf("expected Name 'some-name' but got '%s'", logCache.groupReqs[0].Name)
	}

	if len(resp.SubGroups) != 1 {
		t.Fatalf(`expected to have a SubGroup: %d`, len(resp.SubGroups))
	}

	if !reflect.DeepEqual(resp.SubGroups[0].SourceIDs, []string{"a", "b"}) {
		t.Fatalf(`expected SourceIds to equal: ["a", "b"]: %s`, resp.SubGroups[0].SourceIDs)
	}

	if !reflect.DeepEqual(resp.RequesterIDs, []uint64{1, 2}) {
		t.Fatalf(`expected RequesterIds to equal: [1, 2]: %s`, resp.RequesterIDs)
	}

	logCache.groupErr = errors.New("some-error")
	_, err = client.ShardGroup(context.Background(), "some-name")
	if err == nil {
		t.Fatal("expected err")
	}
}

type stubGrpcGroupReader struct {
	mu        sync.Mutex
	setReqs   []*rpc.SetShardGroupRequest
	addErr    error
	groupReqs []*rpc.ShardGroupRequest
	groupErr  error
	lis       net.Listener
	block     bool

	readReqs []*rpc.ShardGroupReadRequest
	readErr  error
}

func newStubGrpcGroupReader() *stubGrpcGroupReader {
	s := &stubGrpcGroupReader{}
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	s.lis = lis
	srv := grpc.NewServer()
	rpc.RegisterShardGroupReaderServer(srv, s)
	go srv.Serve(lis)

	return s
}

func (s *stubGrpcGroupReader) addr() string {
	return s.lis.Addr().String()
}

func (s *stubGrpcGroupReader) SetShardGroup(ctx context.Context, r *rpc.SetShardGroupRequest) (*rpc.SetShardGroupResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setReqs = append(s.setReqs, r)
	return &rpc.SetShardGroupResponse{}, s.addErr
}

func (s *stubGrpcGroupReader) Read(ctx context.Context, r *rpc.ShardGroupReadRequest) (*rpc.ShardGroupReadResponse, error) {
	if s.block {
		var block chan struct{}
		<-block
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.readReqs = append(s.readReqs, r)

	return &rpc.ShardGroupReadResponse{
		Envelopes: &loggregator_v2.EnvelopeBatch{
			Batch: []*loggregator_v2.Envelope{
				{Timestamp: 99, SourceId: "some-id"},
				{Timestamp: 100, SourceId: "some-id"},
			},
		},
	}, nil
}

func (s *stubGrpcGroupReader) ShardGroup(ctx context.Context, r *rpc.ShardGroupRequest) (*rpc.ShardGroupResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.groupReqs = append(s.groupReqs, r)
	return &rpc.ShardGroupResponse{
		SubGroups: []*rpc.GroupedSourceIds{
			{
				SourceIds: []string{"a", "b"},
			},
		},
		RequesterIds: []uint64{1, 2},
	}, s.groupErr
}
