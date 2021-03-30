package yuubari_go

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/elazarl/goproxy.v1"
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

type ProxyHandler struct {
	*goproxy.ProxyHttpServer
	port             int
	client           http.Client
	errCount         int64
	errCountNotifyCh chan struct{}
	maxRetry         int
	retryInterval    int
}

func (p *ProxyHandler) ProxyWithRetry(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	retryCount := 0
	// this is required otherwise it'll be error: Request.RequestURI can't be set in client requests
	req.RequestURI = ""

	for retryCount <= p.maxRetry {
		log.Debugf("proxy request to %s", req.URL)
		resp, err := p.client.Do(copyRequest(req))
		if err == nil {
			return req, resp
		}
		atomic.AddInt64(&p.errCount, 1)
		p.errCountNotifyCh <- struct{}{}
		if !strings.Contains(err.Error(), "EOF") &&
			!strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
			log.Errorf("unrecoverable error: %s", err)
			return req, goproxy.NewResponse(req, "application/json", 500, "")
		}
		time.Sleep(time.Duration(p.retryInterval) * time.Second)
		retryCount++
	}

	log.Error("reached max retries, abort")
	return req, goproxy.NewResponse(req, "application/json", 500, "")
}

func NewYuubariGoProxyHandler(port int, maxRetry int, retryInterval int, proxy string, onErrorCntIncr func(int64)) *ProxyHandler {
	ret := ProxyHandler{
		ProxyHttpServer: goproxy.NewProxyHttpServer(),
		port:            port,
		client: http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		errCountNotifyCh: make(chan struct{}, 1024),
		maxRetry:         maxRetry,
		retryInterval:    retryInterval,
	}
	ret.OnRequest().DoFunc(ret.ProxyWithRetry)
	go func() {
		for {
			<-ret.errCountNotifyCh
			onErrorCntIncr(ret.errCount)
		}
	}()
	if len(proxy) != 0 {
		ret.SetProxy(proxy)
	}
	return &ret
}

func (p *ProxyHandler) Serve() error {
	log.Info("色々試してみても、いいかしら？")
	log.Debug("デバッグモード オン")
	return http.ListenAndServe(fmt.Sprintf(":%d", p.port), p)
}

func (p *ProxyHandler) SetProxy(proxy string) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		log.Fatal(err)
	}
	p.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}

func (p *ProxyHandler) SetLogPath(path string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
}
