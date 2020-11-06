package main

import (
	"RsDnsTools/controller"
	"flag"
	"runtime"
	"os"
	log "github.com/Sirupsen/logrus"
	"io"
	"os/signal"
	"syscall"
	"fmt"
	"RsDnsTools/util"
)

// For auto version building
//  go build -ldflags "-X main.version=version"
var (
	Version		string
	BuildTime	string
)


func main() {

	var (
		configPath      string
		logPath         string
		isLogVerbose    bool
		processorNumber int

	)
	checkVersion := false

	flag.StringVar(&configPath, "c", "configs/config.json", "config file path")
	flag.StringVar(&logPath, "l", "", "log file path")
	flag.BoolVar(&isLogVerbose, "v", false, "verbose mode")
	flag.IntVar(&processorNumber, "p", runtime.NumCPU(), "number of processor to use")
	flag.BoolVar(&checkVersion, "V", false, "get version")

	flag.Parse()

	if checkVersion{
		fmt.Printf("BuildTag is: %s--%s\r\n", Version, BuildTime)
		return
	}

	config := util.NewConfig(configPath)
	fmt.Printf("%v\n", config)

	if isLogVerbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if logPath != "" {
		lf, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
		if err != nil {
			println("Logfile error: Please check your log file path")
		} else {
			log.SetOutput(io.MultiWriter(lf, os.Stdout))
		}
	}

	log.Info("RSFwd " + Version)

	runtime.GOMAXPROCS(processorNumber)



	server, _ := controller.NewServer(config.ServiceAddress, config)
	server.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping\n", s)

}
