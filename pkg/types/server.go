package types

import (
	"api_gateway/pkg/tcp"
	"context"
	"fmt"
	"net"
)

type Stoppable interface {
	Shutdown(context.Context) error
	Close() error
}

type StoppableServer interface {
	Stoppable
	Serve(listener net.Listener) error
}

// writeCloser returns the given connection, augmented with the WriteCloser
// implementation, if any was found within the underlying conn.
func WriteCloser(conn net.Conn) (tcp.WriteCloser, error) {
	switch typedConn := conn.(type) {
	case *net.TCPConn:
		return typedConn, nil
	default:
		return nil, fmt.Errorf("unknown connection type %T", typedConn)
	}
}

type GrpcStop interface {
	Stop()
	Serve(listener net.Listener) error
}
