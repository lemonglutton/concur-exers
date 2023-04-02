package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	name       string
	connection net.Conn
	commands   chan<- Command
	stopChan   <-chan struct{}
	room       *Room
}

func (c *Client) setUsername(username string) {
	c.name = username
}

func (c *Client) setRoom(room *Room) {
	c.room = room
}

func (c *Client) readInput() error {
	for {
		userInput, err := bufio.NewReader(c.connection).ReadString('\n')

		if err != nil {
			log.Printf("Something went wrong during reading message: %v", err.Error())
			return err
		}
		command, err := NewCommand(userInput, c)
		if err != nil {
			log.Printf("Invalid command: %v", err.Error())
			c.sendError("Invalid command: %v\n", err.Error())
			continue
		}

		select {
		case c.commands <- command:
		case <-c.stopChan:
			return nil

		}
	}
}

func (c *Client) sendError(messagef string, msgArgs ...interface{}) {
	_, err := c.connection.Write([]byte(fmt.Sprintf(messagef, msgArgs...)))
	if err != nil {
		c.close()
	}

}

func (c *Client) sendMessage(messagef string, msgArgs ...interface{}) {
	_, err := c.connection.Write([]byte(fmt.Sprintf(messagef, msgArgs...)))
	if err != nil {
		c.close()
	}
}

func (c *Client) close() {
	c.connection.Close()
}
