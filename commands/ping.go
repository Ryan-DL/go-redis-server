package commands

import (
	"github.com/Ryan-DL/go-redis-server/response"
)

func (ch *CommandHandler) HandlePing() {
	//If we're a PING of len one, we can return with a simple string of "PONG."
	if len(ch.Command) == 1 {
		response.SendSimpleString(ch.Conn, "PONG")
		return
	}

	if len(ch.Command) == 2 {
		response.SendBulkString(ch.Conn, ch.Command[1])
		return
	}
	response.SendError(ch.Conn, "Too many arguments for PING")
}
