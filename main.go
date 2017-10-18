package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var flags *pflag.FlagSet

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

	flags.Parse(os.Args)

	if len(opts.name) == 0 {
		panic(errors.New("you must provide a name"))
	}

	run(Role(*role), opts)
}
