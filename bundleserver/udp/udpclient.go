package udp

import (
	"bufio"
	"errors"
	"net"
	"time"
)

func RequestUdpServer(addr string) (string, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(time.Second * 10))

	req := time.Now().Format(time.RFC850)
	_, err = conn.Write([]byte(req + "\n"))
	if err != nil {
		return "", err
	} else {
		reader := bufio.NewReader(conn)
		b, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			return "", err
		}
		resp := string(b[:len(b)-1])
		if resp != req {
			return "", errors.New("wrong resp")
		}
	}

	return "", nil
}
