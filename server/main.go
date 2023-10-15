package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

type User struct {
	Username string
	Conn     *websocket.Conn
}

func newConnection(activeUsers *[]User) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.QueryParam("username")

		transmitInfoMessage(activeUsers, []byte(username+" connected."))

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		curUser := User{username, ws}
		*activeUsers = append(*activeUsers, curUser)
		defer curUser.Conn.Close()
		for {
			//Read
			_, userMsg, err := curUser.Conn.ReadMessage()
			if err != nil {
				fmt.Println("Connection closed")
				transmitInfoMessage(activeUsers, []byte(curUser.Username+" left."))
				break
			}
			fmt.Println(string(userMsg))
			//Pass message to other users
			transmitMessage(&curUser, activeUsers, userMsg)
		}

		return nil
	}
}

func transmitInfoMessage(activeUsers *[]User, msg []byte) {
	for _, user := range *activeUsers {
		user.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func transmitMessage(currentUser *User, activeUsers *[]User, msg []byte) {
	resMsg := []byte(currentUser.Username + ": ")
	resMsg = append(resMsg, msg...)
	for _, user := range *activeUsers {
		if user.Username != currentUser.Username {
			err := user.Conn.WriteMessage(websocket.TextMessage, resMsg)
			if err != nil {
				fmt.Println("Unable to send message: ", err)
			}
		}
	}
}

func main() {
	server := echo.New()

	activeUsers := make([]User, 0)
	//defining routes
	server.GET("/join", newConnection(&activeUsers))

	server.Start("localhost:8080")
}
