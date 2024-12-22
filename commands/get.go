package commands

import "github.com/Ryan-DL/go-redis-server/response"

func (ch *CommandHandler) HandleGet() {
	if len(ch.Command) < 2 {
		response.SendError(ch.Conn, "Not enough arguments. Expected GET key")
		return
	}

	key := ch.Command[1]

	value, ok := ch.MemoryStore.Get(key)
	if !ok {
		response.SendNullString(ch.Conn)
		return
	}

	response.SendBulkString(ch.Conn, value)
}
