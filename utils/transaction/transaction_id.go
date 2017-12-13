package transaction

import (
  "time"

  "google.golang.org/grpc/metadata"
  "github.com/satori/go.uuid"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
)

const (
  KEY_TRANSACTION = "transaction-id"
)

func TrasactionClientInterceptor() grpc.UnaryClientInterceptor {
  return func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    var	(
      err error
      md  metadata.MD
      ok  bool
    )

    if md, ok = metadata.FromIncomingContext(ctx); !ok {
      md = metadata.New(newMapID())
    } else {
      if _, ok = md[KEY_TRANSACTION]; !ok {
	md = metadata.New(newMapID())
      }
    }

    ctx = metadata.NewOutgoingContext(ctx, md)
    err = invoker(ctx, method, req, resp, cc, opts...)

    return err
  }
}

func TransactionServerInterceptor() grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
    var (
      md  metadata.MD
      ok bool
    )

    if md, ok = metadata.FromIncomingContext(ctx); !ok {
      md = metadata.New(newMapID())
    }

    ctx = metadata.NewIncomingContext(ctx, md)
    resp, err = handler(ctx, req)

    return resp, err
  }
}

func GetTransactionID(ctx context.Context) string {
  var (
    md	metadata.MD
    ok	bool
    id	string
  )

  if md, ok = metadata.FromIncomingContext(ctx); ok {
    if _, ok = md[KEY_TRANSACTION]; ok {
      id = md[KEY_TRANSACTION][0]
    }
  }

  return id
}

func newMapID() map[string]string {
  return map[string]string{
    KEY_TRANSACTION: newID(),
  }
}

func newID() string {
  return uuid.NewV4().String() + "-" + time.Now().Format("20060102150405")
}
