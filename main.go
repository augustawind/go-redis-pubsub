package main

import (
	"os"

	"github.com/spf13/pflag"
)

var flags = pflag.NewFlagSet("pubsub", pflag.ExitOnError)

func main() {
	opts := new(ClientOptions)

	flags.StringVarP(&opts.role, "role", "r", "agent",
		"client's role, one of 'director' or 'agent'")
	flags.StringVarP(&opts.host, "host", "h", "localhost:6379",
		"redis server HOST:PORT")
	flags.StringVarP(&opts.password, "password", "p", "",
		"redis server password")
	flags.IntVarP(&opts.db, "db", "d", 0,
		"redis server database (default 0)")

	flags.Parse(os.Args)

	run(opts)
}
