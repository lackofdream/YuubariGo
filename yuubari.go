package yuubari_go

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/elazarl/goproxy.v1"
)

// removeProxyHeaders Copied from goproxy
func removeProxyHeaders(r *http.Request) {
	r.RequestURI = "" // this must be reset when serving a request with the client
	// If no Accept-Encoding header exists, Transport will add the headers it can accept
	// and would wrap the response body with the relevant reader.
	r.Header.Del("Accept-Encoding")
	// curl can add that, see
	// https://jdebp.eu./FGA/web-proxy-connection-header.html
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Proxy-Authenticate")
	r.Header.Del("Proxy-Authorization")
	// Connection, Authenticate and Authorization are single hop Header:
	// http://www.w3.org/Protocols/rfc2616/rfc2616.txt
	// 14.10 Connection
	//   The Connection general-header field allows the sender to specify
	//   options that are desired for that particular connection and MUST NOT
	//   be communicated by proxies over further connections.
	r.Header.Del("Connection")
}

// craftClientRequest craft client request from received server request
func craftClientRequest(req *http.Request) *http.Request {
	ret := req.Clone(context.Background())
	body, _ := io.ReadAll(req.Body)
	req.Body.Close()
	ret.Body = io.NopCloser(bytes.NewBuffer(body))
	req.Body = io.NopCloser(bytes.NewBuffer(body))
	removeProxyHeaders(ret)
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

// readResp get response body data while keeping response body readable
func readResp(response *http.Response) []byte {
	data, _ := io.ReadAll(response.Body)
	response.Body.Close()
	response.Body = io.NopCloser(bytes.NewBuffer(data))
	return data
}

func (p *ProxyHandler) ProxyWithRetry(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	retryCount := 0

	// check if MITMed connect request
	if len(req.URL.Host) == 0 {
		req.URL.Scheme = "http"
		req.URL.Host = req.Host
	}

	// if KCCacheProxy enabled and URL requested is static files, throw it to KCCacheProxy
	if strings.Contains(req.URL.Path, "/kcs/") || strings.Contains(req.URL.Path, "/kcs2/") {
		log.Debugf("TODO: throw this to KCP: %s", req.URL)
	}

	for retryCount <= p.maxRetry {
		log.Debugf("proxy request to %s", req.URL)
		resp, err := p.client.Do(craftClientRequest(req))
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
	ret.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		log.Debugf("connect received, host: %s", host)
		return goproxy.HTTPMitmConnect, host
	})
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
