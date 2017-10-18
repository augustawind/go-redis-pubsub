package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/spf13/pflag"
)

const (
	defaultChannel = "chaos"
	defaultMessage = "mind exploder"
)

var (
	host     string
	password string
	db       int
	channel  string
	message  string
	counter  int
	pubLog   *log.Logger
	subLog   *log.Logger
)

func newLogger(name string) *log.Logger {
	return log.New(os.Stderr, fmt.Sprintf("[%s] ", name), log.Ltime|log.Lmicroseconds)
}

func choice(options []string) string {
	i := rand.Intn(len(options))
	return options[i]
}

func init() {
	pubLog = newLogger("pub")
	subLog = newLogger("sub")

	pflag.StringVarP(&host, "host", "h", "localhost:6379", "redis host")
	pflag.StringVarP(&password, "password", "p", "", "redis password")
	pflag.IntVarP(&db, "db", "d", 0, "redis database")
	pflag.StringVarP(
		&channel, "channel", "c", defaultChannel, "channel to publish to")
	pflag.StringVarP(
		&message, "message", "m", defaultMessage, "message to publish")
	pflag.Parse()
}

func main() {
	pubConn := getRedisConn()
	subConn := getRedisConn()
	go publish(pubConn)
	time.Sleep(500 * time.Millisecond)
	go subscribe(subConn)
	for {
	}
}

func getRedisConn() redis.Conn {
	conn, err := redis.Dial(
		"tcp",
		host,
		redis.DialPassword(password),
		redis.DialDatabase(db),
	)
	check(err)
	return conn
}

func publish(conn redis.Conn) {
	for {
		msg := fmt.Sprintf("[%04d] %s", counter, message)
		n, err := redis.Int(conn.Do("PUBLISH", channel, msg))
		check(err)
		counter++

		msg = fmt.Sprintf("message received by (%d) recipient", n)
		if n != 1 {
			pubLog.Print(msg, "s")
		} else {
			pubLog.Print(msg)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func subscribe(conn redis.Conn) {
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(channel)

	for counter < 10 {
		switch v := psc.Receive().(type) {
		case redis.Message:
			subLog.Printf("channel: %s | message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			subLog.Printf("channel: %s | kind: %s | count: %d\n", v.Channel, v.Kind, v.Count)
		case error:
			panic(v)
		}
	}
	err := psc.Unsubscribe(channel)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
