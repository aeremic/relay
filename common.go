package main

type ID int

const (
	REGISTER ID = iota
	JOIN
	LEAVE
	MESSAGE
	CHANNELS
	USERS
)

type Command struct {
	Id        ID
	Recipient string
	Sender    string
	Body      []byte
}
