package main

import "fmt"

type Hub struct {
	Channels        map[string]Channel
	Clients         map[string]*Client
	Commands        chan Command
	DeRegistrations chan *Client
	Registrations   chan *Client
}

func (h *Hub) Info(message string) {
	fmt.Println(message)
}

func (h *Hub) Register(c *Client) {
	if _, exists := h.Clients[c.Username]; exists {
		c.Username = ""
		c.Connection.Write([]byte("ERR Username taken"))
	} else {
		h.Clients[c.Username] = c
		c.Connection.Write([]byte("OK\n"))
	}
}

func (h *Hub) Unregister(c *Client) {
	if _, exists := h.Clients[c.Username]; exists {
		delete(h.Channels, c.Username)

		for _, channel := range h.Channels {
			delete(channel.Clients, c)
		}
	}
}

func (h *Hub) JoinChannel(username string, channelname string) {
	if client, ok := h.Clients[username]; ok {
		if channel, ok := h.Channels[channelname]; ok {
			channel.Clients[client] = true
		} else {
			h.Channels[channelname] = Channel{
				Name:    channelname,
				Clients: make(map[*Client]bool),
			}
			h.Channels[channelname].Clients[client] = true
		}
	}
}

func (h *Hub) SendMessage(username string, recipient string, message []byte) {
	if sender, ok := h.Clients[username]; ok {
		switch recipient[0] {
		case '#':
			if channel, ok := h.Channels[recipient]; ok {
				if _, ok := channel.Clients[sender]; ok {
					channel.Broadcast(sender.Username, message)
				}
			}
		case '@':
			if user, ok := h.Clients[recipient]; ok {
				user.Connection.Write(append(message, '\n'))
			}
		}
	}
}

func (h *Hub) ListUsers(username string) {

}

func (h *Hub) ListChannels(username string) {

}

func (h *Hub) LeaveChannel(username string, channelname string) {

}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Registrations:
			h.Register(client)
		case client := <-h.DeRegistrations:
			h.Unregister(client)
		case command := <-h.Commands:
			switch command.Id {
			case JOIN:
				h.JoinChannel(command.Sender, command.Recipient)
			case LEAVE:
				h.LeaveChannel(command.Sender, command.Recipient)
			case MESSAGE:
				h.SendMessage(command.Sender, command.Recipient, command.Body)
			case USERS:
				h.ListUsers(command.Sender)
			case CHANNELS:
				h.ListChannels(command.Sender)
			default:
				h.Info("Command not recognized by hub")
			}
		}
	}
}
