package main

import (
	"flag"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/elazarl/goproxy.v1"
)

var maxRetry int
var retryInterval int

func init() {
	flag.IntVar(&maxRetry, "retry", 3, "max retry times")
	flag.IntVar(&retryInterval, "interval", 5, "retry interval (seconds)")
	flag.Parse()
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	client := http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			retryCount := 1
			req.RequestURI = ""
			resp, err := client.Do(req)
			for err != nil {
				log.Errorf("Network error, already tried %d times", retryCount)
				log.Error(err)
				if retryCount >= maxRetry {
					return req, goproxy.NewResponse(req, "application/json", 500, "")
				}
				time.Sleep(time.Duration(retryInterval) * time.Second)
				resp, err = client.Do(req)
				retryCount++
			}
			return req, resp
		})

	log.Fatal(http.ListenAndServe(":8299", proxy))
}
