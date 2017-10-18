package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	messages []string
	pub      pubsub
	sub      pubsub
)

type pubsub struct {
	*options
	role    role
	conn    redis.Conn
	counter int
	log     *log.Logger
}

type options struct {
	host     string
	password string
	db       int
	channel  string
	messages []string
}

type role string

const (
	isPub role = "pub"
	isSub role = "sub"
)

func newPub(opts *options) pubsub {
	return newPubSub(isPub, opts)
}

func newSub(opts *options) pubsub {
	return newPubSub(isSub, opts)
}

func newPubSub(role role, opts *options) pubsub {
	conn, err := redis.Dial(
		"tcp", opts.host, redis.DialPassword(opts.password), redis.DialDatabase(opts.db))
	check(err)
	return pubsub{
		options: opts,
		role:    role,
		conn:    conn,
		counter: 0,
		log:     log.New(os.Stderr, fmt.Sprintf("[%s] ", role), log.Ltime|log.Lmicroseconds),
	}
}

func (pub pubsub) publish() {
	for {
		msg := fmt.Sprintf("[%04d] %s", pub.counter, choice(pub.messages))
		n, err := redis.Int(pub.conn.Do("PUBLISH", pub.channel, msg))
		check(err)
		pub.counter++

		msg = fmt.Sprintf("message received by (%d) recipient", n)
		if n != 1 {
			pub.log.Print(msg, "s")
		} else {
			pub.log.Print(msg)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (sub pubsub) subscribe() {
	psc := redis.PubSubConn{Conn: sub.conn}
	psc.Subscribe(sub.channel)

	for sub.counter < 10 {
		switch v := psc.Receive().(type) {
		case redis.Message:
			sub.log.Printf(
				"channel: %s | message: %s\n",
				v.Channel, v.Data)
		case redis.Subscription:
			sub.log.Printf(
				"channel: %s | kind: %s | count: %d\n",
				v.Channel, v.Kind, v.Count)
		case error:
			panic(v)
		}
	}
	err := psc.Unsubscribe(sub.channel)
	check(err)
}

func choice(options []string) string {
	i := rand.Intn(len(options))
	return options[i]
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
