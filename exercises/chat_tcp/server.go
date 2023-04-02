package main

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	rooms       sync.Map
	commands    chan Command
	stopChan    chan struct{}
	defaultRoom *Room
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}
	s.run()

	log.Printf("TCP/IP server starts on port 9090")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Some error occured during accepting connection")
			continue
		}

		go func() {
			c := s.newClient(conn)
			err = c.readInput()

			if err != nil {
				log.Printf("Connection has been interuppted: %v", err.Error())
				c.close()
			}
		}()
	}
}

func NewServer(newRooms []*Room) *Server {
	m := sync.Map{}
	for _, room := range newRooms {
		m.Store(room.name, room)
	}

	return &Server{
		commands:    make(chan Command),
		rooms:       m,
		defaultRoom: newRooms[0],
	}
}

func (s *Server) run() {
	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			case c := <-s.commands:
				switch c.cmd {
				case cmdJoin:
					s.changeRoom(c)
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
}

func (s *Server) changeRoom(cmd Command) {
	foundRoom, ok := s.rooms.Load(cmd.input)

	if !ok {
		cmd.client.sendError("Room has not been found")
		return
	}
	room := foundRoom.(*Room)

	oldRoom := cmd.client.room
	oldRoom.leave(cmd.client)
	room.join(cmd.client)
	oldRoom.broadcast("User #%v has left this chat\n", cmd.client.name)
}

func (s *Server) sendMessage(cmd Command) {
	c := cmd.client
	c.room.broadcast("%v: %v\n", cmd.client.name, cmd.input)
}

func (s *Server) quitChat(cmd Command) {

}

func (s *Server) nick(cmd Command) {
	var exists *Client

	s.rooms.Range(func(key, val interface{}) bool {
		room := val.(*Room)
		exists = room.findUser(cmd.input)

		return exists == nil
	})

	// if exists != nil {
	// 	cmd.client.sendError("Provided username #%v already exists in system. Please try different one\n", cmd.input)
	// 	return
	// }

	oldUserName := cmd.client.name
	cmd.client.setUsername(cmd.input)
	cmd.client.room.broadcast("User #%v has changed his username to #%v\n", oldUserName, cmd.input)
}

func (s *Server) newClient(conn net.Conn) *Client {
	c := &Client{name: "anonymous", connection: conn, commands: s.commands, stopChan: make(chan struct{})}
	s.defaultRoom.join(c)

	return c
}
