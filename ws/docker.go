package ws

import (
	"bufio"
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

// TestDockerService init a connection with docker service
func TestDockerService(Type string) {
	// Set up websocket connection
	dockerAddr := os.Getenv("DOCKER_ADDRESS")
	dockerPort := os.Getenv("DOCKER_PORT")
	if len(dockerAddr) == 0 {
		dockerAddr = "localhost"
	}
	if len(dockerPort) == 0 {
		dockerPort = "8888"
	}
	dockerPort = ":" + dockerPort
	dockerAddr = dockerAddr + dockerPort
	url := url.URL{Scheme: "ws", Host: dockerAddr, Path: "/" + Type}
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
			fmt.Print(string(msg))
		}
	}()
	
	go ReadMessage(conn, stopChan)

	switch Type {
	case "debug":
		command := Command{
			UserName:    "golang",
			ProjectName: "test",
			Type:        "debug",
		}
		conn.WriteJSON(&command)
	case "tty":
		command := Command{
			UserName:    "golang",
			ProjectName: "test",
			Type:        "tty",
			ENV:         []string{"GOPATH:/root/go:/home/go"},
			Command:     "go run main.go",
			PWD:         "",
		}
		conn.WriteJSON(&command)
	}
	<-stopChan
}

// ReadMessage read message from stdin
func ReadMessage(conn *websocket.Conn, stopChan chan<- bool) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			stopChan <- true
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
