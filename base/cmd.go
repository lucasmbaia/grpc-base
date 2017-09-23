package base

import (
  "strconv"
  "net"
  "reflect"
  "sync"

  "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
  "github.com/lucasmbaia/grpc-base/config"
  "github.com/lucasmbaia/grpc-base/consul"
  "github.com/lucasmbaia/grpc-base/zipkin"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/reflection"
  "google.golang.org/grpc"
)

type ConfigCMD struct {
  SSL		    bool
  RegisterConsul    bool
  RegisterRest	    bool
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
    wg			sync.WaitGroup
    collector		zipkin.Collector
  )

  wg.Add(2)

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

      opts = append(opts, grpc.Creds(creds))
    }

    if c.RegisterConsul {
      if err = consul.RegisterService(); err != nil {
	errChan <- err
	return
      }
    }

    if config.EnvConfig.TracerServer {
      if collector, err = zipkin.NewCollector(
	config.EnvConfig.ZipkinURL,
	config.EnvConfig.ServiceIPs[0] + ":" + strconv.Itoa(config.EnvConfig.ServicePort),
	config.EnvConfig.ServiceName,
	config.EnvConfig.DebugZipkin,
	config.EnvConfig.SameSpanZipkin,
      ); err != nil {
	errChan <- err
	return
      }

      config.EnvConfig.ZipKinTracer = collector
      opts = append(opts, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(collector.Tracer)))
    }

    s = grpc.NewServer(opts...)

    args = append(args, reflect.ValueOf(s))
    args = append(args, reflect.ValueOf(c.ServerConfig))
    c.ServiceServer.Call(args)

    reflection.Register(s)

    errChan <-s.Serve(listen)
    wg.Done()
  }()

  go func() {
    errChan <-onlyCheck()
    wg.Done()
  }()

  if c.RegisterRest {
    wg.Add(1)

    go func() {
      errChan <-gateway(c.HandlerEndpoint, c.SSL)
      wg.Done()
    }()
  }

  wg.Wait()

  select {
  case e := <-errChan:
    return e
  }
}
