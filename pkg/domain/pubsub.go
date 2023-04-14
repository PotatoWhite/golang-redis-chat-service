package chat

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var logger = log.New(os.Stdout, "chat: ", log.LstdFlags)

func PubsubListener(ctx context.Context, rdb *redis.Client, room string, handleMessage func(string)) {
	ps := rdb.Subscribe(room)
	_, err := ps.Receive()
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := ps.Channel()

	for {
		select {
		case <-ctx.Done():
			logger.Println("closing pubsub listener")
			_ = ps.Close()
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			handleMessage(msg.Payload)
		}
	}
}

func SetupSignalHandler(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
	}()
}
