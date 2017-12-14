package logging

import (
  "path"
  "time"

  "github.com/lucasmbaia/grpc-base/utils/transaction"
  "github.com/sirupsen/logrus"
  "google.golang.org/grpc"
  "golang.org/x/net/context"
  "google.golang.org/grpc/codes"
)

type LoggingStream struct {
  grpc.ServerStream
  LoggingContext context.Context
}

func UnaryServerInterceptor(entry *logrus.Entry) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
    var (
      code    codes.Code
      level   logrus.Level
      fields  logrus.Fields
    )

    resp, err = handler(ctx, req)
    code = defaultErrorToCode(err)
    level = defaultCodeToLevel(code)

    fields = logrus.Fields {
      "grpc.code":  code.String(),
    }

    if err != nil {
      fields[logrus.ErrorKey] = err
    }

    printLog(
      loggerFields(ctx, entry, info.FullMethod),
      level,
      "Finished",
    )

    return resp, err
  }
}

func StreamServerInterceptor(entry *logrus.Entry) grpc.StreamServerInterceptor {
  return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    var (
      code    codes.Code
      level   logrus.Level
      fields  logrus.Fields
      err     error
      str     = &LoggingStream{ServerStream: stream, LoggingContext: stream.Context()}
    )

    err = handler(srv, str)
    code = defaultErrorToCode(err)
    level = defaultCodeToLevel(code)

    fields = logrus.Fields {
      "grpc.code":  code.String(),
    }

    if err != nil {
      fields[logrus.ErrorKey] = err
    }

    printLog(
      loggerFields(stream.Context(), entry, info.FullMethod),
      level,
      "Finished",
    )

    return err
  }
}

func loggerFields(ctx context.Context, entry *logrus.Entry, methodString string) *logrus.Entry {
  return entry.WithFields(
    logrus.Fields {
      "system":           "grpc",
      "span.kind":        "server",
      "grpc.service":     path.Dir(methodString)[1:],
      "grpc.method":      path.Base(methodString),
      "transaction.id":	  transaction.GetTransactionID(ctx),
      "grpc_start_time":  time.Now().Format(time.RFC3339),
    },
  )
}

func printLog(entry *logrus.Entry, level logrus.Level, format string, args ...interface{}) {
  switch level {
  case logrus.DebugLevel:
    entry.Debugf(format, args...)
  case logrus.InfoLevel:
    entry.Infof(format, args...)
  case logrus.WarnLevel:
    entry.Warningf(format, args...)
  case logrus.ErrorLevel:
    entry.Errorf(format, args...)
  case logrus.FatalLevel:
    entry.Fatalf(format, args...)
  case logrus.PanicLevel:
    entry.Panicf(format, args...)
  }
}

