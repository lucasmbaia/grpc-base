package utils

import (
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
