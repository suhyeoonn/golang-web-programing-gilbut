package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var clients []*Client

type Client struct {
	conn *websocket.Conn
	send chan *Message

	roomId string
	user   *User
}

const messageBufferSize = 256

func newClient(conn *websocket.Conn, roomId string, u *User) {
	c := &Client{
		conn:   conn,
		send:   make(chan *Message, messageBufferSize),
		roomId: roomId,
		user:   u,
	}

	clients = append(clients, c)

	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) Close() {
	for i, client := range clients {
		if client == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	close(c.send)

	c.conn.Close()
	log.Printf("close connection. addr : %s", c.conn.RemoteAddr())
}

func (c *Client) readLoop() {
	// 메시지 수신 대기
	for {
		m, err := c.read()
		if err != nil {
			log.Println("read message error:", err)
			break
		}

		m.create()
		broadcast(m)
	}
	c.Close()
}

func (c *Client) writeLoop() {
	// 클라이언트의 send 채널 메시지 수신 대기
	for msg := range c.send {
		if c.roomId == msg.RoomId.Hex() {
			fmt.Println("write Loop~~")
			c.write(msg)
		}
	}
}

func broadcast(m *Message) {
	for _, client := range clients {
		client.send <- m
	}
}

func (c *Client) read() (*Message, error) {
	var msg *Message

	if err := c.conn.ReadJSON(&msg); err != nil {
		return nil, err
	}

	msg.CreatedAt = time.Now()
	msg.User = c.user

	log.Println("read from websocket:", msg)

	return msg, nil
}

func (c *Client) write(m *Message) error {
	log.Println("write to websocket:", m)

	return c.conn.WriteJSON(m)
}
