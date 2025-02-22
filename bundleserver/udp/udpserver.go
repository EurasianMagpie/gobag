package udp

import (
	"fmt"
	"net"
)

func StartUdpServer(addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		fmt.Println("[UdpServer] failed to resolve udp addr, err:", err)
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("[UdpServer] failed to listen udp, err:", err)
		return err
	}
	fmt.Println("[UdpServer] listening on ", udpAddr)

	for {
		var buf [512]byte
		n, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println("[UdpServer] failed to read from udp, err:", err)
			return err
		}

		fmt.Println("[UdpServer] read n", n)
		fmt.Println("[UdpServer] request:", string(buf[0:n]))

		conn.WriteToUDP(buf[0:n], addr)
	}

	return nil
}
