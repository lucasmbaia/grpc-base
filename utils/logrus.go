package utils

import (
  "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
  "github.com/lucasmbaia/grpc-base/utils/logging"
)

func InitLogrus(unaryServerInterceptor *[]grpc.UnaryServerInterceptor, streamServerInterceptor *[]grpc.StreamServerInterceptor) {
  var (
    logrusEntry	  *logrus.Entry
  )

  logrusEntry = logrus.NewEntry(logrus.New())

  *unaryServerInterceptor = append(*unaryServerInterceptor, logging.UnaryServerInterceptor(logrusEntry))
  *streamServerInterceptor = append(*streamServerInterceptor, logging.StreamServerInterceptor(logrusEntry))

  return
}
