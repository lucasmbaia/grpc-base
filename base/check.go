package base

import (
  "net/http"
  "github.com/lucasmbaia/grpc-base/config"
)

func onlyCheck() error {
  mux := http.NewServeMux()
  mux.HandleFunc(config.EnvConfig.EndPointCheck, func(w http.ResponseWriter, r *http.Request){})
  return http.ListenAndServe(":" + config.EnvConfig.PortUrlCheck, mux)
}
