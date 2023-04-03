package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

const (
	host = "localhost"
	port = "9090"
)

type Server struct {
	rooms       sync.Map
	commands    chan Command
	defaultRoom *Room
}

func (s *Server) Listen() {
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	s.run()

	log.Printf("TCP/IP server starts on port %v", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error: Server could accept incoming connection %v\n", err)
			continue
		}

		go func() {
			c := s.newClient(conn)
			err = c.readInput()

			if err != nil {
				log.Printf("Error:[#%v] Couldn't recieve message from user. Probably connection was interuppted %v\n", c.name, err)
				s.quitChat(c)
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
		for c := range s.commands {
			switch c.cmd {
			case cmdJoin:
				s.changeRoom(c)
			case cmdMessage:
				s.sendMessage(c)
			case cmdQuit:
				s.quitChat(c.client)
			case cmdUsername:
				s.nick(c)
			case cmdRooms:
				s.listRooms(c)
			}
		}
	}()
}

func (s *Server) changeRoom(cmd Command) {
	foundRoom, ok := s.rooms.Load(cmd.input)

	if !ok {
		cmd.client.sendMessage("Room has not been found. Please use /rooms command to check what rooms are available\n")
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

func (s *Server) quitChat(c *Client) {
	defer c.close()

	oldRoom := c.room
	oldRoom.leave(c)
	oldRoom.broadcast("User #%v has left chat\n", c.name)
}

func (s *Server) listRooms(cmd Command) {
	var roomsNames []string

	s.rooms.Range(func(key, val interface{}) bool {
		room := val.(*Room)
		roomsNames = append(roomsNames, room.name)
		return true
	})
	cmd.client.sendMessage("This is a list of rooms: %v\n", strings.Join(roomsNames, ", "))
}

func (s *Server) nick(cmd Command) {
	var exists *Client

	s.rooms.Range(func(key, val interface{}) bool {
		room := val.(*Room)
		exists = room.findUser(cmd.input)

		return exists == nil
	})

	if exists != nil {
		cmd.client.sendMessage("Provided username #%v already exists in system. Please try different one\n", cmd.input)
		return
	}

	oldUserName := cmd.client.name
	cmd.client.setUsername(cmd.input)
	cmd.client.room.broadcast("User #%v has changed his username to #%v\n", oldUserName, cmd.input)
}

func (s *Server) newClient(conn net.Conn) *Client {
	c := &Client{name: "anonymous", connection: conn, commands: s.commands}
	s.defaultRoom.join(c)

	return c
}
