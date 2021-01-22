package main

import (
	"github.com/AdikaStyle/go-report-builder/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)

	config, err := cmd.ParseConfig()
	panicOnError(err)

	module := cmd.NewModule(config)

	panicOnError(module.Server.Start())
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
