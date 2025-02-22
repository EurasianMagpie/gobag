package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"
)

func RequestTcpServer(addr string) (string, error) {
	d := net.Dialer{Timeout: time.Second * 5}
	conn, err := d.Dial("tcp", addr)
	//conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	req := time.Now().Format(time.RFC850)
	if _, err := conn.Write([]byte(req + "\n")); err != nil {
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

func RunTcpClient(addr string) (string, error) {
	d := net.Dialer{Timeout: time.Second * 5}
	conn, err := d.Dial("tcp", addr)
	//conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("[TcpClient] dail failed,", err)
		return "", err
	}
	defer conn.Close()

	for {
		req := time.Now().Format(time.RFC850)
		if _, err := conn.Write([]byte(req + "\n")); err != nil {
			fmt.Println("[TcpClient] write error,", err)
			return "", err
		} else {
			reader := bufio.NewReader(conn)
			b, err := reader.ReadBytes(byte('\n'))
			if err != nil {
				fmt.Println("[TcpClient] read error,", err)
				return "", err
			}
			resp := string(b[:len(b)-1])
			if resp != req {
				fmt.Println("[TcpClient] wrong resp")
				return "", errors.New("wrong resp")
			} else {
				fmt.Println("[TcpClient] recv:", resp)
			}
		}
		time.Sleep(time.Second * 1)
	}

	return "", nil
}
