package commands

import (
	"net"

	"github.com/Ryan-DL/go-redis-server/util"
)

type CommandHandler struct {
	Conn        net.Conn
	Command     []string
	MemoryStore *util.ValueStore
}

func NewCommandHandler(conn net.Conn, command []string, memoryStore *util.ValueStore) *CommandHandler {
	return &CommandHandler{
		Conn:        conn,
		Command:     command,
		MemoryStore: memoryStore,
	}
}
