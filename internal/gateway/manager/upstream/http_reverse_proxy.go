package upstream

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	"api_gateway/pkg/buffer"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func (f *Factory) NewLoadBalanceReverseProxy(lb loadbalancer.LoadBalance, upstream *config.Upstream) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.URL.String())
		if err != nil || nextAddr == "" {
			log.Error().Msgf("Load Balance Poll is empty")
			return
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			log.Error().Msgf("Error When parse Balance Poll url")
			return
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return &httputil.ReverseProxy{
		Director:      director,
		FlushInterval: 100 * time.Millisecond,
		BufferPool:    buffer.NewBufferPool(),
		Transport: &http.Transport{
			DialContext:           dialer.DialContext,
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConnsPerHost:   upstream.MaxIdleConnsPerHost,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ReadBufferSize:        64 * 1024,
			WriteBufferSize:       64 * 1024,
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
			http.NotFoundHandler().ServeHTTP(writer, request)
			log.Error().Msgf("Reverse Proxy Error: %v", err)
		}}
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
