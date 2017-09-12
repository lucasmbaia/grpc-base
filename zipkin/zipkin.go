package zipkin

import (
  "context"

  "github.com/lucasmbaia/grpc-base/config"
  opentracing "github.com/opentracing/opentracing-go"
  zipkinTracer "github.com/openzipkin/zipkin-go-opentracing"
)

type Collector struct {
  Conn	  zipkinTracer.Collector
  Tracer  opentracing.Tracer
}

type Span struct {
  Span	opentracing.Span
  Ctx	context.Context
}

func NewCollector(url, hostPort, endPoint string, debug bool) (Collector, error) {
  var (
    collector Collector
    err	      error
  )

  if collector.Conn, err = zipkinTracer.NewHTTPCollector(url); err != nil {
    return collector, err
  }

  if collector.Tracer, err = zipkinTracer.NewTracer(
    zipkinTracer.NewRecorder(collector.Conn, debug, hostPort, endPoint),
    zipkinTracer.ClientServerSameSpan(true),
    zipkinTracer.TraceID128Bit(true),
  ); err != nil {
    return collector, err
  }

  opentracing.SetGlobalTracer(collector.Tracer)

  return collector, nil
}

func (c Collector) Close() {
  c.Conn.Close()
}

func (c Collector) OpenFatherSpan(ctx context.Context, name string, tags map[string]string, parent opentracing.Span) Span {
  var (
    span	Span
  )

  if parent == nil {
    span.Span = c.Tracer.StartSpan(name)
  } else {
    span.Span = c.Tracer.StartSpan(name, opentracing.ChildOf(parent.Context()))
  }

  span.Ctx = opentracing.ContextWithSpan(ctx, span.Span)

  for key, value := range tags {
    span.Span.SetTag(key, value)
  }

  return span
}

func OpenChildSpan(ctx context.Context, name string, tags map[string]string, parent opentracing.Span) Span {
  var (
    span  Span
  )

  if parent == nil {
    span.Span, span.Ctx = opentracing.StartSpanFromContext(ctx, name)
  } else {
    span.Span, span.Ctx = opentracing.StartSpanFromContext(ctx, name, opentracing.ChildOf(parent.Context()))
  }

  for key, value := range tags {
    span.Span.SetTag(key, value)
  }

  return span
}

func (s Span) Event(event []string) {
  for _, value := range event {
    s.Span.LogEvent(value)
  }
}

func (s Span) Tag(key, value string) {
  s.Span.SetTag(key, value)
}

func OpenZipkin(ctx context.Context, name string) (Collector, Span, error) {
  var (
    collector Collector
    span      Span
    err       error
    tags      = make(map[string]string)
  )

  if config.EnvConfig.Tracer {
    if collector, err = NewCollector(config.EnvConfig.ZipkinURL, "0.0.0.0:0", config.EnvConfig.ServiceName, true); err != nil {
      return collector, span, err
    }

    tags = map[string]string{
      "Host":     config.EnvConfig.ServiceIPs[0],
      "Hostname": config.EnvConfig.Hostname,
    }

    span = OpenChildSpan(ctx, name, tags, nil)
  }

  return collector, span, nil
}
