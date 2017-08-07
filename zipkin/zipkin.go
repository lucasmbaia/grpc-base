package zipkin

import (
  "log"

  opentracing "github.com/opentracing/opentracing-go"
  zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

type Collector struct {
  Conn	  zipkin.Collector
  Tracer  opentracing.Tracer
}

func NewCollector(url, hostPort, endPoint string, debug bool) (Collector, error) {
  var (
    collector Collector
    err	      error
  )

  if collector.Conn, err = zipkin.NewHTTPCollector(url); err != nil {
    return collector, nil
  }

  if collector.Tracer, err = zipkin.NewTracer(zipkin.NewRecorder(collector.Conn, debug, hostPort, endPoint)); err != nil {
    return collector, nil
  }

  apply := &StartSpanOptions{Tags: map[string]interface{}{"teste":"teste"}}
  log.Println(collector.Tracer.StartSpan("Parent", opentracing.StartSpanOption{opentracing.Apply(apply)}))
    return collector, nil
  }
