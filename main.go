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
	"github.com/Ryan-DL/go-redis-server/util"
)

func handleConnection(conn net.Conn, memoryStore *util.ValueStore) {
	defer func() {
		log.Printf("Closing connection from %s", conn.RemoteAddr())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	for {
		prefix, err := reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading prefix from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}

		if prefix != '*' { // All redis commands are of an array type
			util.SendError(conn, "Protocol error: expected '*', got '"+string(prefix)+"'")
			return
		}

		// Read the number of elements in the array
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

			// Read bulk string length
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

			// Read the actual string
			buf := make([]byte, bulkLen+2) // +2 for \r\n
			_, err = io.ReadFull(reader, buf)
			if err != nil {
				util.SendError(conn, "Protocol error: unable to read bulk string")
				return
			}
			arg := string(buf[:bulkLen])
			command = append(command, arg)
		}

		commandHandler := commands.NewCommandHandler(conn, command, memoryStore)

		handleCommand(commandHandler)
	}
}

func handleCommand(cmd *commands.CommandHandler) {
	fmt.Println("Running command: ", cmd.Command)

	switch command := cmd.Command[0]; command {
	case "PING":
		fmt.Printf("DEBUG: cmd.Command[0] = %q\n", cmd.Command[0])
		cmd.HandlePing()
	case "SET":
		cmd.HandleSet()
	case "GET":
		cmd.HandleGet()
	default:
		util.SendError(cmd.Conn, "Unknown command: "+command)
	}
}

func main() {
	memoryStore := util.NewValueStore(10 * time.Second)

	//todo: change port
	listener, err := net.Listen("tcp", ":7000")
	if err != nil {
		log.Fatalf("Failed to listen on port 6379: %v", err)
	}
	defer listener.Close()
	log.Println("Redis is now running on on port 7000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
		go handleConnection(conn, memoryStore)
	}
}
