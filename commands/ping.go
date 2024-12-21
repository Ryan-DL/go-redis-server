package commands

import (
	"github.com/Ryan-DL/go-redis-server/util"
)

func (ch *CommandHandler) HandlePing() {
	//If we're a PING of len one, we can return with a simple string of "PONG."
	if len(ch.Command) == 1 {
		util.SendSimpleString(ch.Conn, "PONG")
		return
	}

	if len(ch.Command) == 2 {
		util.SendBulkString(ch.Conn, ch.Command[1])
		return
	}
	util.SendError(ch.Conn, "Too many arguments for PING")
}
