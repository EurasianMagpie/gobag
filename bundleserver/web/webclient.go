package web

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func RequestWebSocketServer(addr string) (string, error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws/echo"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		//log.Println("dial failed, error:", err)
		return "", err
	}
	defer c.Close()

	req := time.Now().Format(time.RFC850)
	err = c.WriteMessage(websocket.TextMessage, []byte(req))
	if err != nil {
		//log.Println("write req, error", err)
		return "", err
	}
	_, resp, err := c.ReadMessage()
	if err != nil {
		//log.Println("read resp, error:", err)
		return "", err
	}
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		//log.Println("write close, error", err)
		return "", err
	}
	if string(resp) != req {
		//log.Println("wrong resp")
		return "", errors.New("wrong resp")
	}
	return "", nil
}

func RequestHttpServer(addr string) (string, error) {
	u := url.URL{Scheme: "http", Host: addr, Path: "/hi"}
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return "", nil
}
