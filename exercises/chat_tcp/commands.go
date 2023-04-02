package main

import (
	"errors"
	"strings"
)

type commandId int

const (
	cmdUsername commandId = iota
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

	if len(args) != 2 {
		return Command{}, errors.New("malformed command")
	}
	unprocessedCmd := strings.ReplaceAll(args[0], " ", "")

	var cmd commandId
	switch unprocessedCmd {
	case "/nick":
		cmd = 0
	case "/join":
		cmd = 1
	case "/msg":
		cmd = 2
	case "/quit":
		cmd = 3
	default:
		return Command{}, errors.New("invalid command")
	}

	return Command{client: client, cmd: cmd, input: args[1]}, nil
}
