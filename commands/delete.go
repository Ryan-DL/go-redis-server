package commands

import (
	"github.com/Ryan-DL/go-redis-server/util"
)

func (ch *CommandHandler) HandleDelete() {
	if len(ch.Command) < 2 {
		util.SendError(ch.Conn, "Not enough arguments. Expected DEL key [key ...]")
		return
	}

	keys := ch.Command[1:]

	deletedCount := 0

	for _, key := range keys {
		if ch.MemoryStore.Delete(key) {
			deletedCount++
		}
	}
	util.SendInteger(ch.Conn, deletedCount)
}
