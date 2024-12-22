package commands

import (
	"github.com/Ryan-DL/go-redis-server/util"
)

func (ch *CommandHandler) HandleExists() {
	if len(ch.Command) < 2 {
		util.SendError(ch.Conn, "Not enough arguments. Expected EXISTS key [key ...]")
		return
	}

	keys := ch.Command[1:]

	existsCount := 0

	for _, key := range keys {
		if _, ok := ch.MemoryStore.Get(key); ok {
			existsCount++
		}
	}

	util.SendInteger(ch.Conn, existsCount)
}
