package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

type Client struct {
	Username string

	Connection net.Conn
	Outbound   chan<- *Command
	Register   chan<- *Client
	Unregister chan<- *Client
}

func (c *Client) RegisterNew(args []byte) error {
	username := bytes.TrimSpace(args)

	if username[0] != '@' {
		return fmt.Errorf("Username must begin with @")
	}

	if len(username) == 0 {
		return fmt.Errorf("Username cannot be blank")
	}

	c.Username = string(username)
	c.Register <- c

	return nil
}

func (c *Client) Error(error error) {
	c.Connection.Write([]byte("ERR " + error.Error() + "\n"))
}

func (c *Client) Handle(message []byte) {
	command := bytes.ToUpper(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))
	args := bytes.TrimSpace(bytes.TrimPrefix(message, command))

	switch string(command) {
	case "REG":
		if error := c.RegisterNew(args); error != nil {
			c.Error(error)
		}
	case "JOIN":
		if error := c.Join(args); error != nil {
			c.Error(error)
		}
	case "LEAVE":
		if error := c.Leave(args); error != nil {
			c.Error(error)
		}
	case "MSG":
		if error := c.Message(args); error != nil {
			c.Error(error)
		}
	case "CHNS":
		c.Channels()
	case "USRS":
		c.Users()
	default:
		c.Error(fmt.Errorf("Command not recognized %s", command))
	}
}

func (c *Client) Read() error {
	for {
		message, error := bufio.NewReader(c.Connection).ReadBytes('\n')
		if error == io.EOF {
			c.Unregister <- c
			return nil
		}

		if error != nil {
			return error
		}

		c.Handle(message)
	}
}
