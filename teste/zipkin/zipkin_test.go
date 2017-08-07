package main

import (
  "log"
  "testing"
  "github.com/lucasmbaia/grpc-base/zipkin"
)

func TestNewCollector(t *testing.T) {
  var(
    err	error
    collector zipkin.Collector
  )

  if collector, err = zipkin.NewCollector("http://172.16.95.113:9411/api/v1/spans", "127.0.0.1:0", "base", false); err != nil {
    log.Fatalf("Erro to create new collector: ", err)
  }

  log.Println(collector)
}
