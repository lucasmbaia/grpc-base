package logging

import (
  "google.golang.org/grpc/codes"
  "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
)

type errorToCode func(err error) codes.Code

func defaultErrorToCode(err error) codes.Code {
  return grpc.Code(err)
}

func defaultCodeToLevel(code codes.Code) logrus.Level {
  switch code {
  case codes.OK:
    return logrus.InfoLevel
  case codes.Canceled:
    return logrus.InfoLevel
  case codes.Unknown:
    return logrus.ErrorLevel
  case codes.InvalidArgument:
    return logrus.InfoLevel
  case codes.DeadlineExceeded:
    return logrus.WarnLevel
  case codes.NotFound:
    return logrus.InfoLevel
  case codes.AlreadyExists:
    return logrus.InfoLevel
  case codes.PermissionDenied:
    return logrus.WarnLevel
  case codes.Unauthenticated:
    return logrus.InfoLevel
  case codes.ResourceExhausted:
    return logrus.WarnLevel
  case codes.FailedPrecondition:
    return logrus.WarnLevel
  case codes.Aborted:
    return logrus.WarnLevel
  case codes.OutOfRange:
    return logrus.WarnLevel
  case codes.Unimplemented:
    return logrus.ErrorLevel
  case codes.Internal:
    return logrus.ErrorLevel
  case codes.Unavailable:
    return logrus.WarnLevel
  case codes.DataLoss:
    return logrus.ErrorLevel
  default:
    return logrus.ErrorLevel
  }
}

