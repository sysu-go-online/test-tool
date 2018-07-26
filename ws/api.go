package ws

import (
	"fmt"
	"net/url"
	"time"

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
		Command:     "set",
		BreakPoints: "main.cpp:6",
	}
	conn.WriteJSON(&clientmsg)

	timer := time.NewTimer(3 * time.Second)
	time.NewTimer(5 * time.Second)
	<-timer.C
	clientmsg.Command = "run"
	conn.WriteJSON(&clientmsg)

	timer = time.NewTimer(3 * time.Second)
	<-timer.C
	clientmsg.Command = "next"
	conn.WriteJSON(&clientmsg)

	timer = time.NewTimer(3 * time.Second)
	<-timer.C
	clientmsg.Command = "continue"
	conn.WriteJSON(&clientmsg)

	timer = time.NewTimer(3 * time.Second)
	<-timer.C
	clientmsg.Command = "quit"
	conn.WriteJSON(&clientmsg)
	<-stopChan
}
