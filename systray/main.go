package main

import (
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"yuubari_go"
)

var maxRetry int
var retryInterval int
var port int

func init() {
	flag.IntVar(&port, "port", 8099, "listen port")
	flag.IntVar(&maxRetry, "retry", 3, "max retry times")
	flag.IntVar(&retryInterval, "interval", 5, "retry interval (seconds)")
	flag.Parse()
}
func main() {
	iconData, _ := faviconIcoBytes()
	proxy := yuubari_go.NewYuubariGoProxyHandler(port, maxRetry, retryInterval)
	go proxy.Serve()
	systray.Run(func() {
		systray.SetIcon(iconData)
		systray.SetTitle("YuubariGo!")
		systray.SetTooltip(fmt.Sprintf("YuubariGo! (%d)", proxy.ErrCount))
		mQuit := systray.AddMenuItem("Quit", "Quit")
		go func() {
			<-mQuit.ClickedCh
			systray.Quit()
		}()
		go func() {
			for {
				<-proxy.ErrCountNotifyCh
				systray.SetTooltip(fmt.Sprintf("YuubariGo! (%d)", proxy.ErrCount))
			}
		}()
	}, func() {
	})
}
