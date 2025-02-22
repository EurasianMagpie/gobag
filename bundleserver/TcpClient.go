package main

import (
	"bundleserver/tcp"
	"flag"
)

var tcpaddr = flag.String("tcp", "[::]:7771", "tcp service address")

func main() {
	flag.Parse()
	tcp.RunTcpClient(*tcpaddr)
}
