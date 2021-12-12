package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Message struct {
	Discards [][]Discard `json:"discards"`
}

type Discard struct {
	Tile        string `json:"tile"`
	IsTsumogiri bool   `json:"isTsumogiri"`
	IsRiichi    bool   `json:"isRiichi"`
	IsRedFive   bool   `json:"isRedFive"`
}

func handleConnections(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var res string
		err := ws.ReadJSON(&res)
		if err != nil {
			fmt.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
	}
	return nil
}

func handleMessages() {
	fmt.Println("handle Messages is called")
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
