package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Ryan-DL/go-redis-server/commands"
	"github.com/Ryan-DL/go-redis-server/config"
	"github.com/Ryan-DL/go-redis-server/util"
)

func handleConnection(conn net.Conn, memoryStore *util.ValueStore, password string) {
	defer func() {
		log.Printf("Closing connection from %s", conn.RemoteAddr())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	authenticated := false

	for {
		prefix, err := reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading prefix from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		if prefix != '*' { // all redis commands are of an array type
			util.SendError(conn, "Protocol error: expected '*', got '"+string(prefix)+"'")
			return
		}

		// extract number of args
		line, err := reader.ReadString('\n')
		if err != nil {
			util.SendError(conn, "Protocol error: unable to read array length for command")
			return
		}
		line = strings.TrimSpace(line)
		numArgs, err := strconv.Atoi(line)
		if err != nil {
			util.SendError(conn, "Protocol error: invalid array length")
			return
		}

		command := make([]string, 0, numArgs)
		for i := 0; i < numArgs; i++ {
			bulkPrefix, err := reader.ReadByte()
			if err != nil {
				util.SendError(conn, "Protocol error: unable to read bulk string prefix")
				return
			}
			if bulkPrefix != '$' {
				util.SendError(conn, "Protocol error: expected '$', got '"+string(bulkPrefix)+"'")
				return
			}

			bulkLenStr, err := reader.ReadString('\n')
			if err != nil {
				util.SendError(conn, "Protocol error: unable to read bulk string length")
				return
			}
			bulkLenStr = strings.TrimSpace(bulkLenStr)
			bulkLen, err := strconv.Atoi(bulkLenStr)
			if err != nil {
				util.SendError(conn, "Protocol error: invalid bulk string length")
				return
			}

			// read the actual string
			buf := make([]byte, bulkLen+2) // +2 for \r\n
			_, err = io.ReadFull(reader, buf)
			if err != nil {
				util.SendError(conn, "Protocol error: unable to read bulk string")
				return
			}
			arg := string(buf[:bulkLen])
			command = append(command, arg)
		}

		// if the connection is not authenticated, and the password is set.
		if !authenticated && password != "" {
			if len(command) == 2 && strings.ToUpper(command[0]) == "AUTH" {
				if command[1] == password {
					authenticated = true
					util.SendSimpleString(conn, "OK")
				} else {
					util.SendError(conn, "ERR invalid password")
				}
				continue
			} else {
				util.SendError(conn, "NOAUTH Authentication required.")
				continue
			}
		}

		// handle other commands after authentication
		commandHandler := commands.NewCommandHandler(conn, command, memoryStore)
		handleCommand(commandHandler)
	}
}

func handleCommand(cmd *commands.CommandHandler) {
	cmd.Command[0] = strings.ToUpper(cmd.Command[0])

	switch command := cmd.Command[0]; command {
	case "PING":
		fmt.Printf("DEBUG: cmd.Command[0] = %q\n", cmd.Command[0])
		cmd.HandlePing()
	case "SET":
		cmd.HandleSet()
	case "GET":
		cmd.HandleGet()
	case "DEL":
		cmd.HandleDelete()
	case "EXISTS":
		cmd.HandleExists()
	case "EXPIRE":
		cmd.HandleExpire()
	case "TTL":
		cmd.HandleTTL()
	case "RENAME":
		cmd.HandleRename()
	case "APPEND":
		cmd.HandleAppend()
	case "INCR":
		cmd.HandleIncr()
	case "DECR":
		cmd.HandleDecr()
	case "INFO":
		cmd.HandleInfo()
	default:
		util.SendError(cmd.Conn, "Unknown command: "+command)
	}
}

func main() {
	cfg := config.LoadConfig()

	memoryStore := util.NewValueStore(10 * time.Second)

	var port string
	if cfg.RedisPort != nil {
		port = fmt.Sprintf(":%d", *cfg.RedisPort)
	} else {
		port = ":6379"
	}
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer listener.Close()

	log.Printf("Redis is now running on port %s", port)

	var password string
	if cfg.RedisPassword != nil {
		password = *cfg.RedisPassword
	} else {
		password = ""
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		go handleConnection(conn, memoryStore, password)
	}
}
