package main

import (
	"fmt"
	"strings"
)

type commandId int

const (
	cmdUsername commandId = iota
	cmdJoin
	cmdMessage
	cmdQuit
	cmdRooms
)

type Command struct {
	client *Client
	cmd    commandId
	input  string
}

func NewCommand(msg string, client *Client) (Command, error) {
	parsedMsg := strings.Trim(msg, "\r\n")
	args := strings.SplitN(parsedMsg, " ", 2)

	if args[0] != "/quit" && args[0] != "/rooms" && len(args) != 2 {
		return Command{}, fmt.Errorf("malformed command, input string was %v after conversion %v", msg, args)
	}
	unprocessedCmd := strings.ReplaceAll(args[0], " ", "")

	var c Command
	switch unprocessedCmd {
	case "/nick":
		c = Command{client: client, cmd: 0, input: args[1]}
	case "/join":
		c = Command{client: client, cmd: 1, input: args[1]}
	case "/msg":
		c = Command{client: client, cmd: 2, input: args[1]}
	case "/quit":
		c = Command{client: client, cmd: 3, input: ""}
	case "/rooms":
		c = Command{client: client, cmd: 4, input: ""}
	default:
		return Command{}, fmt.Errorf("there is no such command available. The parsed command was %v", unprocessedCmd)
	}

	return c, nil
}
