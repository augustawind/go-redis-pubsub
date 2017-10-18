package main

import (
	"fmt"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
)

type ClientOptions struct {
	role     string
	host     string
	password string
	db       int
	//pubChannels []string
	//subChannels []string
}

func (opt *ClientOptions) newConn() redis.Conn {
	conn, err := redis.Dial(
		"tcp", opt.host, redis.DialPassword(opt.password), redis.DialDatabase(opt.db))
	check(err)
	return conn
}

func (opt *ClientOptions) newLogger(id uuid.UUID) *log.Logger {
	return log.New(os.Stderr, fmt.Sprintf("[%s:%s] ", opt.role, id.String()[:8]), log.Ltime|log.Lmicroseconds)
}

type Node struct {
	*ClientOptions
	id      uuid.UUID
	conn    redis.Conn
	subConn redis.PubSubConn
	log     *log.Logger
}

func NewNode(opts *ClientOptions) *Node {
	id := uuid.NewV4()
	return &Node{
		id:      id,
		conn:    opts.newConn(),
		subConn: redis.PubSubConn{Conn: opts.newConn()},
		log:     opts.newLogger(id),
	}
}

func (m *Node) Publish(channel string, msg string) int {
	n, err := redis.Int(m.conn.Do("PUBLISH", channel, msg))
	check(err)

	m.log.Printf("---- publish ----\n\tsend MSG\t=\t<< %s >>\n\tthru CHANNEL\t=\t<< %s >>\n\t# of CLIENTS\t=\t<< %d >>",
		msg, channel, n)
	// TODO: wait for feedback from subscribers before sending to channel
	return n
}

func (m *Node) Subscribe(channel string, c chan []byte) {
	m.subConn.Subscribe(channel)

	for {
		switch v := m.subConn.Receive().(type) {
		case redis.Message:
			m.log.Printf("--- subscribe ---\n\treceive MSG\t=\t<< %s >>\n\tthru CHANNEL\t=\t<< %s >>",
				v.Data, v.Channel)
			c <- v.Data
		case redis.PMessage:
			m.log.Printf("--- subscribe ---\n\treceive PMSG\t=\t<< %s >>\n\tthru CHANNEL\t=\t<< %s >>\n\twith PATTERN\t=\t<< %s >>",
				v.Data, v.Channel, v.Pattern)
			c <- v.Data
		case redis.Subscription:
			m.log.Printf("--- subscribe ---\n\treceive ACTION\t=\t<< %s >>\n\tthru CHANNEL\t=\t<< %s >>\n\t# of CLIENTS\t=\t<< %d >>",
				v.Kind, v.Channel, v.Count)
		case error:
			panic(v)
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
