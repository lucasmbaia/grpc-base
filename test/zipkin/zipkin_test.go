package main

import (
  "context"
  "testing"
  "log"

  "github.com/lucasmbaia/grpc-base/zipkin"
)

func TestNewCollector(t *testing.T) {
  var (
    err	error
    collector zipkin.Collector
  )

  if collector, err = zipkin.NewCollector("http://172.16.95.113:9411/api/v1/spans", "127.0.0.1:0", "base", false); err != nil {
    log.Fatalf("Erro to create new collector: ", err)
  }

  log.Println(collector)
}

func TestFlowZipkin(t *testing.T) {
  var (
    collector zipkin.Collector
    err	      error
    tags      = make(map[string]string)
    ctx	      = context.Background()
  )

  if collector, err = newCollector("merda"); err != nil {
    log.Fatalf("Error to create new collector: ", err)
  }

  s := collector.OpenFatherSpan(ctx, "Start", tags, nil)

  collector2, _ := newCollector("god")
  s.Event([]string{"Pocs"})
  s2 := collector2.OpenChildSpan(s.Ctx, "Pocs", tags, nil)
  sc := collector2.OpenChildSpan(s2.Ctx, "cururu", tags, s2.Span)

  sc.Span.Finish()
  s2.Span.Finish()
  collector2.Close()

  collector3, _ := newCollector("war")
  s.Event([]string{"Coco"})
  s3 := collector3.OpenChildSpan(s.Ctx, "Coco", tags, nil)
  sc2 := collector3.OpenChildSpan(s3.Ctx, "pega", tags, nil)

  sc2.Span.Finish()
  s3.Span.Finish()
  collector3.Close()

  s.Span.Finish()
  collector.Close()
}

func newCollector(name string) (zipkin.Collector, error) {
  return zipkin.NewCollector("http://172.16.95.113:9411/api/v1/spans", "0.0.0.0:0", name, true)
}

/*func TestPorra(t *testing.T) {
  var (
    tags  = make(map[string]string)
  )

  //no fromhttprequest manda o god
  shit, _ := newCollector("shit")
  god, _ := newCollector("god")

  shitSpan := shit.OpenFatherSpan(context.Background(), "Start", tags, nil)
  shitSpan.Span.LogEvent("Call God")

  s, ctx := opentracing.StartSpanFromContext(shitSpan.Ctx, "God")

  spanMidleware := opentracing.SpanFromContext(ctx)

  spanMidleware.SetTag("Tamo", "Junto")

  req, _ := http.NewRequest("GET", "http://myservice/", nil)

  if err := shit.Tracer.Inject(
    spanMidleware.Context(),
    opentracing.TextMap,
    opentracing.HTTPHeadersCarrier(req.Header),
  ); err != nil {
    log.Println(err)
  }

  wireContext, _ := shit.Tracer.Extract(
    opentracing.TextMap,
    opentracing.HTTPHeadersCarrier(req.Header),
  )

  spanFrom := god.Tracer.StartSpan("God", ext.RPCServerOption(wireContext))
  opentracing.ContextWithSpan(ctx, spanFrom)

  spanFrom.Finish()
  spanMidleware.Finish()
  s.Finish()
  shitSpan.Span.Finish()
  god.Close()
  shit.Close()
}*/

