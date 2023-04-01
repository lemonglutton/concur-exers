package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

const welcomeMsg = "Welcome to the simple TCP/IP chat. This is a list of commands You can use.\n" +
	"First of all to start using this chat You need to provide us with Your username.\n" + "Username must consist at least with 8 chars\n" +
	"Once Your username is set you start in you will join #general channel by default"

type Server struct {
	rooms    sync.Map
	commands <-chan Command
	stopChan chan struct{}
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}
	<-s.run()

	log.Printf("TCP/IP server starts on port 9090")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Some erroc occured during accepting connection")
			continue
		}

		go func() {
			c := s.newClient(conn)
			err = c.readInput()

			if err != nil {
				log.Printf("Connection has been interuppted: %v", err.Error())
			}
			c.close()
		}()
	}
}

func NewServer(newRooms []Room) *Server {
	m := sync.Map{}
	for _, room := range newRooms {
		m.Store(room.name, room)
	}

	return &Server{
		commands: make(<-chan Command),
		rooms:    m,
	}
}

func (s *Server) run() <-chan struct{} {
	goroutineIsReadyChan := make(chan struct{})

	go func() {
		close(goroutineIsReadyChan)

		for {
			select {
			case <-s.stopChan:
				return
			case c := <-s.commands:
				if c.cmd != cmdUsername && !c.client.usernameSet {
					c.client.sendError("You haven't set Your username. Try again with /nick command.")
				}
				switch c.cmd {
				case cmdRooms:
					s.listRooms(c)
				case cmdJoin:
					s.joinRoom(c)
				case cmdMessage:
					s.sendMessage(c)
				case cmdQuit:
					s.quitChat(c)
				case cmdUsername:
					s.nick(c)
				}
			}
		}
	}()

	return goroutineIsReadyChan
}

func (s *Server) listRooms(cmd Command) {
	var listOfRooms []string

	s.rooms.Range(func(key, _ interface{}) bool {
		room, _ := key.(Room)
		listOfRooms = append(listOfRooms, room.name)

		return true
	})
	cmd.client.sendMessage(strings.Join(listOfRooms, ","))
}

func (s *Server) joinRoom(cmd Command) {
	foundRoom, ok := s.rooms.Load(cmd.input)
	if !ok {
		cmd.client.sendError("")
		return
	}

	room, _ := foundRoom.(Room)
	room.join(cmd.client)
	cmd.client.sendMessage("You have joined a room!")
}

func (s *Server) sendMessage(cmd Command) {

	// err := c.room.broadcast()
}

func (s *Server) quitChat(cmd Command) {}

func (s *Server) nick(cmd Command) {
	var err error
	var exists *Client

	s.rooms.Range(func(key, _ interface{}) bool {
		room, _ := key.(Room)
		exists, err = room.findUser(cmd.input)

		return exists == nil
	})

	if err != nil {
		cmd.client.sendError("provided username already exists in system, try different one")
		return
	}
	cmd.client.setUsername(cmd.input)
}

func (s *Server) newClient(conn net.Conn) *Client {
	c := Client{name: "anonymous", connection: conn, stopChan: make(chan struct{}), room: nil}
	c.sendMessage(welcomeMsg)

	return &c
}
