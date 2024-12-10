package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"testing"
)

func TestInitLogger(t *testing.T) {
	l := &logrus.Logger{}
	l.SetOutput(io.MultiWriter(os.Stdout))
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	l.SetLevel(logrus.ErrorLevel)

	l.Error("test")
	l.Error("test")
}
