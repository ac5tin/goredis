package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// ConsumeFunc consumes message at the channel.
type ConsumeFunc func(channel string, message []byte) error

// Publish to a channel with string message
func (r *Client) Publish(channel, message string) (int, error) {
	c := r.pool.Get()
	defer c.Close()
	n, err := redis.Int(c.Do("PUBLISH", channel, message))
	if err != nil {
		return 0, fmt.Errorf("redis publish %s %s, err: %v", channel, message, err)
	}
	return n, nil
}

// Subscribe subscribes messages at the channels.
func (r *Client) Subscribe(ctx context.Context, consume ConsumeFunc, ready *chan bool, channel ...string) error {
	psc := redis.PubSubConn{Conn: r.pool.Get()}

	log.Printf("redis pubsub subscribe channel: %v", channel)
	if err := psc.Subscribe(redis.Args{}.AddFlat(channel)...); err != nil {
		return err
	}
	done := make(chan error, 1)
	// start a new goroutine to receive message
	go func() {
		defer psc.Close()
		for {
			switch msg := psc.Receive().(type) {
			case error:
				done <- fmt.Errorf("redis pubsub receive err: %v", msg)
				return
			case redis.Message:
				if err := consume(msg.Channel, msg.Data); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				if msg.Count == 0 {
					// all channels are unsubscribed
					done <- nil
					return
				}
			}
		}
	}()
	*ready <- true
	// health check
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			if err := psc.Unsubscribe(); err != nil {
				return fmt.Errorf("redis pubsub unsubscribe err: %v", err)
			}
			return nil
		case err := <-done:
			return err
		case <-tick.C:
			if err := psc.Ping(""); err != nil {
				return err
			}
		}
	}

}
