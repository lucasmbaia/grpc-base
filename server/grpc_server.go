package server

import (
	"google.golang.org/grpc"
	"reflect"
	"google.golang.org/grpc/credentials"
)

const (
	defaultPortGRPC = 5222
)

type ServerGRPC struct {
	server		*grpc.Server
	port		int
	certificate	string
	key		string
	services	[]Services
}

type ServerGRPCConfig struct {
	Port		int
	Certificate	string
	Key		string
	Services	[]Services
}

type Services struct {
	Server		interface{}
	Contract	interface{}
}

func NewServerGRPC(cfg ServerGRPCConfig) (s ServerGRPC, err error) {
	if cfg.Port == 0 {
		s.port = defaultPortGRPC
	} else {
		s.port = cfg.Port
	}

	s.certificate = cfg.Certificate
	s.key = cfg.Key

	for _, service := range cfg.Services {
		if service.Server == nil {
			err = errors.New("Service is empty")
		}

		if service.Contract == nil {
			err = errors.New("Contract is empty")
		}

		if err != nil {
			return
		}
	}

	s.services = cfg.Services

	return
}

func (s *ServerGRPC) Listen() (err error) {
	var (
		l	  net.Listener
		grpcOpts  = []grpc.ServerOption{}
		file	  *os.File
		grpcCreds credentials.TransportCredentials
	)

	if l, err = net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(s.port))); err != nil {
		return
	}

	if s.certificate != "" && s.key != "" {
		if grpcCreds, err = credentials.NewServerTLSFromFile(s.certificate, s.key); err != nil {
			return
		}

		grpcOpts = append(grpcOpts, grpc.Creds(grpcCreds))
	}

	s.server = grpc.NewServer(grpcOpts...)

	for _, service := range s.services {
		var args []reflect.Value

		args = append(args, reflect.ValueOf(s.server))
		args = append(args, reflect.ValueOf(service.Contract))

		reflect.Indirect(reflect.ValueOf(service.Server)).Call(args)
		reflection.Register(s.server)
	}

	if err = s.server.Serve(l); err != nil {
		return
	}

	return
}
