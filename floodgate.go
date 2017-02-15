package floodgate

import (
	"io"
	"log"
	"net"
	"path/filepath"
	"time"
)

type Proxy struct {
	Type    string
	SrcAddr string
	DstAddr string
	Out     chan io.ReadCloser
	Timeout time.Duration

	srcAbsoluteAddr string
	listener        net.Listener
}

type ProxyReadCloser struct {
	io.Reader
	src net.Conn
	dst net.Conn
}

func (p ProxyReadCloser) Close() error {
	p.src.Close()
	p.dst.Close()
	return nil
}

func (f *Proxy) StartServer() error {
	l, err := net.Listen("unix", f.SrcAddr)
	if err != nil {
		return err
	}

	absoluteSocketPath, err := filepath.Abs(f.SrcAddr)
	if err != nil {
		return err
	}

	log.Println("Listening on unix://" + absoluteSocketPath)
	f.listener = l
	f.srcAbsoluteAddr = absoluteSocketPath
	return nil
}

func (f *Proxy) Start() {
	if f.listener == nil {
		log.Print("nil listener: please call StartServer() first")
		return
	}

	for {
		conn, err := f.listener.Accept()
		if err != nil {
			log.Printf("couldn't accept: %s", err)
			return
		}

		log.Print("notifysock: client connected")

		go f.connectAndProxy(conn)
	}
}

func (f *Proxy) Close() {
	log.Printf("deleting socket at %s", f.srcAbsoluteAddr)
	f.listener.Close()
}

func (f *Proxy) connectAndProxy(srcConn net.Conn) {
	dstConn, err := net.DialTimeout("unix", f.DstAddr, f.Timeout)
	if err != nil {
		log.Fatalf("couldn't connect to %s: %s", f.DstAddr, err.Error())
	}
	r := ProxyReadCloser{
		Reader: io.TeeReader(srcConn, dstConn),
		src:    srcConn,
		dst:    dstConn,
	}
	f.Out <- r
}
