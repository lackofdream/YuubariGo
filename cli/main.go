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
var expedNotify bool
var tgToken string
var tgUser int64
var teitokuName string

func init() {
	flag.IntVar(&port, "port", 8099, "listen port")
	flag.IntVar(&maxRetry, "retry", 3, "max retry times")
	flag.IntVar(&retryInterval, "interval", 5, "retry interval (seconds)")
	flag.BoolVar(&debug, "debug", false, "enable debug log")
	flag.StringVar(&proxy, "proxy", "", "backend proxy url")
	flag.StringVar(&kcp, "kcp", "", "KCCacheProxy url")
	flag.BoolVar(&expedNotify, "expedNotify", false, "enable exped_notify plugin")
	flag.StringVar(&tgToken, "tgToken", "", "telegram bot token")
	flag.StringVar(&teitokuName, "teitokuName", "", "teitoku nickname for exped_notify")
	flag.Int64Var(&tgUser, "tgUser", 0, "telegram user id to be notified")
	flag.Parse()
}

func main() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	instance := yuubari_go.NewYuubariGoProxyHandler(port, maxRetry, retryInterval, proxy, kcp, func(errCnt int64) {
		log.Warnf("error count: %d", errCnt)
	})
	if expedNotify {
		instance = yuubari_go.MakeResourceLogged(yuubari_go.MakeNotifiable(instance, proxy, tgToken, tgUser, teitokuName))
	}
	log.Fatal(instance.Serve())
}
