package tcp

import (
	"context"
	"net"

	"github.com/libs4go/errors"
	"github.com/overlaynetwork/onet-go"
)

type tcpTransport struct{}

func (transport *tcpTransport) String() string {
	return transport.Protocol()
}

func (transport *tcpTransport) Protocol() string {
	return "tcp"
}

func (transport *tcpTransport) Listen(network *onet.OverlayNetwork) (onet.Listener, error) {

	tcpAddr, err := network.NavtiveAddr.ResolveNetAddr()

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport listen on %s error", network.NavtiveAddr)
	}

	listen, err := net.Listen(tcpAddr.Network(), tcpAddr.String())

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport listen on %s error", network.NavtiveAddr)
	}

	return onet.ToOnetListener(listen, network)
}

func (transport *tcpTransport) Dial(ctx context.Context, network *onet.OverlayNetwork) (onet.Conn, error) {

	tcpAddr, err := network.NavtiveAddr.ResolveNetAddr()

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport listen on %s error", network.NavtiveAddr)
	}

	conn, err := net.Dial(tcpAddr.Network(), tcpAddr.String())

	if err != nil {
		return nil, errors.Wrap(err, "tcp transport conn to %s error", network.NavtiveAddr)
	}

	return onet.ToOnetConn(conn, network)
}

func init() {
	if err := onet.RegisterTransport(&tcpTransport{}); err != nil {
		panic(err)
	}
}
