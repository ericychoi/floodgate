package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ericychoi/floodgate"
)

func main() {
	var src, dst string
	var timeoutMS int64
	flag.StringVar(&src, "src", "", "source unix socket")
	flag.StringVar(&dst, "dst", "", "destination unix socket")
	flag.Int64Var(&timeoutMS, "timeout", 5000, "timeout in ms")
	flag.Parse()

	outChan := make(chan io.ReadCloser)
	defer close(outChan)
	fg := &floodgate.Proxy{
		Type:    "unix",
		SrcAddr: src,
		DstAddr: dst,
		Out:     outChan,
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	}
	err := fg.StartServer()
	if err != nil {
		log.Fatal(err)
	}
	defer fg.Close()
	go fg.Start()

	go func() {
		for {
			r := <-fg.Out
			br := bufio.NewReader(r)
			for {
				line, err := br.ReadBytes('\n')
				if err != nil {
					r.Close()
					break
				}
				fmt.Print(string(line))
			}
		}
	}()

	// Wait for TERM or INT
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigChan
	log.Printf("caught %v", sig)
}
