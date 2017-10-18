package main

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"math/rand"
)

var flags *pflag.FlagSet

const (
	defaultChannel = "FooChannel"
	chanTD         = "executive"
	chanBU         = "grassroots"
)

const (
	ACKNOWLEDGEMENT = "ACK"
	MESSAGE = "foobar"
)

func main() {
	flags = pflag.NewFlagSet("pubsub", pflag.ExitOnError)

	opts := new(ClientOptions)
	role := flags.StringP("role", "r", "sub",
		"client's role, one of 'pub' or 'sub'")
	flags.StringVarP(&opts.name, "name", "n", "",
		"unique string identifier")
	flags.StringVarP(&opts.host, "host", "h", "localhost:6379",
		"redis server HOST:PORT")
	flags.StringVarP(&opts.password, "password", "p", "",
		"redis server password")
	flags.IntVarP(&opts.db, "db", "d", 0,
		"redis server database (default 0)")
	//flags.StringVarP(&opts.channel, "channel", "c", defaultChannel,
	//	"redis channel name")

	flags.Parse(os.Args)

	if len(opts.name) == 0 {
		panic(errors.New("you must provide a name"))
	}

	node := NewNode(opts)

	switch *role {
	case "pub":
		//var nNodes int

		c := make(chan []byte)
		go node.Subscribe(chanBU, c)

		for {
			msgSent := randomMsg()
			node.Publish(chanTD, msgSent)
			//nNodes = node.Publish(chanTD, MESSAGE)
			//for nNodes == 0 {
			//	nNodes = node.Publish(chanTD, MESSAGE)
			//}

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
	case "sub":
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
		panic(errors.New("positional argument ROLE not found; must provide 'pub' or 'sub'"))
	default:
		panic(errors.Errorf("invalid role '%s': must be one of 'pub' or 'sub'"))
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
