package main

import (
	"github.com/ipoluianov/gomisc/logger"
	"github.com/ipoluianov/map_u00_io/app"
	"github.com/ipoluianov/map_u00_io/application"
)

func main() {
	name := "aneth_eth"
	application.Name = name
	application.ServiceName = name
	application.ServiceDisplayName = name
	application.ServiceDescription = name
	application.ServiceRunFunc = app.RunAsService
	application.ServiceStopFunc = app.StopService
	logger.Init(logger.CurrentExePath() + "/logs")

	if !application.TryService() {
		app.RunDesktop()
	}
}
