package utils

import (
  //"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
  //"github.com/grpc-ecosystem/go-grpc-middleware/tags"
  //"github.com/sirupsen/logrus"
  "google.golang.org/grpc"
)

func InitLogrus(unaryServerInterceptor *[]grpc.UnaryServerInterceptor, streamServerInterceptor *[]grpc.StreamServerInterceptor) {
  /*var (
    logrusLogger  *logrus.Logger
    custom	  grpc_logrus.CodeToLevel
    opts	  []grpc_logrus.Option
    logrusEntry	  *logrus.Entry
  )

  logrusEntry = logrus.NewEntry(logrusLogger)

  opts = []grpc_logrus.Option {
    grpc_logrus.WithLevels(custom),
  }

  grpc_logrus.ReplaceGrpcLogger(logrusEntry)

  *unaryServerInterceptor = append(*unaryServerInterceptor, grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)))
  *unaryServerInterceptor = append(*unaryServerInterceptor, grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...))
  *streamServerInterceptor = append(*streamServerInterceptor, grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)))
  *streamServerInterceptor = append(*streamServerInterceptor, grpc_logrus.StreamServerInterceptor(logrusEntry, opts...))*/

  return
}
