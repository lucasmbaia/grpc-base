package logging

import (
  "testing"
  "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
)

func TestUnaryServerInterceptor(t *testing.T) {
  _ = grpc.NewServer(
    grpc.UnaryInterceptor(UnaryServerInterceptor(logrus.NewEntry(logrus.New()))),
    grpc.StreamInterceptor(StreamServerInterceptor(logrus.NewEntry(logrus.New()))),
  )
}
