package yuubari_go

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/elazarl/goproxy.v1"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

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

func (p *ProxyHandler) ProxyWithRetry(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	retryCount := 1
	req.RequestURI = ""
	reqCopy := copyRequest(req)
	resp, err := p.client.Do(req)
	for err != nil {
		atomic.AddInt64(&p.ErrCount, 1)
		p.ErrCountNotifyCh <- struct{}{}
		log.Infof("error count: %d", atomic.LoadInt64(&p.ErrCount))
		if !strings.Contains(err.Error(), "EOF") || retryCount >= p.maxRetry {
			log.Error("reached max retries, abort")
			log.Error(err)
			return req, goproxy.NewResponse(req, "application/json", 500, "")
		}
		time.Sleep(time.Duration(p.retryInterval) * time.Second)
		resp, err = p.client.Do(copyRequest(reqCopy))
		retryCount++
	}
	return req, resp
}

type ProxyHandler struct {
	*goproxy.ProxyHttpServer
	port             int
	client           http.Client
	ErrCount         int64
	ErrCountNotifyCh chan struct{}
	maxRetry         int
	retryInterval    int
}

func NewYuubariGoProxyHandler(port int, maxRetry int, retryInterval int) *ProxyHandler {
	ret := ProxyHandler{
		ProxyHttpServer: goproxy.NewProxyHttpServer(),
		port:            port,
		client: http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		ErrCountNotifyCh: make(chan struct{}, 1),
		maxRetry:         maxRetry,
		retryInterval:    retryInterval,
	}
	ret.OnRequest().DoFunc(ret.ProxyWithRetry)
	return &ret
}

func (p *ProxyHandler) Serve() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", p.port), p)
}
