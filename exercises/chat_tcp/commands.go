package main

import (
	"errors"
	"strings"
)

type commandId int

const (
	cmdRooms commandId = iota
	cmdUsername
	cmdJoin
	cmdMessage
	cmdQuit
)

type Command struct {
	client *Client
	cmd    commandId
	input  string
}

func NewCommand(msg string, client *Client) (Command, error) {
	parsedMsg := strings.Trim(msg, "\r\n")
	args := strings.SplitN(parsedMsg, " ", 2)
	unprocessedCmd := strings.TrimSpace(args[0])

	var cmd commandId
	switch unprocessedCmd {
	case "/rooms":
		cmd = 0
	case "/msg":
		cmd = 2
	case "/join":
		cmd = 1
	case "/quit":
		cmd = 3
	case "/nick":
		cmd = 4
	default:
		return Command{}, errors.New("invalid command")
	}

	return Command{client: client, cmd: cmd, input: args[1]}, nil
}
