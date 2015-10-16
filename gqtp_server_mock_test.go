package goroo

import (
	"net"
)

type gqtpServer struct {
	Address  string
	Listener net.Listener
	handler  func(conn net.Conn)
	quit     chan bool
}

func (g *gqtpServer) start() {
	for {
		conn, err := g.Listener.Accept()
		if err != nil {
			return
		}
		go func() {
			g.handler(conn)
			g.quit <- true
		}()
	}
}

func (g *gqtpServer) Close() {
	<-g.quit
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

func newGqtpServer(handler func(conn net.Conn)) *gqtpServer {
	l := newGqtpLocalListener()
	gs := &gqtpServer{l.Addr().String(), l, handler, make(chan bool)}
	go gs.start()
	return gs
}
