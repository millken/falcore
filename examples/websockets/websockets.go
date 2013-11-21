package main

import (
	"fmt"
	"github.com/fitstar/falcore"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile|log.LstdFlags)
	port := 7575
	// setup pipeline
	pipeline := falcore.NewPipeline()

	// upstream
	pipeline.Upstream.PushBack(helloFilter)

	// setup server
	server := falcore.NewServer(port, pipeline)

	server.WebsocketHandler = WebsocketHandler
	// start the server
	// this is normally blocking forever unless you send lifecycle commands
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Could not start server:", err)
	}	
}

func WebsocketHandler(req *falcore.Request, ws *websocket.Conn) {
	tc := time.Tick(4*time.Second)
	go WebsocketReader(ws)
	for  {
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello from websocket!"))
		if err != nil {
			falcore.Error("send err: %v\n", err)
			return
		}
		<- tc
	}
}

func WebsocketReader(ws *websocket.Conn) {
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			falcore.Error("read err: %v\n", err)
			return
		}
		falcore.Info("Got Message: %v\n", string(data))
	}
}

var helloFilter = falcore.NewRequestFilter(func(req *falcore.Request) *http.Response {
	return falcore.StringResponse(req.HttpRequest, 200, nil, page)
})
