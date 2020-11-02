package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	reflect "reflect"
	"testing"

	"github.com/hypertrace/goagent/sdk/google.golang.org/grpc/internal/helloworld"
	"github.com/hypertrace/goagent/sdk/internal/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var _ helloworld.GreeterServer = server{}

type server struct {
	err          error
	replyHeader  metadata.MD
	replyTrailer metadata.MD
	*helloworld.UnimplementedGreeterServer
}

func (s server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	var reply *helloworld.HelloReply
	if s.err == nil {
		reply = &helloworld.HelloReply{Message: fmt.Sprintf("Hello %s", req.GetName())}
	}

	if err := grpc.SetTrailer(ctx, s.replyTrailer); err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("unable to send reply trailer: %v", err))
	}

	if err := grpc.SendHeader(ctx, s.replyHeader); err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("unable to send reply headers: %v", err))
	}

	return reply, s.err
}

// createDialer creates a connection to be used as context dialer in GRPC
// communication.
func createDialer(s *grpc.Server) func(context.Context, string) (net.Conn, error) {
	const bufSize = 1024 * 1024

	listener := bufconn.Listen(bufSize)
	conn := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	return conn
}

// jsonEqual compares the JSON from two strings.
func jsonEqual(a, b string) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal([]byte(a), &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(b), &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func TestResolveInOutType(t *testing.T) {
	assert.Equal(t, "response", resolveInType(true))
	assert.Equal(t, "request", resolveOutType(true))
	assert.Equal(t, "request", resolveInType(false))
	assert.Equal(t, "response", resolveOutType(false))
}

var _ stats.Handler = &mockHandler{}

type mockHandler struct {
	Spans []*mock.Span
}

func (h *mockHandler) HandleConn(context.Context, stats.ConnStats) {}

func (h *mockHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (h *mockHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {}

// TagRPC implements per-RPC context management.
func (h *mockHandler) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	s := mock.NewSpan()
	h.Spans = append(h.Spans, s)
	return mock.ContextWithSpan(ctx, s)
}