package panic

import (
  "log"
  "testing"
  "google.golang.org/grpc"
  "golang.org/x/net/context"
  "github.com/getsentry/raven-go"
)

type serverStream struct {
  grpc.ServerStream
  ctx context.Context
}

func (s *serverStream) Context() context.Context {
  return s.ctx
}

func (s *serverStream) SendMsg(m interface{}) error {
  return nil
}

func (s *serverStream) RecvMsg(m interface{}) error {
  return nil
}

func init() {
  raven.SetDSN("https://7275d55b562741898c85d24607986002:b16e18ae6e7d4ffeb2049914184db1c1@sentry.io/256694")
}

func TestPanicUnaryInterceptor(t *testing.T) {
  var (
    err	error
    ctx	= context.Background()
  )

  unaryInfo := &grpc.UnaryServerInfo {
    FullMethod:	"Test.UnaryMethod",
  }

  unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
    panic("test panic")
  }

  if _, err = PanicUnaryInterceptor(ctx, "teste", unaryInfo, unaryHandler); err == nil {
    log.Fatalf("Unexpected success")
  }

  log.Println("Message of panic: ", grpc.ErrorDesc(err))
}

func TestPanicStreamInterceptor(t *testing.T) {
  var (
    err	    error
    service = struct{}{}
    stream  = &serverStream{ctx: context.Background()}
  )

  streamInfo := &grpc.StreamServerInfo {
    FullMethod:	    "Test.StreamMethod",
    IsServerStream: true,
  }

  streamHandler := func(srv interface{}, stream grpc.ServerStream) error {
    panic("test panic")
  }

  if err = PanicStreamInterceptor(service, stream, streamInfo, streamHandler); err == nil {
    log.Fatalf("Unexpected success")
  }
}
