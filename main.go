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
		"eat the rich",
	}
)

func init() {
	flags = pflag.NewFlagSet("pubsub", pflag.PanicOnError)

	opts := new(options)
	flags.StringVarP(&opts.host, "host", "h", "localhost:6379", "redis host")
	flags.StringVarP(&opts.password, "password", "p", "", "redis password")
	flags.IntVarP(&opts.db, "db", "d", 0, "redis database")
	flags.StringVarP(&opts.channel, "channel", "c", defaultChannel, "channel to publish to")
	flags.StringSliceVarP(
		&opts.messages, "messages", "m", defaultMessages, "message pool to be published from")
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
