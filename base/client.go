package base

import (
  "google.golang.org/grpc"
  "github.com/lucasmbaia/grpc-base/config"
  "github.com/lucasmbaia/grpc-base/zipkin"
  "google.golang.org/grpc/credentials"
  "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
  "github.com/lucasmbaia/grpc-base/utils/transaction"
  "github.com/lucasmbaia/grpc-base/utils"
)

type Config struct {
  Collector zipkin.Collector
}

func (c Config) ClientConnect() (*grpc.ClientConn, error) {
  var (
    opts		    []grpc.DialOption
    creds		    credentials.TransportCredentials
    err			    error
    unaryClientInterceptor  []grpc.UnaryClientInterceptor
  )

  if config.EnvConfig.GrpcSSL {
    if creds, err = credentials.NewClientTLSFromFile(config.EnvConfig.CAFile, config.EnvConfig.ServerNameAuthority); err != nil {
      return new(grpc.ClientConn), err
    }

    opts = []grpc.DialOption{
      grpc.WithTransportCredentials(creds),
      grpc.WithUnaryInterceptor(transaction.TrasactionClientInterceptor()),
    }
  } else {
    opts = []grpc.DialOption{
      grpc.WithInsecure(),
      grpc.WithUnaryInterceptor(transaction.TrasactionClientInterceptor()),
    }
  }

  unaryClientInterceptor = append(unaryClientInterceptor, transaction.TrasactionClientInterceptor())

  if config.EnvConfig.TracerClient {
    unaryClientInterceptor = append(unaryClientInterceptor, otgrpc.OpenTracingClientInterceptor(c.Collector.Tracer))
  }

  opts = append(opts, grpc.WithUnaryInterceptor(utils.ClientUnaryInterceptor(unaryClientInterceptor...)))
  return grpc.Dial(config.EnvLocal.LinkerdURL, opts...)
}
