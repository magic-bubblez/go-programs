package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to start listener: %v", err)
	}
	log.Println("server listening on port 8080")
	for {
		conn, err := listener.Accept() //waits unitil a client connects
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		log.Printf("new client connected: %s", conn.RemoteAddr())

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close() //closing connection to avoid reserving file descriptors
	defer log.Printf("client disconnected: %s", conn.RemoteAddr())

	// a buffer to read bytes into (like wc)
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("error from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		log.Printf("received from %s: %q\n", conn.RemoteAddr(), buf[:n])

		// echo it back to the client
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("write error to %s: %v", conn.RemoteAddr(), err)
			return
		}
	}
}
