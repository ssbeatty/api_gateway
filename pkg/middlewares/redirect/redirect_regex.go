package redirect

import (
	"api_gateway/pkg/middlewares/logs"
	"context"
	"net/http"
	"strings"
)

const TypeRegexName = "RedirectRegex"

type Regex struct {
	// Regex defines the regex used to match and capture elements from the request URL.
	Regex string `json:"regex,omitempty"`
	// Replacement defines how to modify the URL to have the new target URL.
	Replacement string `json:"replacement,omitempty"`
	// Permanent defines whether the redirection is permanent (301).
	Permanent bool `json:"permanent,omitempty"`
}

// NewRedirectRegex creates a redirect middleware.
func NewRedirectRegex(ctx context.Context, next http.Handler, conf *Regex, name string) (http.Handler, error) {
	logger := logs.GetLogger(ctx, name, TypeRegexName)
	logger.Debug().Msg("Creating middleware")
	logger.Debug().Msgf("Setting up redirection from %s to %s", conf.Regex, conf.Replacement)

	return newRedirect(next, conf.Regex, conf.Replacement, conf.Permanent, rawURL, name)
}

func (b *Regex) Schema() (string, error) {

	return "", nil
}

func rawURL(req *http.Request) string {
	scheme := schemeHTTP
	host := req.Host
	port := ""
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

	return strings.Join([]string{scheme, "://", host, port, uri}, "")
}
