package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"yuubari_go"
)

var maxRetry int
var retryInterval int
var port int
var debug bool

func init() {
	flag.IntVar(&port, "port", 8099, "listen port")
	flag.IntVar(&maxRetry, "retry", 3, "max retry times")
	flag.IntVar(&retryInterval, "interval", 5, "retry interval (seconds)")
	flag.BoolVar(&debug, "debug", false, "enable debug log")
	flag.Parse()
}

func main() {
	log.Infoln("色々試してみても、いいかしら？")
	proxy := yuubari_go.NewYuubariGoProxyHandler(port, maxRetry, retryInterval)
	proxy.Verbose = debug
	log.Fatal(proxy.Serve())
}
