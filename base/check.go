package base

import (
  "strings"
  "net/http"
  "github.com/lucasmbaia/grpc-base/config"
)

func onlyCheck() error {
  var (
    mux	  *http.ServeMux
    url	  = config.EnvConfig.EndPointCheck
  )

  if strings.Split(url, "/")[0] != "" {
    url = "/" + url
  }

  mux = http.NewServeMux()
  mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request){})
  return http.ListenAndServe(":" + config.EnvConfig.PortUrlCheck, mux)
}
