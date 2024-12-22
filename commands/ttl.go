package commands

import (
	"time"

	"github.com/Ryan-DL/go-redis-server/response"
)

func (ch *CommandHandler) HandleTTL() {
	if len(ch.Command) < 2 {
		response.SendError(ch.Conn, "Not enough arguments. Expected TTL key")
		return
	}

	key := ch.Command[1]

	_, ok := ch.MemoryStore.Get(key)
	if !ok {
		response.SendInteger(ch.Conn, -2) // Key does not exist
		return
	}

	expiry, hasExpiry := ch.MemoryStore.GetExpiry(key)
	if !hasExpiry {
		response.SendInteger(ch.Conn, -1)
		return
	}

	ttl := int(time.Until(expiry).Seconds())
	if ttl < 0 {
		response.SendInteger(ch.Conn, -2)
		return
	}

	response.SendInteger(ch.Conn, ttl)
}
