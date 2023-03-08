package inflightconn

import (
	"api_gateway/pkg/middlewares/logs"
	"api_gateway/pkg/tcp"
	"context"
	"fmt"
	"net"
	"sync"
)

const TypeName = "InFlightConnTCP"

type inFlightConn struct {
	name           string
	next           tcp.Handler
	maxConnections int64

	mu          sync.Mutex
	connections map[string]int64 // current number of connections by remote IP.
}

type TCPInFlightConn struct {
	// Amount defines the maximum amount of allowed simultaneous connections.
	// The middleware closes the connection if there are already amount connections opened.
	Amount int64 `json:"amount,omitempty"`
}

func (b *TCPInFlightConn) Schema() (string, error) {

	return "", nil
}

// New creates a max connections middleware.
// The connections are identified and grouped by remote IP.
func New(ctx context.Context, next tcp.Handler, config *TCPInFlightConn, name string) (tcp.Handler, error) {
	logger := logs.GetLogger(ctx, name, TypeName)
	logger.Debug().Msg("Creating middleware")

	return &inFlightConn{
		name:           name,
		next:           next,
		connections:    make(map[string]int64),
		maxConnections: config.Amount,
	}, nil
}

// ServeTCP serves the given TCP connection.
func (i *inFlightConn) ServeTCP(conn tcp.WriteCloser) {
	logger := logs.GetLogger(context.Background(), i.name, TypeName)

	ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		logger.Error().Err(err).Msg("Cannot parse IP from remote addr")
		conn.Close()
		return
	}

	if err = i.increment(ip); err != nil {
		logger.Error().Err(err).Msg("Connection rejected")
		conn.Close()
		return
	}

	defer i.decrement(ip)

	i.next.ServeTCP(conn)
}

// increment increases the counter for the number of connections tracked for the
// given IP.
// It returns an error if the counter would go above the max allowed number of
// connections.
func (i *inFlightConn) increment(ip string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.connections[ip] >= i.maxConnections {
		return fmt.Errorf("max number of connections reached for %s", ip)
	}

	i.connections[ip]++

	return nil
}

// decrement decreases the counter for the number of connections tracked for the
// given IP.
// It ensures that the counter does not go below zero.
func (i *inFlightConn) decrement(ip string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.connections[ip] <= 0 {
		return
	}

	i.connections[ip]--
}
