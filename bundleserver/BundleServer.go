package main

import (
	"bundleserver/tcp"
	"bundleserver/udp"
	"bundleserver/web"
	"context"
	"flag"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var webaddr = flag.String("ws", "[::]:7770", "web websocket service address")
var tcpaddr = flag.String("tcp", "[::]:7771", "tcp service address")
var udpaddr = flag.String("udp", "[::]:7772", "udp service address")

func StartBundleServer() error {
	ctx := context.Background()
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return web.StartWebServer(*webaddr)
	})
	group.Go(func() error {
		return tcp.StartTcpServer(*tcpaddr)
	})
	group.Go(func() error {
		return udp.StartUdpServer(*udpaddr)
	})

	var groupErr error = nil
	if err := group.Wait(); err != nil {
		fmt.Printf("[bundleServer] Error: %v\n", err)
		groupErr = err
	} else {
		fmt.Println("[bundleServer] All server successfully")
	}

	return groupErr
}

func main() {
	flag.Parse()
	_ = StartBundleServer()
}
