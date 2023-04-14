package main

import (
	chat "chat-service/pkg/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var logger = log.New(os.Stdout, "main: ", log.LstdFlags)

var upgrader = websocket.Upgrader{}
var chatRooms = make(map[string]*chat.ChatRoom)
var rdb *redis.Client
var redisAddr = "localhost:6379"
var mu sync.Mutex

func main() {
	_redisAddr := os.Getenv("APP_REDIS_ADDR")
	if _redisAddr != "" {
		redisAddr = _redisAddr
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	initRedisOrExit()
	chat.SetupSignalHandler(cancel)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleWebSocket)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatal(err)
	}
}

func initRedisOrExit() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		logger.Fatal(err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Println(err)
		return
	}

	defer conn.Close()

	ctx, cancel := context.WithCancel(rdb.Context())
	defer cancel()

	handlePubsubMessage := func(message string) {
		mu.Lock()
		defer mu.Unlock()
		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}

	var user *chat.User

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var data map[string]interface{}
		json.Unmarshal(message, &data)

		action := data["action"].(string)

		switch action {
		case "join":
			roomName := data["room"].(string)
			user = &chat.User{
				ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
				Nickname: data["nickname"].(string),
			}

			room := getOrCreateChatRoom(roomName)
			room.AddClient(conn)

			rdb.Publish(roomName, fmt.Sprintf("%s joined the room.", user.Nickname))

			go chat.PubsubListener(ctx, rdb, roomName, handlePubsubMessage)
		case "message":
			roomName := data["room"].(string)
			text := data["text"].(string)

			room := getOrCreateChatRoom(roomName)

			rdb.Publish(roomName, fmt.Sprintf("%s: %s", user.Nickname, text))

			room.Broadcast(messageType, message)
			room.RemoveClient(conn)
		default:
			logger.Println("Unknown action: ", action)
		}
	}
}

func getOrCreateChatRoom(name string) *chat.ChatRoom {
	room, exists := chatRooms[name]
	if !exists {
		room = chat.NewChatRoom(name)
		chatRooms[name] = room
	}
	return room
}
