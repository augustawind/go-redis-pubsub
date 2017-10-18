package main

import (
	"os"
	"time"

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

func init() {
	flags = pflag.NewFlagSet("pubsub", pflag.ExitOnError)

	opts := new(options)
	flags.StringVarP(&opts.host, "host", "h", "localhost:6379",
		"redis server HOST:PORT")
	flags.StringVarP(&opts.password, "password", "p", "",
		"redis server password")
	flags.IntVarP(&opts.db, "db", "d", 0,
		"redis server database (default 0)")
	flags.StringVarP(&opts.channel, "channel", "c", defaultChannel,
		"redis channel name")
	flags.StringSliceVarP(&opts.messages, "messages", "m", nil,
		"messages to publish, chosen randomly (default messages included)")
	flags.Parse(os.Args)

	pub = newPub(opts)
	sub = newSub(opts)
}

func main() {
	go pub.publish()
	time.Sleep(500 * time.Millisecond)
	go sub.subscribe()
	for {
	}
}
