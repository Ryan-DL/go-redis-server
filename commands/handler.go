package commands

import (
	"net"

	"github.com/Ryan-DL/go-redis-server/cache"
)

type CommandHandler struct {
	Conn        net.Conn
	Command     []string
	MemoryStore *cache.ValueStore
}

func NewCommandHandler(conn net.Conn, command []string, memoryStore *cache.ValueStore) *CommandHandler {
	return &CommandHandler{
		Conn:        conn,
		Command:     command,
		MemoryStore: memoryStore,
	}
}
