package app

import (
	"fmt"

	"github.com/ipoluianov/gomisc/logger"
	"github.com/ipoluianov/map_u00_io/httpserver"
)

func Start() {
	logger.Println("Start begin")
	TuneFDs()

	httpserver.Instance.Start()

	logger.Println("Start end")
}

func Stop() {
}

func RunDesktop() {
	logger.Println("Running as console application")
	Start()
	fmt.Scanln()
	logger.Println("Console application exit")
}

func RunAsService() error {
	Start()
	return nil
}

func StopService() {
	Stop()
}
