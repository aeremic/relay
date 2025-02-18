package main

type Channel struct {
	Name    string
	Clients map[*Client]bool
}

func (c *Channel) Broadcast(source string, messageFromSource []byte) {
	messageToBroadcast := append([]byte(source), ": "...)
	messageToBroadcast = append(messageToBroadcast, messageFromSource...)
	messageToBroadcast = append(messageToBroadcast, '\n')

	for client := range c.Clients {
		client.Connection.Write(messageToBroadcast)
	}
}
