package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	wsServer      *WebsocketServer = nil
	wsServerMutex                  = sync.Mutex{}
)

type WebsocketServer struct {
	upgrader websocket.Upgrader
	wsEcho   *websocket.Conn
}

func NewWebsocketServer() *WebsocketServer {
	return &WebsocketServer{
		upgrader: websocket.Upgrader{},
		wsEcho:   nil,
	}
}

func (s *WebsocketServer) handleWsEcho(c *gin.Context) {
	//upgrade get request to websocket protocol
	ws, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	s.wsEcho = ws
	for {
		//Read Message from client
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		//If client message is ping will return pong
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//Response message to client
		err = ws.WriteMessage(mt, message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func (s *WebsocketServer) handleWsSendAny(c *gin.Context) {
	//upgrade get request to websocket protocol
	ws, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	if ws != nil {
		msg := strconv.FormatInt(time.Now().UnixMilli(), 10)
		c.JSON(200, gin.H{
			"now": msg,
		})
		eventParams := make(map[string]interface{}, 2)
		eventParams["now"] = msg
		jsonBytes, err := json.Marshal(eventParams)
		if err != nil {
			fmt.Println("[ginWsSendAny] err:", err)
			return
		}
		err = ws.WriteMessage(websocket.TextMessage, jsonBytes)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *WebsocketServer) handleHtmlHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"wsuri": "ws://" + c.Request.Host + "/ws/echo",
	})
}

func doDirectHandleWsEcho(c *gin.Context) {
	wsServer.handleWsEcho(c)
}

func doDirectHandleWsSendAny(c *gin.Context) {
	wsServer.handleWsSendAny(c)
}

func doDirectHandleHtmlHome(c *gin.Context) {
	wsServer.handleHtmlHome(c)
}

func StartWebServer(addr string) error {
	if wsServer == nil {
		wsServerMutex.Lock()
		if wsServer == nil {
			wsServer = NewWebsocketServer()
		}
		wsServerMutex.Unlock()
	}
	if wsServer == nil {
		log.Println("NewWebsocketServer failed")
		return errors.New("NewWebsocketServer failed")
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/home", doDirectHandleHtmlHome)
	r.GET("/ws/echo", doDirectHandleWsEcho)
	r.GET("/ws/send", doDirectHandleWsSendAny)
	r.GET("/hi", func(c *gin.Context) {
		//log.Printf("[xx] hi, %p\n", c.Request)
		c.JSON(http.StatusOK, gin.H{
			"message": "hi, there~",
		})
	})
	r.Run(addr)

	return nil
}
