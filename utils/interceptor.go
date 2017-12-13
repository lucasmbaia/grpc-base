package utils

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc"
)

func ServerUnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
  var (
    lastInterceptor int
    size	    = len(interceptor)
  )

  if size == 1 {
    return interceptor[0]
  } else if size > 1 {
    lastInterceptor = size - 1

    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
      var (
	unaryHandler  grpc.UnaryHandler
	count	      int
      )

      unaryHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
	if count == lastInterceptor {
	  return handler(currentCtx, currentReq)
	}

	count++
	return interceptor[count](currentCtx, currentReq, info, unaryHandler)
      }

      return interceptor[0](ctx, req, info, unaryHandler)
    }
  }

  return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    return handler(ctx, req)
  }
}

func ServerStreamInterceptor(interceptor ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
  var (
    lastInterceptor int
    size	    = len(interceptor)
  )

  if size == 1 {
    return interceptor[0]
  } else if size > 1{
    lastInterceptor = size - 1

    return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
      var (
	streamHandler grpc.StreamHandler
	count	      int
      )

      streamHandler = func(currentSrv interface{}, currentStream grpc.ServerStream) error {
	if count == lastInterceptor {
	  return handler(currentSrv, currentStream)
	}

	count++
	return interceptor[count](currentSrv, currentStream, info, streamHandler)
      }

      return interceptor[0](srv, stream, info, streamHandler)
    }
  }

  return func(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    return handler(srv, stream)
  }
}

func ClientUnaryInterceptor(interceptor ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
  var (
    lastInterceptor int
    size	    = len(interceptor)
  )

  if size == 1 {
    return interceptor[0]
  } else {
    lastInterceptor = size - 1

    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
      var (
	unaryHandler  grpc.UnaryInvoker
	count	      int
      )

      unaryHandler = func(currentCtx context.Context, currentMethod string, currentReq, currentRepl interface{}, currentConn *grpc.ClientConn, currentOpts ...grpc.CallOption) error {
	if count == lastInterceptor {
	  return invoker(currentCtx, currentMethod, currentReq, currentRepl, currentConn, currentOpts...)
	}

	count++
	return interceptor[count](currentCtx, currentMethod, currentReq, currentRepl, currentConn, unaryHandler, currentOpts...)
      }

      return interceptor[0](ctx, method, req, reply, cc, unaryHandler, opts...)
    }
  }

  return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    return invoker(ctx, method, req, reply, cc, opts...)
  }
}

func ClientStreamInterceptor(interceptor ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
  var (
    lastInterceptor int
    size	    = len(interceptor)
  )

  if size == 1 {
    return interceptor[0]
  } else {
    lastInterceptor = size - 1

    return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
      var (
	streamHandler grpc.Streamer
	count	      int
      )

      streamHandler = func(currentCtx context.Context, currentDesc *grpc.StreamDesc, currentConn *grpc.ClientConn, currentMethod string, currentOpts ...grpc.CallOption) (grpc.ClientStream, error) {
	if count == lastInterceptor {
	  return streamer(currentCtx, currentDesc, currentConn, currentMethod, currentOpts...)
	}

	count++
	return interceptor[count](currentCtx, currentDesc, currentConn, currentMethod, streamHandler, currentOpts...)
      }

      return interceptor[0](ctx, desc, cc, method, streamHandler, opts...)
    }
  }

  return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
    return streamer(ctx, desc, cc, method, opts...)
  }
}
