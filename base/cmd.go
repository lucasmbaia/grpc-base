package base

import (
  "strconv"
  "net"
  "reflect"
  "sync"

  "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
  "github.com/lucasmbaia/grpc-base/utils/transaction"
  "github.com/lucasmbaia/grpc-base/config"
  "github.com/lucasmbaia/grpc-base/utils"
  "github.com/lucasmbaia/grpc-base/utils/panic"
  "github.com/lucasmbaia/grpc-base/consul"
  "github.com/lucasmbaia/grpc-base/zipkin"
  "github.com/getsentry/raven-go"
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

func init() {
  if config.EnvConfig.SentryUrl != "" {
    //raven.SetDSN("https://7275d55b562741898c85d24607986002:b16e18ae6e7d4ffeb2049914184db1c1@sentry.io/256694")
    raven.SetDSN(config.EnvConfig.SentryUrl)
  }
}

func (c ConfigCMD) Run() error {
  var (
    listen		    net.Listener
    err			    error
    errChan		    = make(chan error, 1)
    creds		    credentials.TransportCredentials
    opts		    []grpc.ServerOption
    s			    *grpc.Server
    args		    []reflect.Value
    wg			    sync.WaitGroup
    collector		    zipkin.Collector
    unaryServerInterceptor  []grpc.UnaryServerInterceptor
    streamServerInterceptor []grpc.StreamServerInterceptor
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
      unaryServerInterceptor = append(unaryServerInterceptor, otgrpc.OpenTracingServerInterceptor(collector.Tracer))
    }

    utils.InitLogrus(&unaryServerInterceptor, &streamServerInterceptor)
    unaryServerInterceptor = append(unaryServerInterceptor, panic.PanicUnaryInterceptor)
    unaryServerInterceptor = append(unaryServerInterceptor, transaction.TransactionServerInterceptor())
    streamServerInterceptor = append(streamServerInterceptor, panic.PanicStreamInterceptor)

    opts = append(opts, grpc.UnaryInterceptor(utils.UnaryInterceptor(unaryServerInterceptor...)))
    opts = append(opts, grpc.StreamInterceptor(utils.StreamInterceptor(streamServerInterceptor...)))

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
