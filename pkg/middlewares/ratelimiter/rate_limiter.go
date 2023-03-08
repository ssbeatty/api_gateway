// Package ratelimiter implements a rate limiting and traffic shaping middleware with a set of token buckets.
package ratelimiter

import (
	"api_gateway/pkg/middlewares/logs"
	"context"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

const (
	TypeName = "RateLimiterType"
)

// rateLimiter implements rate limiting and traffic shaping with a set of token buckets;
// one for each traffic source. The same parameters are applied to all the buckets.
type rateLimiter struct {
	name string
	rate *IPRateLimiter[string]

	next http.Handler
}

type RateLimit struct {
	Average int `json:"average,omitempty"`
	// Every Second rate
	Period int `json:"period,omitempty"`
}

func (b *RateLimit) Schema() (string, error) {

	return "", nil
}

func getClientIP(req *http.Request) string {
	var clientIP string

	forwardedFor := req.Header.Get("X-Forwarded-For")
	realIP := req.Header.Get("X-Real-Ip")
	if realIP != "" || forwardedFor != "" {
		clientIP = strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
		if clientIP == "" {
			clientIP = strings.TrimSpace(realIP)
		}
	} else {
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
			return ip
		}
	}

	return clientIP
}

// New returns a rate limiter middleware.
func New(ctx context.Context, next http.Handler, config *RateLimit, name string) (http.Handler, error) {
	logger := logs.GetLogger(ctx, name, TypeName)
	logger.Debug().Msg("Creating middleware")

	limit := rate.Every(time.Duration(config.Period) * time.Second)
	return &rateLimiter{
		name: name,
		rate: NewIPRateLimiter[string](limit, config.Average),
		next: next,
	}, nil
}

func (rl *rateLimiter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := logs.GetLogger(req.Context(), rl.name, TypeName)
	ctx := logger.WithContext(req.Context())

	limiter := rl.rate.GetLimiter(getClientIP(req))

	if !limiter.Allow() {
		rl.serveDelayError(ctx, rw, time.Second)
		return
	}

	rl.next.ServeHTTP(rw, req)
}

func (rl *rateLimiter) serveDelayError(ctx context.Context, w http.ResponseWriter, delay time.Duration) {
	w.Header().Set("Retry-After", fmt.Sprintf("%.0f", math.Ceil(delay.Seconds())))
	w.Header().Set("X-Retry-In", delay.String())
	w.WriteHeader(http.StatusTooManyRequests)

	if _, err := w.Write([]byte(http.StatusText(http.StatusTooManyRequests))); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Could not serve 429")
	}
}
