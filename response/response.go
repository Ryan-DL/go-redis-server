package response

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

type DataType interface {
	Serialize() string
}

type SimpleString string

func (s SimpleString) Serialize() string {
	return "+" + string(s) + "\r\n"
}

type ErrorType string

func (e ErrorType) Serialize() string {
	return "-" + string(e) + "\r\n"
}

type IntegerType int

func (i IntegerType) Serialize() string {
	return ":" + strconv.Itoa(int(i)) + "\r\n"
}

type BulkStringType string

func (b BulkStringType) Serialize() string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(b), string(b))
}

type NullBulkString struct{}

func (n NullBulkString) Serialize() string {
	return "$-1\r\n"
}

type ArrayType []DataType

func (a ArrayType) Serialize() string {
	if a == nil {
		return "*-1\r\n"
	}
	response := "*" + strconv.Itoa(len(a)) + "\r\n"
	for _, elem := range a {
		response += elem.Serialize()
	}
	return response
}

func writeResponse(conn net.Conn, resp DataType) {
	response := resp.Serialize()
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error sending RESP: %v", err)
	}
}

// Helper Functions to write responses
func SendSimpleString(conn net.Conn, msg string) {
	response := SimpleString(msg)
	writeResponse(conn, response)
}

func SendError(conn net.Conn, msg string) {
	response := ErrorType(msg)
	writeResponse(conn, response)
}

func SendInteger(conn net.Conn, value int) {
	response := IntegerType(value)
	writeResponse(conn, response)
}

func SendBulkString(conn net.Conn, msg string) {
	response := BulkStringType(msg)
	writeResponse(conn, response)
}

func SendNullString(conn net.Conn) {
	response := NullBulkString{}
	writeResponse(conn, response)
}
