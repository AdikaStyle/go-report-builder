package main

import (
	"github.com/AdikaStyle/go-report-builder/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("starting go-report-builder...")
	logrus.Info("parsing config...")
	config, err := cmd.ParseConfig()
	panicOnError(err)
	logrus.Infof("loaded config: %+v", config)

	logrus.Info("building Module...")
	module := cmd.NewModule(config)
	logrus.Info("module built successfully")

	logrus.Info("server started at port: %d...", config.ServerPort)
	panicOnError(module.Server.Start())
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
