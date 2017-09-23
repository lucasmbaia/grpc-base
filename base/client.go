package base

import (
  "google.golang.org/grpc"
  "github.com/lucasmbaia/grpc-base/config"
  "github.com/lucasmbaia/grpc-base/zipkin"
  "google.golang.org/grpc/credentials"
)

type Config struct {
  Collector zipkin.Collector
}

func (c Config) ClientConnect() (*grpc.ClientConn, error) {
  var (
    opts  []grpc.DialOption
    creds credentials.TransportCredentials
    err   error
  )

  if config.EnvConfig.GrpcSSL {
    if creds, err = credentials.NewClientTLSFromFile(config.EnvConfig.CAFile, config.EnvConfig.ServerNameAuthority); err != nil {
      return new(grpc.ClientConn), err
    }

    opts = []grpc.DialOption{
      grpc.WithTransportCredentials(creds),
    }
  } else {
    opts = []grpc.DialOption{
      grpc.WithInsecure(),
    }
  }

  if config.EnvConfig.TracerClient {
    opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(c.Collector.Tracer)))
  }

  return grpc.Dial(config.EnvLocal.LinkerdURL, opts...)
}
