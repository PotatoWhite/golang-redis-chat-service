```shell
go get github.com/go-redis/redis
go get github.com/gorilla/websocket
go get -u github.com/pinpoint-apm/pinpoint-go-agent

```


# redis 확인
```shell
redis-cli

localhost:6379> subscribe roomname
Reading messages... (press Ctrl-C to quit)
1) "subscribe"
2) "roomname"
3) (integer) 1
```