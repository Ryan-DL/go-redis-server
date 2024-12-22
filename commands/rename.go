package commands

import (
	"time"

	"github.com/Ryan-DL/go-redis-server/util"
)

func (ch *CommandHandler) HandleRename() {
	if len(ch.Command) < 3 {
		util.SendError(ch.Conn, "Not enough arguments. Expected RENAME key newkey")
		return
	}

	key := ch.Command[1]
	newKey := ch.Command[2]

	// retrieve the key's value and ensure it exists
	value, exists := ch.MemoryStore.Get(key)
	if !exists {
		util.SendError(ch.Conn, "ERR no such key")
		return
	}

	expiry, hasExpiry := ch.MemoryStore.GetExpiry(key)

	// rename by first deleting the newKey
	ch.MemoryStore.Delete(newKey)

	if hasExpiry {
		remainingTTL := time.Until(expiry)
		ch.MemoryStore.Set(newKey, value, remainingTTL)
	} else {
		ch.MemoryStore.Set(newKey, value, 0) // No TTL
	}

	ch.MemoryStore.Delete(key)

	util.SendSimpleString(ch.Conn, "OK")
}
