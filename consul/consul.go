package consul

import (
  "strings"
  "strconv"

  _consul "github.com/hashicorp/consul/api"
  "github.com/lucasmbaia/grpc-base/config"
)

const (
  URL_CHECK = "http://{url_check}:{port_check}/{endpoint_check}"
)

type consulConfig struct {
  ConsulURL	string
  ID            string
  Name          string
  Address       string
  UrlCheck      string
  Port          int
  IntervalCheck int
  Timeout       int
  Deregister    int
}

func RegisterService() error {
  var (
    urlCheck  *strings.Replacer
    service   consulConfig
  )

  urlCheck = strings.NewReplacer("{url_check}", config.EnvConfig.ServiceIPs[0], "{port_check}", config.EnvConfig.PortUrlCheck, "{endpoint_check}", config.EnvConfig.EndPointCheck)

  service = consulConfig{
    ConsulURL:	    config.EnvConfig.ConsulURL,
    ID:		    config.EnvConfig.Hostname,
    Name:	    config.EnvConfig.ServiceName,
    Address:	    config.EnvConfig.ServiceIPs[0],
    Port:	    config.EnvConfig.ServicePort,
    UrlCheck:	    urlCheck.Replace(URL_CHECK),
    IntervalCheck:  config.EnvConfig.ServiceIntervalCheck,
    Timeout:	    config.EnvConfig.ServiceTimeout,
    Deregister:	    config.EnvConfig.ServiceDeregister,
  }

  return register(service)
}

func initConsul(host string) (*_consul.Client, error) {
  var (
    config  = _consul.DefaultConfig()
  )

  config.Address = host

  return _consul.NewClient(config)
}

func register(c consulConfig) error {
  var (
    client  *_consul.Client
    err	    error
    service _consul.AgentServiceRegistration
  )

  if client, err = initConsul(c.ConsulURL); err != nil {
    return err
  }

  service = _consul.AgentServiceRegistration{
    ID:	      c.ID,
    Name:     c.Name,
    Port:     c.Port,
    Address:  c.Address,
    Check:    &_consul.AgentServiceCheck{
      HTTP:			      c.UrlCheck,
      Interval:			      strconv.Itoa(c.IntervalCheck) + "s",
      Timeout:			      strconv.Itoa(c.Timeout) + "s",
      DeregisterCriticalServiceAfter: strconv.Itoa(c.Deregister) + "s",
    },
  }

  return client.Agent().ServiceRegister(&service)
}
