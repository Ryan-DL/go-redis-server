package commands

import (
	"strconv"
	"time"

	"github.com/Ryan-DL/go-redis-server/response"
)

// The official documentation provdies several options on the logic to expire a key.
// Because this is for fun, I'm going to leave this as an exercise to the reader.
// https://redis.io/docs/latest/commands/expire/

func (ch *CommandHandler) HandleExpire() {
	if len(ch.Command) < 3 {
		response.SendError(ch.Conn, "Not enough arguments. Expected EXPIRE key seconds")
		return
	}

	key := ch.Command[1]
	seconds, err := strconv.Atoi(ch.Command[2])
	if err != nil || seconds < 0 {
		response.SendError(ch.Conn, "Invalid seconds argument")
		return
	}

	value, exists := ch.MemoryStore.Get(key)
	if !exists {
		response.SendInteger(ch.Conn, 0)
		return
	}

	ttl := time.Duration(seconds) * time.Second
	ch.MemoryStore.Set(key, value, ttl)
	response.SendInteger(ch.Conn, 1)
}
