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

func (r *Room) findUser(searchedUsername string) (*Client, error) {
	foundClient, ok := r.members.Load(searchedUsername)

	if !ok {
		return nil, nil
	}

	client, ok := foundClient.(Client)
	if !ok {
		return nil, fmt.Errorf("%v client has not been found", searchedUsername)
	}

	return &client, nil
}

func (r *Room) broadcast(sender *Client, msg string) error {
	r.members.Range(func(key, val interface{}) bool {
		username, _ := key.(string)
		client, _ := val.(Client)

		if username != sender.name {
			client.sendMessage(msg)
		}

		return true
	})
	return nil
}

func (r *Room) join(c *Client) {
	r.members.Store(c.name, c)
}

func (r *Room) leave(c *Client) {
	r.members.Delete(c.name)
}
