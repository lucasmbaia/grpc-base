package utils

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc"
)

func UnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
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

func StreamInterceptor(interceptor ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
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
