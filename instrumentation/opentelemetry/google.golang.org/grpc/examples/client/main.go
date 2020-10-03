// +build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	traceablegrpc "github.com/traceableai/goagent/instrumentation/opentelemetry/google.golang.org/grpc"
	"github.com/traceableai/goagent/instrumentation/opentelemetry/google.golang.org/grpc/examples"
	pb "github.com/traceableai/goagent/instrumentation/opentelemetry/google.golang.org/grpc/examples/helloworld"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	examples.InitTracer("grpc-client")

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			traceablegrpc.EnrichUnaryClientInterceptor(
				otelgrpc.UnaryClientInterceptor(global.TraceProvider().Tracer("ai.traceable")),
			),
		),
	)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %v", r.GetMessage())
}