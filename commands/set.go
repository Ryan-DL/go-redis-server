package commands

import (
	"strconv"
	"time"

	"github.com/Ryan-DL/go-redis-server/response"
)

func (ch *CommandHandler) HandleSet() {
	if len(ch.Command) < 3 {
		response.SendError(ch.Conn, "Not enough arguments. Expected SET key value")
	}

	key := ch.Command[1]
	value := ch.Command[2]

	// I'm choosing not to support other times besides seconds.
	expiration := false
	for _, v := range ch.Command {
		if v == "EX" {
			expiration = true
		}
	}

	if expiration {
		last := ch.Command[len(ch.Command)-1]
		secondsToAdd, err := strconv.ParseInt(last, 10, 64)
		if err != nil {
			response.SendError(ch.Conn, "Unable to parse requested expiration")
			return
		}
		expirey := time.Duration(secondsToAdd) * time.Second
		ch.MemoryStore.Set(key, value, expirey)
	} else {
		ch.MemoryStore.Set(key, value, 0)
	}

	response.SendSimpleString(ch.Conn, "OK")
}
