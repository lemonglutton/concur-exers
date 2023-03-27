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
	name       string
	connection net.Conn
}

type Server struct {
	users sync.Map
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}
	defer l.Close()

	log.Printf("TCP/IP server starts on port 9090")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Some erroc occured during accepting connection")
			return
		}

		go func() {
			client := s.createUser(conn)
			s.handleUserConnection(client)
		}()
	}
}

func NewServer() *Server {
	return &Server{users: sync.Map{}}
}

func (s *Server) run() {

}

func (s *Server) handleUserConnection(c *Client) {
	for {
		userInput, err := bufio.NewReader(c.connection).ReadString('\n')
		if err != nil {
			log.Printf("Connection was interuppted: %v", err.Error())
		}
		log.Println("Otrzymalem wiadomosc")

		s.users.Range(func(key, value interface{}) bool {
			c.connection.Write([]byte(fmt.Sprintf("\n%v: %v\n", c.name, userInput)))
			return true
		})
	}
}

func (s *Server) createUser(conn net.Conn) *Client {
	for {
		userInput, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("error occured during creating user")
		}

		username := strings.Trim(userInput, "\n")
		if len(username) <= 8 {
			conn.Write([]byte("Username was too short please provide it once again\n"))
			continue
		}

		var busyName bool
		s.users.Range(func(username, connection interface{}) bool {
			if username == userInput {
				conn.Write([]byte("This username has been already taken. Please choose different one"))

				busyName = false
				return busyName
			}
			busyName = true
			return busyName
		})

		if !busyName {
			defer conn.Write([]byte("Successful user creation! Welcome on board! \n"))
			c := Client{name: username, connection: conn}
			s.users.Store(userInput, c)
			return &c
		}
	}
}
