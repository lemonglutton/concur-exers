package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	name        string
	usernameSet bool
	connection  net.Conn
	commands    chan<- Command
	stopChan    <-chan struct{}
}

func (c *Client) setUsername(username string) {
	c.name = username
}

func (c *Client) readInput() error {
	for {
		userInput, err := bufio.NewReader(c.connection).ReadString('\n')

		if err != nil {
			log.Printf("Connection was interuppted: %v", err.Error())
			return err
		}
		command, err := NewCommand(userInput, c)
		if err != nil {
			log.Printf("Invalid command: %v", err.Error())
			continue
		}

		for {
			select {
			case c.commands <- command:
			case <-c.stopChan:
				return nil

			}
		}
	}
}

func (c *Client) sendError(errMsg string) {
	_, err := c.connection.Write([]byte(fmt.Sprintf("Error occured: %v", errMsg)))
	if err != nil {
		c.close()
	}

}

func (c *Client) sendMessage(msg string) {
	_, err := c.connection.Write([]byte(fmt.Sprintf("%v", msg)))
	if err != nil {
		c.close()
	}
}

func (c *Client) close() {
	c.connection.Close()

}
