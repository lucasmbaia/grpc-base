package base

import (
  "strconv"
  "net"
  "reflect"
  "log"

  "github.com/lucasmbaia/grpc-base/config"
  "github.com/lucasmbaia/grpc-base/consul"
  "google.golang.org/grpc/credentials"
  //"google.golang.org/grpc/reflection"
  "google.golang.org/grpc"
)

type ConfigCMD struct {
  SSL		    bool
  RegisterConsul    bool
  ServerConfig	    interface{}
  HandlerEndpoint   reflect.Value
  ServiceServer	    reflect.Value
}

func (c ConfigCMD) Run() error {
  var (
    listen		net.Listener
    err			error
    errChan		= make(chan error, 1)
    creds		credentials.TransportCredentials
    opts		[]grpc.ServerOption
    s			*grpc.Server
    args		[]reflect.Value
  )

  go func() {
    if listen, err = net.Listen(config.EnvConfig.TypeConnection, ":" + strconv.Itoa(config.EnvConfig.ServicePort)); err != nil {
      errChan <- err
      return
    }

    if c.SSL {
      if creds, err = credentials.NewServerTLSFromFile(config.EnvConfig.CertFile, config.EnvConfig.KeyFile); err != nil {
	errChan <- err
	return
      }

      opts = []grpc.ServerOption{
	grpc.Creds(creds),
      }
    }

    if c.RegisterConsul {
      if err = consul.RegisterService(); err != nil {
	errChan <- err
	return
      }
    }

    s = grpc.NewServer(opts...)

    args = append(args, reflect.ValueOf(s))
    args = append(args, reflect.ValueOf(c.ServerConfig))
    c.ServiceServer.Call(args)

    //reflection.Register(s)

    errChan <-s.Serve(listen)
    //errChan <-gateway(c.HandlerEndpoint, c.SSL)
  }()

  select {
  case e := <-errChan:
    return e
  }
}
