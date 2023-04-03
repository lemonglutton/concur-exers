package main

import (
	"fmt"
	"sync"
)

type Room struct {
	name    string
	members sync.Map
}

func NewRoom(name string) *Room {
	return &Room{name: name, members: sync.Map{}}
}

func (r *Room) findUser(searchedUsername string) *Client {
	var foundClient *Client

	r.members.Range(func(key, value interface{}) bool {
		client := value.(*Client)

		if client.name == searchedUsername {
			foundClient = client
			return false
		}
		return true
	})
	return foundClient
}

func (r *Room) broadcast(messagef string, messageArgs ...interface{}) {
	r.members.Range(func(key, val interface{}) bool {
		c, _ := val.(*Client)
		c.sendMessage(messagef, messageArgs...)
		return true
	})
}

func (r *Room) join(c *Client) {
	r.members.Store(c.connection, c)
	c.setRoom(r)
	c.sendMessage(fmt.Sprintf("You have joined %v room!\n\n"+
		"/join [xxx] - use this command change current room\n"+
		"/msg [xxx] - use this command send a message to other members of room\n"+
		"/quit [xxx] - use this command to quit chat\n"+
		"/nick [xxx] - use this command to change your nick from anonymous\n"+
		"/rooms [xxx] - use this command to change your nick from anonymous\n\n", r.name))
}

func (r *Room) leave(c *Client) {
	r.members.Delete(c.connection)
}
