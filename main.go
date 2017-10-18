package main

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var flags *pflag.FlagSet

var (
	defaultChannel  = "revolutionary-vanguard"
	defaultMessages = []string{
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
)

func main() {
	flags = pflag.NewFlagSet("pubsub", pflag.ExitOnError)

	opts := new(clientOptions)
	role := flags.StringP("role", "r", "sub",
		"client's role, one of 'pub' or 'sub'")
	flags.StringVarP(&opts.host, "host", "h", "localhost:6379",
		"redis server HOST:PORT")
	flags.StringVarP(&opts.password, "password", "p", "",
		"redis server password")
	flags.IntVarP(&opts.db, "db", "d", 0,
		"redis server database (default 0)")
	flags.StringVarP(&opts.channel, "channel", "c", defaultChannel,
		"redis channel name")

	pubPrompt := flags.BoolP("prompt", "i", false,
		"wait for user input before each publish")
	pubMessages := flags.StringSliceP("messages", "m", nil,
		"messages to publish, chosen randomly (default messages included)")

	flags.Parse(os.Args)

	switch *role {
	case "pub":
		pub := newPub(opts, *pubPrompt, *pubMessages)
		c := make(chan int)
		go pub.publish(c)
		_ = <-c
		pub.interact()
	case "sub":
		sub := newSub(opts)
		c := make(chan bool)
		go sub.subscribe(c)
		sub.log.Printf("Waiting for message...")
		start := time.Now().Unix()
		_ = <-c
		end := time.Now().Unix()
		sub.log.Printf("Spent %d seconds waiting.", end-start)
	case "":
		panic(errors.New("positional argument ROLE not found; must provide 'pub' or 'sub'"))
	default:
		panic(errors.Errorf("invalid role '%s': must be one of 'pub' or 'sub'"))
	}
}
