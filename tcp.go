package tcp

import (
	"context"
	"net"
	"sync"

	"github.com/libs4go/errors"
	"github.com/overlaynetwork/onet-go"
)

type tcpTransport struct {
	sync.RWMutex
	listeners map[string]net.Listener
}

func newTCPTransport() *tcpTransport {
	return &tcpTransport{
		listeners: make(map[string]net.Listener),
	}
}

func (transport *tcpTransport) String() string {
	return transport.Protocol()
}

func (transport *tcpTransport) Protocol() string {
	return "tcp"
}

func (transport *tcpTransport) Server(ctx context.Context, network *onet.OverlayNetwork, addr *onet.Addr, next onet.Next) (onet.Conn, error) {

	tcpAddr, _, err := addr.ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	transport.RLock()
	listener, ok := transport.listeners[addr.String()]
	transport.RUnlock()

	if !ok {
		var err error
		listener, err = net.Listen(tcpAddr.Network(), tcpAddr.String())

		if err != nil {
			return nil, err
		}

		transport.Lock()
		transport.listeners[addr.String()] = listener
		transport.Unlock()
	}

	conn, err := listener.Accept()

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport listen on %s error", tcpAddr)
	}

	return onet.ToOnetConn(conn, network, addr)
}

func (transport *tcpTransport) Client(ctx context.Context, network *onet.OverlayNetwork, addr *onet.Addr, next onet.Next) (onet.Conn, error) {

	tcpAddr, _, err := addr.ResolveNetAddr()

	if err != nil {
		return nil, err
	}

	conn, err := net.Dial(tcpAddr.Network(), tcpAddr.String())

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport conn to %s error", tcpAddr)
	}

	return onet.ToOnetConn(conn, network, addr)
}

func (transport *tcpTransport) Close(network *onet.OverlayNetwork, addr *onet.Addr, next onet.NextClose) error {
	transport.Lock()
	delete(transport.listeners, addr.String())
	transport.Unlock()

	return next()

}

func init() {
	if err := onet.RegisterTransport(newTCPTransport()); err != nil {
		panic(err)
	}
}
