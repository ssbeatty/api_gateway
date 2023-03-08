package redirect

import (
	"api_gateway/pkg/middlewares/logs"
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
)

const (
	TypeSchemeName  = "RedirectScheme"
	uriPattern      = `^(https?:\/\/)?(\[[\w:.]+\]|[\w\._-]+)?(:\d+)?(.*)$`
	xForwardedProto = "X-Forwarded-Proto"
)

type Scheme struct {
	// Scheme defines the scheme of the new URL.
	Scheme string `json:"scheme,omitempty" toml:"scheme,omitempty" yaml:"scheme,omitempty" export:"true"`
	// Port defines the port of the new URL.
	Port string `json:"port,omitempty" toml:"port,omitempty" yaml:"port,omitempty" export:"true"`
	// Permanent defines whether the redirection is permanent (301).
	Permanent bool `json:"permanent,omitempty" toml:"permanent,omitempty" yaml:"permanent,omitempty" export:"true"`
}

func (b *Scheme) Schema() (string, error) {

	return "", nil
}

// NewRedirectScheme creates a new RedirectScheme middleware.
func NewRedirectScheme(ctx context.Context, next http.Handler, conf *Scheme, name string) (http.Handler, error) {
	logger := logs.GetLogger(ctx, name, TypeSchemeName)
	logger.Debug().Msg("Creating middleware")
	logger.Debug().Msgf("Setting up redirection to %s %s", conf.Scheme, conf.Port)

	if len(conf.Scheme) == 0 {
		return nil, errors.New("you must provide a target scheme")
	}

	port := ""
	if len(conf.Port) > 0 && !(conf.Scheme == schemeHTTP && conf.Port == "80" || conf.Scheme == schemeHTTPS && conf.Port == "443") {
		port = ":" + conf.Port
	}

	return newRedirect(next, uriPattern, conf.Scheme+"://${2}"+port+"${4}", conf.Permanent, clientRequestURL, name)
}

func clientRequestURL(req *http.Request) string {
	scheme := schemeHTTP
	host, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		host = req.Host
	} else {
		port = ":" + port
	}
	uri := req.RequestURI

	if match := uriRegexp.FindStringSubmatch(req.RequestURI); len(match) > 0 {
		scheme = match[1]

		if len(match[2]) > 0 {
			host = match[2]
		}

		if len(match[3]) > 0 {
			port = match[3]
		}

		uri = match[4]
	}

	if req.TLS != nil {
		scheme = schemeHTTPS
	}

	if xProto := req.Header.Get(xForwardedProto); xProto != "" {
		// When the initial request is a connection upgrade request,
		// X-Forwarded-Proto header might have been set by a previous hop to ws(s),
		// even though the actual protocol used so far is HTTP(s).
		// Given that we're in a middleware that is only used in the context of HTTP(s) requests,
		// the only possible valid schemes are one of "http" or "https", so we convert back to them.
		switch {
		case strings.EqualFold(xProto, "ws"):
			scheme = schemeHTTP
		case strings.EqualFold(xProto, "wss"):
			scheme = schemeHTTPS
		default:
			scheme = xProto
		}
	}

	if scheme == schemeHTTP && port == ":80" || scheme == schemeHTTPS && port == ":443" {
		port = ""
	}

	return strings.Join([]string{scheme, "://", host, port, uri}, "")
}
