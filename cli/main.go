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
var proxy string
var kcp string

func init() {
	flag.IntVar(&port, "port", 8099, "listen port")
	flag.IntVar(&maxRetry, "retry", 3, "max retry times")
	flag.IntVar(&retryInterval, "interval", 5, "retry interval (seconds)")
	flag.BoolVar(&debug, "debug", false, "enable debug log")
	flag.StringVar(&proxy, "proxy", "", "backend proxy url")
	flag.StringVar(&kcp, "kcp", "", "KCCacheProxy url")
	flag.Parse()
}

func main() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	proxy := yuubari_go.NewYuubariGoProxyHandler(port, maxRetry, retryInterval, proxy, kcp, func(errCnt int64) {
		log.Warnf("error count: %d", errCnt)
	})
	log.Fatal(proxy.Serve())
}
