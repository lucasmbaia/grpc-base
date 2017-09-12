package utils

import (
  "encoding/json"
  "net"
)

func GetIPs() ([]string, error) {
  var (
    addrs []net.Addr
    err   error
    ips   []string
  )

  if addrs, err = net.InterfaceAddrs(); err != nil {
    return ips, err
  }

  for _, addr := range addrs {
    if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
      if ip.IP.To4() != nil {
	ips = append(ips, ip.IP.String())
      }
    }
  }

  return ips, nil
}

func ConvertArgsToString(i interface{}) string {
  var (
    body  []byte
    err   error
  )

  if body, err = json.Marshal(i); err != nil {
    return err.Error()
  }

  return string(body)
}
