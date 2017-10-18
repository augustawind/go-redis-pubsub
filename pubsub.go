package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/garyburd/redigo/redis"
)

var reader = bufio.NewReader(os.Stdin)

type pub struct {
	conn     redis.Conn
	log      *log.Logger
	channel  string
	messages []string
	prompt   bool
	counter  int
}

type sub struct {
	conn    redis.PubSubConn
	log     *log.Logger
	channel string
}

type clientOptions struct {
	host     string
	password string
	db       int
	channel  string
}

func (opt *clientOptions) newConn() redis.Conn {
	conn, err := redis.Dial(
		"tcp", opt.host, redis.DialPassword(opt.password), redis.DialDatabase(opt.db))
	check(err)
	return conn
}

func newLogger(name string) *log.Logger {
	return log.New(os.Stderr, fmt.Sprintf("[%s] ", name), log.Ltime|log.Lmicroseconds)
}

func newPub(opts *clientOptions, prompt bool, messages []string) *pub {
	if len(messages) == 0 {
		messages = defaultMessages
	}
	return &pub{
		conn:     opts.newConn(),
		log:      newLogger("pub"),
		channel:  opts.channel,
		messages: messages,
		prompt:   prompt,
		counter:  0,
	}
}

func newSub(opts *clientOptions) *sub {
	return &sub{
		conn:    redis.PubSubConn{opts.newConn()},
		log:     newLogger("sub"),
		channel: opts.channel,
	}
}

func (p *pub) publish(c chan int) {
	n, err := redis.Int(p.conn.Do(
		"PUBLISH", p.channel, p.newMessage(),
	))
	check(err)

	p.log.Printf("message received by ( %d ) recipient(s)", n)
	p.counter++
	c <- n
}

func (p *pub) newMessage() string {
	return fmt.Sprintf("#%04d :: %s", p.counter, choice(p.messages))
}

func (s *sub) subscribe(c chan bool) {
	s.conn.Subscribe(s.channel)

	for {
		switch v := s.conn.Receive().(type) {
		case redis.Message:
			s.log.Printf("channel: %s | message: %s\n",
				v.Channel, v.Data)
			c <- true
		case redis.PMessage:
			s.log.Printf("channel: %s | message: %s\n",
				v.Channel, v.Data)
			c <- true
		case redis.Subscription:
			s.log.Printf("channel: %s | kind: %s | count: %d\n",
				v.Channel, v.Kind, v.Count)
		case error:
			panic(v)
		}
	}
	// TODO: remove this, keep it rolling
	err := s.conn.Unsubscribe(s.channel)
	check(err)
}

func (p *pub) interact() {
	if p.prompt {
		fmt.Print("Press <Enter>...")
		_, err := reader.ReadString('\n')
		check(err)
	}
}

func choice(xs []string) string {
	i := rand.Intn(len(xs))
	return xs[i]
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
