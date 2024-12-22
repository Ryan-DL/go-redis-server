package commands

import (
	"fmt"

	"github.com/Ryan-DL/go-redis-server/response"
)

func (ch *CommandHandler) HandleInfo() {
	info := fmt.Sprintf(`
# Server
redis_version: redis-experimental
total_keys: %d
`,
		len(ch.MemoryStore.GetKeys()),
	)

	response.SendBulkString(ch.Conn, info)
}
