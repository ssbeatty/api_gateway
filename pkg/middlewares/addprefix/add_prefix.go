package addprefix

import (
	"api_gateway/pkg/middlewares/base"
	"api_gateway/pkg/middlewares/logs"
	"context"
	"fmt"
	"net/http"
)

const (
	TypeName = "AddPrefix"
)

// AddPrefix is a middleware used to add prefix to an URL request.
type addPrefix struct {
	next   http.Handler
	prefix string
	name   string
}

type AddPrefix struct {
	base.Config
	// Prefix is the string to add before the current path in the requested URL.
	// It should include a leading slash (/).
	Prefix string `json:"prefix,omitempty"`
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, config *AddPrefix, name string) (http.Handler, error) {
	logs.GetLogger(ctx, name, TypeName).Debug().Msg("Creating middleware")
	var result *addPrefix

	if len(config.Prefix) > 0 {
		result = &addPrefix{
			prefix: config.Prefix,
			next:   next,
			name:   name,
		}
	} else {
		return nil, fmt.Errorf("prefix cannot be empty")
	}

	return result, nil
}

func (a *addPrefix) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := logs.GetLogger(req.Context(), a.name, TypeName)

	oldURLPath := req.URL.Path
	req.URL.Path = ensureLeadingSlash(a.prefix + req.URL.Path)
	logger.Debug().Msgf("URL.Path is now %s (was %s).", req.URL.Path, oldURLPath)

	if req.URL.RawPath != "" {
		oldURLRawPath := req.URL.RawPath
		req.URL.RawPath = ensureLeadingSlash(a.prefix + req.URL.RawPath)
		logger.Debug().Msgf("URL.RawPath is now %s (was %s).", req.URL.RawPath, oldURLRawPath)
	}
	req.RequestURI = req.URL.RequestURI()

	a.next.ServeHTTP(rw, req)
}

func ensureLeadingSlash(str string) string {
	if str == "" {
		return str
	}

	if str[0] == '/' {
		return str
	}

	return "/" + str
}
