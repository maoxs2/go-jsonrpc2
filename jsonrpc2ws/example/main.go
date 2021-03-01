package main

import (
	"github.com/c0mm4nd/go-jsonrpc2/jsonrpc2ws"
	"github.com/gorilla/websocket"
	"log"
	"time"

	"github.com/c0mm4nd/go-jsonrpc2"
)

type MyJsonHandler struct {
}

func (h *MyJsonHandler) Handle(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
	//result, _ := json.Marshal(map[string]interface{}{"ok": true})
	return jsonrpc2.NewJsonRpcSuccess(msg.ID, nil) // never use []byte{}
}

func main() {
	server, _ := jsonrpc2ws.NewServer(jsonrpc2ws.ServerConfig{
		Addr: "127.0.0.1:8888",
	})

	server.RegisterJsonRpcHandler("check", new(MyJsonHandler))
	server.RegisterJsonRpcHandleFunc("checkAgain", func(msg *jsonrpc2.JsonRpcMessage) *jsonrpc2.JsonRpcMessage {
		//result, _ := json.Marshal(map[string]interface{}{"ok": true})
		return jsonrpc2.NewJsonRpcSuccess(msg.ID, nil)
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	client, _ := jsonrpc2ws.NewClient(jsonrpc2ws.ClientConfig{
		Addr: "127.0.0.1:8888",
		Path: "/",
	})

	du := time.Tick(10 * time.Second)
	for {
		select {
		case <-du:
			msgType := websocket.TextMessage

			msg := jsonrpc2.NewJsonRpcRequest(1, "hello", nil)
			err := client.WriteMessage(msgType, msg)
			if err != nil {
				panic(err)
			}

			_, msg, err = client.ReadMessage()
			if err != nil {
				panic(err)
			}

			log.Println("reply:", msg)
		}
	}
}
