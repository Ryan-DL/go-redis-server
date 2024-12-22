package commands

import (
	"fmt"

	"github.com/Ryan-DL/go-redis-server/util"
)

func (ch *CommandHandler) HandleInfo() {
	info := fmt.Sprintf(`
# Server
redis_version: redis-experimental
total_keys: %d
`,
		len(ch.MemoryStore.GetKeys()),
	)

	util.SendBulkString(ch.Conn, info)
}
