package goroo

import (
	"io"
	"net"
)

type gqtpServer struct {
	Address  string
	Listener net.Listener
	handler  func(io.Writer, io.Reader)
}

func (g *gqtpServer) start() {
	for {
		conn, err := g.Listener.Accept()
		if err != nil {
			continue
		}
		g.handler(conn, conn)
		conn.Close()
	}
}

func (g *gqtpServer) Close() {
	g.Listener.Close()
}

func newGqtpLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(err)
		}
	}
	return l
}

func newGqtpServer(handler func(io.Writer, io.Reader)) *gqtpServer {
	l := newGqtpLocalListener()
	gs := &gqtpServer{l.Addr().String(), l, handler}
	go gs.start()
	return gs
}
