package main

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	log "github.com/sirupsen/logrus"
	"yuubari_go"
)

var maxRetry int
var retryInterval int
var port int
var debug bool
var proxy string

func init() {
	flag.IntVar(&port, "port", 8099, "listen port")
	flag.IntVar(&maxRetry, "retry", 3, "max retry times")
	flag.IntVar(&retryInterval, "interval", 5, "retry interval (seconds)")
	flag.BoolVar(&debug, "debug", false, "enable debug log")
	flag.StringVar(&proxy, "proxy", "", "backend proxy url")
	flag.Parse()
}
func main() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	iconData, _ := faviconIcoBytes()
	systray.Run(func() {
		systray.SetIcon(iconData)
		systray.SetTitle("YuubariGo!")
		systray.SetTooltip("YuubariGo!")
		mQuit := systray.AddMenuItem("Quit", "Quit")
		go func() {
			<-mQuit.ClickedCh
			systray.Quit()
		}()
		proxy := yuubari_go.NewYuubariGoProxyHandler(port, maxRetry, retryInterval, proxy, func(errCnt int64) {
			systray.SetTitle(fmt.Sprintf("YuubariGo! (%d)", errCnt))
			systray.SetTooltip(fmt.Sprintf("YuubariGo! (%d)", errCnt))
		})
		proxy.SetLogPath("YuubariGo.log")
		go proxy.Serve()
	}, func() {
	})
}
