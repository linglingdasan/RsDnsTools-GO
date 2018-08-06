package main

import (
	"github.com/miekg/dns"
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
var version string

func main() {

	var (
		configPath      string
		logPath         string
		isLogVerbose    bool
		processorNumber int

	)

	flag.StringVar(&configPath, "c", "./config.json", "config file path")
	flag.StringVar(&logPath, "l", "", "log file path")
	flag.BoolVar(&isLogVerbose, "v", false, "verbose mode")
	flag.IntVar(&processorNumber, "p", runtime.NumCPU(), "number of processor to use")

	flag.Parse()

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

	log.Info("RSFwd " + version)

	runtime.GOMAXPROCS(processorNumber)



	server, _ := NewServer("0.0.0.0:5553")
	server.Run()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, stopping\n", s)

}

//当请求进入时，需要判断转发到哪个fwd上（根据源IP和查询域名），需要确认ecs信息的添加方式（根据源ip所属区域）
func dnsFwd(w dns.ResponseWriter, req *dns.Msg) {

}