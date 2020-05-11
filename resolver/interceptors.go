package resolver

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func withClientUnaryInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientUnaryInterceptor)
}

func clientUnaryInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Logic before invoking the invoker
	start := time.Now()
	clientDeadline := time.Now().Add(time.Duration(5000) * time.Millisecond)
	_ctx, _ := context.WithDeadline(ctx, clientDeadline)

	// Calls the invoker to execute RPC
	err := invoker(_ctx, method, req, reply, cc, opts...)

	// Logic after invoking the invoker
	log.Printf("Invoked RPC method=%s; Duration=%s; Error=%v", method, time.Since(start), err)

	return err
}
