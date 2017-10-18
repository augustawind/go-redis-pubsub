package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

const (
	chanTD = "director-agent"
	chanBU = "agent-director"
)

const (
	ACKNOWLEDGEMENT = "ACK"
)

type Role string

const (
	RoleDirector Role = "director"
	RoleAgent    Role = "agent"
)

func InvalidRole(msg string) error {
	return errors.Errorf("%s: must be one of 'director' or 'agent'", msg)
}

func run(role Role, opts *ClientOptions) {
	node := NewNode(opts)

	switch role {
	case RoleDirector:
		c := make(chan []byte)
		go node.Subscribe(chanBU, c)

		for {
			msgSent := randomMsg()
			node.Publish(chanTD, msgSent)

			msgReceived := <-c
			switch string(msgReceived) {
			case ACKNOWLEDGEMENT:
				node.log.Println("Acknowledged!")
			case msgSent:
				node.log.Println("Confirmed! Messages match.")
			default:
				node.log.Printf("Uh, oh! Unexpected response '%s'.", msgReceived)
			}
		}
	case RoleAgent:
		c := make(chan []byte)
		go node.Subscribe(chanTD, c)

		node.Publish(chanBU, ACKNOWLEDGEMENT)

		for {
			node.log.Printf("Waiting for message...")
			start := time.Now().Unix()

			msg := <-c
			end := time.Now().Unix()
			node.log.Printf("Received message in %d seconds.", end-start)

			node.log.Printf("Confirming with sender...")
			for node.Publish(chanBU, string(msg)) == 0 {
				time.Sleep(100 * time.Millisecond)
			}
		}
	case "":
		panic(InvalidRole("empty role"))
	default:
		panic(InvalidRole(fmt.Sprintf("invalid role '%s", role)))
	}

}

func randomMsg() string {
	i := rand.Intn(len(bagOfMessages))
	return bagOfMessages[i]
}

var bagOfMessages = []string{
	"there is no ethical consumption under late capitalism",
	"the history of society is the history of class struggle",
	"profit is wage theft",
	"a spectre is haunting europe",
	"WORKING MEN OF ALL COUNTRIES, UNITE!",
	"you have nothing to lose but your chains",
	"periodic economic crisis is unavoidable in a capitalist economy",
	"the streets will run red with the blood of the bourgeoisie",
	"eat the rich",
}
