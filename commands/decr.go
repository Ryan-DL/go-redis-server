package commands

import (
	"strconv"
	"time"

	"github.com/Ryan-DL/go-redis-server/response"
)

func (ch *CommandHandler) HandleDecr() {
	if len(ch.Command) < 2 {
		response.SendError(ch.Conn, "Not enough arguments. Expected DECR key")
		return
	}

	key := ch.Command[1]

	currentValue, exists := ch.MemoryStore.Get(key)
	if !exists {
		// Create new key and initialize it to 0, then decrement
		ch.MemoryStore.Set(key, "-1", 0) // No expiration for a new key
		response.SendInteger(ch.Conn, -1)
		return
	}

	currentInt, err := strconv.ParseInt(currentValue, 10, 64)
	if err != nil {
		response.SendError(ch.Conn, "ERR value is not an integer or out of range")
		return
	}

	newValue := currentInt - 1

	expiry, hasExpiry := ch.MemoryStore.GetExpiry(key)
	if hasExpiry {
		remainingTTL := time.Until(expiry)
		ch.MemoryStore.Set(key, strconv.FormatInt(newValue, 10), remainingTTL)
	} else {
		ch.MemoryStore.Set(key, strconv.FormatInt(newValue, 10), 0) // No expiration
	}

	response.SendInteger(ch.Conn, int(newValue))
}
