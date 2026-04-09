package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Client struct {
	name string
	conn net.Conn
	out  chan []byte
}

type Room struct {
	clients map[*Client]bool
	mu      sync.Mutex
}

func NewRoom() *Room {
	return &Room{
		clients: make(map[*Client]bool),
	}
}

func (r *Room) Add(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = true
}

func (r *Room) Remove(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}

func (r *Room) Broadcast(msg []byte, sender *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.clients {
		if c == sender {
			continue
		}
		select {
		case c.out <- msg:
		default:
		}
	}
}

func writeLoop(c *Client) {
	for msg := range c.out {
		_, err := c.conn.Write(msg)
		if err != nil {
			return
		}
	}
}

func handleClient(conn net.Conn, room *Room) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	conn.Write([]byte("enter your name: "))
	nameLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name := strings.TrimSpace(nameLine)
	if name == "" {
		return
	}

	client := &Client{
		name: name,
		conn: conn,
		out:  make(chan []byte, 16),
	}

	room.Add(client)
	go writeLoop(client)

	log.Printf("%s joined from %s", name, conn.RemoteAddr())
	room.Broadcast([]byte(fmt.Sprintf("*** %s joined ***\n", name)), client)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}
		msg := []byte(fmt.Sprintf("[%s]: %s\n", name, text))
		room.Broadcast(msg, client)
	}

	room.Remove(client)
	close(client.out)
	room.Broadcast([]byte(fmt.Sprintf("*** %s left ***\n", name)), client)
	log.Printf("%s disconnected", name)
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}
	defer listener.Close()

	log.Println("chat server listening on :8080")

	room := NewRoom()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleClient(conn, room)
	}
}
