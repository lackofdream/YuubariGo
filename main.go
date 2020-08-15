package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/elazarl/goproxy.v1"
)

var maxRetry int
var retryInterval int
var port int
var client http.Client
var proxy *goproxy.ProxyHttpServer
var debug bool
var errCount int64

func init() {
	flag.IntVar(&port, "port", 8099, "listen port")
	flag.IntVar(&maxRetry, "retry", 3, "max retry times")
	flag.IntVar(&retryInterval, "interval", 5, "retry interval (seconds)")
	flag.BoolVar(&debug, "debug", false, "enable debug log")
	flag.Parse()
	client = http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	proxy = goproxy.NewProxyHttpServer()
	log.SetLevel(log.InfoLevel)
	if debug {
		proxy.Verbose = true
		log.SetLevel(log.DebugLevel)
	}
}

func copyRequest(req *http.Request) *http.Request {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = req.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	ret := req.Clone(context.Background())
	ret.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return ret
}

func ProxyWithRetry(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	retryCount := 1
	req.RequestURI = ""
	reqCopy := copyRequest(req)
	resp, err := client.Do(req)
	for err != nil {
		atomic.AddInt64(&errCount, 1)
		log.Infof("error count: %d", atomic.LoadInt64(&errCount))
		if !strings.Contains(err.Error(), "EOF") || retryCount >= maxRetry {
			log.Error("reached max retries, abort")
			log.Error(err)
			return req, goproxy.NewResponse(req, "application/json", 500, "")
		}
		time.Sleep(time.Duration(retryInterval) * time.Second)
		resp, err = client.Do(copyRequest(reqCopy))
		retryCount++
	}
	return req, resp
}

func main() {
	proxy.OnRequest().DoFunc(ProxyWithRetry)
	log.Infoln("色々試してみても、いいかしら？")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), proxy))
}
