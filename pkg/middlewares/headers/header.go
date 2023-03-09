package headers

import (
	"api_gateway/pkg/middlewares/base"
	"api_gateway/pkg/middlewares/logs"
	"net/http"
	"strings"
)

const (
	TypeName = "Headers"
)

// Header is a middleware that helps setup a few basic security features.
// A single headerOptions struct can be provided to configure which features should be enabled,
// and the ability to override a few of the default values.
type Header struct {
	next             http.Handler
	hasCustomHeaders bool
	headers          *Headers
	name             string
}

type Headers struct {
	base.Config
	CustomRequestHeaders  map[string]string `json:"customRequestHeaders,omitempty"`
	CustomResponseHeaders map[string]string `json:"customResponseHeaders,omitempty"`
}

// NewHeader constructs a new header instance from supplied frontend header struct.
func NewHeader(next http.Handler, cfg *Headers, name string) (*Header, error) {
	return &Header{
		next:    next,
		headers: cfg,
		name:    name,
	}, nil
}

func (s *Header) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if s.hasCustomHeaders {
		s.modifyCustomRequestHeaders(req)
	}

	// If there is a next, call it.
	if s.next != nil {
		s.next.ServeHTTP(newResponseModifier(rw, req, s.PostRequestModifyResponseHeaders), req)
	}
}

// modifyCustomRequestHeaders sets or deletes custom request headers.
func (s *Header) modifyCustomRequestHeaders(req *http.Request) {
	logger := logs.GetLogger(req.Context(), s.name, TypeName)
	// Loop through Custom request headers
	for header, value := range s.headers.CustomRequestHeaders {
		switch {
		case value == "":
			req.Header.Del(header)

		case strings.EqualFold(header, "Host"):
			req.Host = value

		default:
			req.Header.Set(header, value)
			logger.Debug().Msgf("modify old header: %s: %s", header, value)
		}
	}
}

// PostRequestModifyResponseHeaders set or delete response headers.
// This method is called AFTER the response is generated from the backend
// and can merge/override headers from the backend response.
func (s *Header) PostRequestModifyResponseHeaders(res *http.Response) error {
	// Loop through Custom response headers
	for header, value := range s.headers.CustomResponseHeaders {
		if value == "" {
			res.Header.Del(header)
		} else {
			res.Header.Set(header, value)
		}
	}

	varyHeader := res.Header.Get("Vary")
	if varyHeader == "Origin" {
		return nil
	}

	if varyHeader != "" {
		varyHeader += ","
	}
	varyHeader += "Origin"

	res.Header.Set("Vary", varyHeader)
	return nil
}
