package commands

import (
	"strconv"
	"time"

	"github.com/Ryan-DL/go-redis-server/util"
)

func (ch *CommandHandler) HandleIncr() {
	if len(ch.Command) < 2 {
		util.SendError(ch.Conn, "Not enough arguments. Expected INCR key")
		return
	}

	key := ch.Command[1]

	currentValue, exists := ch.MemoryStore.Get(key)
	if !exists {
		// create new key and initialize it to 0 and increment
		ch.MemoryStore.Set(key, "1", 0) // no expiration for a new key
		util.SendInteger(ch.Conn, 1)
		return
	}

	// check we're dealing with an integer
	currentInt, err := strconv.ParseInt(currentValue, 10, 64)
	if err != nil {
		util.SendError(ch.Conn, "ERR value is not an integer or out of range")
		return
	}

	newValue := currentInt + 1

	expiry, hasExpiry := ch.MemoryStore.GetExpiry(key)
	if hasExpiry {
		remainingTTL := time.Until(expiry)
		ch.MemoryStore.Set(key, strconv.FormatInt(newValue, 10), remainingTTL)
	} else {
		ch.MemoryStore.Set(key, strconv.FormatInt(newValue, 10), 0) // no expiration
	}

	util.SendInteger(ch.Conn, int(newValue))
}
