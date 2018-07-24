package ws

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

// TestAPIDebug test debug service in api service
func TestAPIDebug() {
	addr := "localhost"
	port := ":8080"
	addr = addr + port
	url := url.URL{Scheme: "ws", Host: addr, Path: "/ws/debug"}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	stopChan := make(chan bool, 1)
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				stopChan <- true
				return
			}
			fmt.Println(string(msg))
		}
	}()

	clientmsg := ClientDebugMessage{
		Command: "run",
		BreakPoints: []string{"main.cpp:6"},
	}
	conn.WriteJSON(&clientmsg)
	// clientmsg.Command = "run"
	// conn.WriteJSON(&clientmsg)
	// clientmsg.Command = "next"
	// conn.WriteJSON(&clientmsg)
	// clientmsg.Command = "continue"
	// conn.WriteJSON(&clientmsg)
	clientmsg.Command = "finish"
	conn.WriteJSON(&clientmsg)
	<-stopChan
}
