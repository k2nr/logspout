package tls

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/gliderlabs/logspout/adapters/raw"
	"github.com/gliderlabs/logspout/router"
)

const (
	writeBuffer = 1024 * 1024
)

func init() {
	router.AdapterTransports.Register(new(tlsTransport), "tls")
	// convenience adapters around raw adapter
	router.AdapterFactories.Register(rawTLSAdapter, "tls")
}

func rawTLSAdapter(route *router.Route) (router.LogAdapter, error) {
	route.Adapter = "raw+tls"
	return raw.NewRawAdapter(route)
}

type tlsTransport int

func (_ *tlsTransport) Dial(addr string, options map[string]string) (net.Conn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Println("tls:", err)
		return nil, err
	}
	tcp_conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Println("tls:", err)
		return nil, err
	}
	err = tcp_conn.SetWriteBuffer(writeBuffer)
	if err != nil {
		log.Println("tls:", err)
		return nil, err
	}

	conn := tls.Client(tcp_conn, &tls.Config{ServerName: addr})
	return conn, nil
}
