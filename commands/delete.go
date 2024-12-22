package commands

import (
	"github.com/Ryan-DL/go-redis-server/response"
)

func (ch *CommandHandler) HandleDelete() {
	if len(ch.Command) < 2 {
		response.SendError(ch.Conn, "Not enough arguments. Expected DEL key [key ...]")
		return
	}

	keys := ch.Command[1:]

	deletedCount := 0

	for _, key := range keys {
		if ch.MemoryStore.Delete(key) {
			deletedCount++
		}
	}
	response.SendInteger(ch.Conn, deletedCount)
}
