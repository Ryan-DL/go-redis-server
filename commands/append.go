package commands

import (
	"time"

	"github.com/Ryan-DL/go-redis-server/response"
)

func (ch *CommandHandler) HandleAppend() {
	if len(ch.Command) < 3 {
		response.SendError(ch.Conn, "Not enough arguments. Expected APPEND key value")
		return
	}

	key := ch.Command[1]
	appendValue := ch.Command[2]

	currentValue, exists := ch.MemoryStore.Get(key)
	expiry, hasExpiry := ch.MemoryStore.GetExpiry(key)

	if !exists {
		ch.MemoryStore.Set(key, appendValue, 0) // No expiration for a new key
		response.SendInteger(ch.Conn, len(appendValue))
		return
	}

	// Append the value if the key exists
	newValue := currentValue + appendValue

	// Set the updated value with the same expiration (if any)
	if hasExpiry {
		remainingTTL := time.Until(expiry)
		ch.MemoryStore.Set(key, newValue, remainingTTL)
	} else {
		ch.MemoryStore.Set(key, newValue, 0)
	}

	response.SendInteger(ch.Conn, len(newValue))
}
