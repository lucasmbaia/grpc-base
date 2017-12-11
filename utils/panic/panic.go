package panic

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "github.com/getsentry/raven-go"
  "github.com/lucasmbaia/grpc-base/config"
  "errors"
)

type PanicHandler func(context.Context, interface{})

var _ grpc.UnaryServerInterceptor = PanicUnaryInterceptor
var _ grpc.StreamServerInterceptor = PanicStreamInterceptor
var addHandlers	[]PanicHandler

func panicError(err interface{}) error {
  sentry(err)
  return grpc.Errorf(codes.Internal, "panic: %v", err)
}

func PanicUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
  defer handleCrash(ctx, func(ctx context.Context, r interface{}) {
    err = panicError(r)
  })

  return handler(ctx, req)
}

func PanicStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
  defer handleCrash(stream.Context(), func(ctx context.Context, r interface{}) {
    err = panicError(r)
  })

  return handler(srv, stream)
}

func handleCrash(ctx context.Context, handler PanicHandler) {
  if r := recover(); r != nil {
    handler(ctx, r)

    if addHandlers != nil {
      for _, fn := range addHandlers {
	fn(ctx, r)
      }
    }
  }
}

func sentry(r interface{}) {
  if config.EnvConfig.SentryUrl != "" {
    switch err := r.(type) {
    case string:
      raven.CaptureErrorAndWait(errors.New(err), nil)
    case error:
      raven.CaptureErrorAndWait(err, nil)
    default:
      raven.CaptureErrorAndWait(errors.New("Unknow error"), nil)
    }
  }

  return
}
