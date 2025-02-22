package main

import (
	"bundleserver/tcp"
	"bundleserver/udp"
	"bundleserver/web"
	"flag"
	"fmt"
)

var host = flag.String("host", "204.79.197.200", "icmp host")
var webaddr = flag.String("ws", "[::]:7770", "web websocket service address")
var tcpaddr = flag.String("tcp", "[::]:7771", "tcp service address")
var udpaddr = flag.String("udp", "[::]:7772", "udp service address")

type testFunc func(string) (string, error)

type caseParam struct {
	name string
	fn   testFunc
	p1   string
}

func runCase(name string, fn testFunc, p1 string) {
	fmt.Printf("[ %s   ] -> %s\n", name, p1)
	r, err := fn(p1)
	if err != nil {
		fmt.Println("[   FAIL ] error, ", err)
	} else {
		fmt.Println("[     OK ] succeed, ", r)
	}
}

func runBundleClient() error {
	var cases = [...]caseParam{
		//{"ICMP", icmp.RunPinger, *host},
		{"TCP ", tcp.RequestTcpServer, *tcpaddr},
		{"UDP ", udp.RequestUdpServer, *udpaddr},
		{"HTTP", web.RequestHttpServer, *webaddr},
		{"WS  ", web.RequestWebSocketServer, *webaddr},
	}

	for i, c := range &cases {
		fmt.Printf(">>> case %d ----------\n", i+1)
		runCase(c.name, c.fn, c.p1)
	}
	return nil
}

func main() {
	flag.Parse()
	runBundleClient()
}
