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
			return err
		}
		command, err := NewCommand(userInput, c)
		if err != nil {
			log.Printf("Error:[#%v] Invalid command: %v\n", c.name, err)

			c.sendMessage("Your command was invalid please correct it." +
				"/join [xxx] - use this command change current room\n" +
				"/msg [xxx] - use this command send a message to other members of room\n" +
				"/quit [xxx] - use this command to quit chat\n" +
				"/nick [xxx] - use this command to change your nick from anonymous\n" +
				"/rooms [xxx] - use this command to change your nick from anonymous\n\n")
			continue
		}
		c.commands <- command
	}
}

func (c *Client) sendMessage(messagef string, msgArgs ...interface{}) {
	_, err := c.connection.Write([]byte(fmt.Sprintf(messagef, msgArgs...)))
	if err != nil {
		log.Printf("Error:[#%v] Message has not been sent: %v\n", c.name, err)
		c.close()
	}
}

func (c *Client) close() {
	c.connection.Close()
}
