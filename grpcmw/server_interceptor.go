package grpcmw

import (
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// StreamServerInterceptor represents a server interceptor for gRPC methods that
// return a stream. It allows chaining of `grpc.StreamServerInterceptor`
// and other `StreamServerInterceptor`.
type StreamServerInterceptor interface {
	// Interceptor chains all added interceptors into a single
	// `grpc.StreamServerInterceptor`.
	Interceptor() grpc.StreamServerInterceptor
	// AddGRPCInterceptor adds given interceptors to the chain.
	AddGRPCInterceptor(i ...grpc.StreamServerInterceptor) StreamServerInterceptor
	// AddInterceptor is a convenient way for adding `StreamServerInterceptor`
	// to the chain of interceptors.
	AddInterceptor(i ...StreamServerInterceptor) StreamServerInterceptor
}

// UnaryServerInterceptor represents a server interceptor for gRPC methods that
// return a single value. It allows chaining of `grpc.UnaryServerInterceptor`
// and other `UnaryServerInterceptor`.
type UnaryServerInterceptor interface {
	// Interceptor chains all added interceptors into a single
	// `grpc.UnaryServerInterceptor`.
	Interceptor() grpc.UnaryServerInterceptor
	// AddGRPCInterceptor adds given interceptors to the chain.
	AddGRPCInterceptor(i ...grpc.UnaryServerInterceptor) UnaryServerInterceptor
	// AddInterceptor is a convenient way for adding `UnaryServerInterceptor`
	// to the chain of interceptors.
	AddInterceptor(i ...UnaryServerInterceptor) UnaryServerInterceptor
}

type streamServerInterceptor struct {
	interceptors []grpc.StreamServerInterceptor
	lock         *sync.RWMutex
}

type unaryServerInterceptor struct {
	interceptors []grpc.UnaryServerInterceptor
	lock         *sync.RWMutex
}

// NewStreamServerInterceptor returns a new `StreamServerInterceptor`.
// It initializes its interceptor chain with `arr`.
// This implementation is thread-safe.
func NewStreamServerInterceptor(arr ...grpc.StreamServerInterceptor) StreamServerInterceptor {
	return &streamServerInterceptor{
		interceptors: arr,
		lock:         &sync.RWMutex{},
	}
}

func chainStreamServerInterceptor(current grpc.StreamServerInterceptor, info *grpc.StreamServerInfo, next grpc.StreamHandler) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return current(srv, stream, info, next)
	}
}

// Interceptor chains all added interceptors into a single
// `grpc.StreamServerInterceptor`.
//
// The `handler` passed to each interceptor is either the next interceptor or,
// for the last element of the chain, the target method.
func (si streamServerInterceptor) Interceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO: Find a more efficient way
		interceptor := handler
		si.lock.RLock()
		for idx := len(si.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainStreamServerInterceptor(si.interceptors[idx], info, interceptor)
		}
		si.lock.RUnlock()
		return interceptor(srv, ss)
	}
}

// AddGRPCInterceptor adds `arr` to the chain of interceptors.
func (si *streamServerInterceptor) AddGRPCInterceptor(arr ...grpc.StreamServerInterceptor) StreamServerInterceptor {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.interceptors = append(si.interceptors, arr...)
	return si
}

// AddInterceptor is a convenient way for adding `StreamServerInterceptor`
// to the chain of interceptors. It only calls the method `Interceptor`
// for each of them and append the return value to the chain.
func (si *streamServerInterceptor) AddInterceptor(arr ...StreamServerInterceptor) StreamServerInterceptor {
	si.lock.Lock()
	defer si.lock.Unlock()
	for _, i := range arr {
		si.interceptors = append(si.interceptors, i.Interceptor())
	}
	return si
}

// NewUnaryServerInterceptor returns a new `UnaryServerInterceptor`.
// It initializes its interceptor chain with `arr`.
// This implementation is thread-safe.
func NewUnaryServerInterceptor(arr ...grpc.UnaryServerInterceptor) UnaryServerInterceptor {
	return &unaryServerInterceptor{
		interceptors: arr,
		lock:         &sync.RWMutex{},
	}
}

func chainUnaryServerInterceptor(current grpc.UnaryServerInterceptor, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return current(ctx, req, info, next)
	}
}

// Interceptor chains all added interceptors into a single
// `grpc.UnaryServerInterceptor`.
//
// The `handler` passed to each interceptor is either the next interceptor or,
// for the last element of the chain, the target method.
func (ui *unaryServerInterceptor) Interceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Find a more efficient way
		interceptor := handler
		ui.lock.RLock()
		for idx := len(ui.interceptors) - 1; idx >= 0; idx-- {
			interceptor = chainUnaryServerInterceptor(ui.interceptors[idx], info, interceptor)
		}
		ui.lock.RUnlock()
		return interceptor(ctx, req)
	}
}

// AddGRPCInterceptor adds `arr` to the chain of interceptors.
func (ui *unaryServerInterceptor) AddGRPCInterceptor(arr ...grpc.UnaryServerInterceptor) UnaryServerInterceptor {
	ui.lock.Lock()
	defer ui.lock.Unlock()
	ui.interceptors = append(ui.interceptors, arr...)
	return ui
}

// AddInterceptor is a convenient way for adding `UnaryServerInterceptor`
// to the chain of interceptors. It only calls the method `Interceptor`
// for each of them and append the return value to the chain.
func (ui *unaryServerInterceptor) AddInterceptor(arr ...UnaryServerInterceptor) UnaryServerInterceptor {
	ui.lock.Lock()
	defer ui.lock.Unlock()
	for _, i := range arr {
		ui.interceptors = append(ui.interceptors, i.Interceptor())
	}
	return ui
}
