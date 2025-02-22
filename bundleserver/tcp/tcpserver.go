package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func StartTcpServer(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("[TcpServer] failed to create listener, err:", err)
		return err
	}
	fmt.Println("[TcpServer] listening on ", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[TcpServer] failed to accept connection, err:", err)
			continue
		}

		go handleConnection(conn)
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("[TcpServer] failed to read data, err:", err)
			}
			return
		}
		fmt.Printf("[TcpServer] request: %s", bytes)

		conn.Write(bytes)
	}
}
