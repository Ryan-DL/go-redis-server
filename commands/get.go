package commands

import (
	"github.com/Ryan-DL/go-redis-server/util"
)

func (ch *CommandHandler) HandleGet() {
	if len(ch.Command) < 2 {
		util.SendError(ch.Conn, "Not enough arguments. Expected GET key")
		return
	}

	key := ch.Command[1]

	value, ok := ch.MemoryStore.Get(key)
	if !ok {
		util.SendError(ch.Conn, "Key not found or expired")
		return
	}

	util.SendBulkString(ch.Conn, value)
}
