package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func ServerListener(wsConn *websocket.Conn, exit chan struct{}) {
	for {
		_, msg, err := wsConn.ReadMessage() //Reading messages from server

		if err != nil {
			fmt.Println("Connection closed")
			break
		}
		fmt.Println(string(msg))
	}

	close(exit)
}

func ServerWriter(wsConn *websocket.Conn, exit chan struct{}) {
	reader := bufio.NewReader(os.Stdin)
	for {
		userMsg, err := reader.ReadString('\r') //Reading user's input from console
		if err != nil {
			fmt.Println("Unable to read message: ", err)
		}
		userMsg = strings.Trim(userMsg, "\r\n") //Cleaning user's message
		if userMsg == "/leave" {
			fmt.Println("You left the chat")
			close(exit)
			break
		} else if userMsg != "" {
			// jsonMsg, _ := json.Marshal(userMsg)
			wsConn.WriteMessage(websocket.TextMessage, []byte(userMsg))
		}
	}
}

func initialMenu() string {
	fmt.Println("1. Connect to server\n2. Exit")

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}

		if input == '1' {
			break
		} else if input == '2' {
			os.Exit(0)
		} else {
			fmt.Println("Invalid input")
		}
	}

	fmt.Print("Enter your username: ")
	username := strings.Trim(func(bufR *bufio.Reader) string {
		for {
			u, _ := bufR.ReadString('\n')
			if u != "\r\n" {
				return u
			}
		}
	}(reader), "\r\n")

	return username
}

func main() {
	username := initialMenu()

	wsConn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/join?username="+username, nil)
	if err != nil {
		log.Fatal("Unable to connect: ", err)
	}

	exit := make(chan struct{})

	go func() {
		ServerListener(wsConn, exit)
	}()

	ServerWriter(wsConn, exit)

	wsConn.Close()

}
