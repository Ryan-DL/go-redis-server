## Go Redis Server

A partial, in-memory, RESP2, weekend implementation of a Redis server in Go. Probably not ideal for production.

## Build and Execute

```bash
go run main.go # Just run 
# Docker Build

```

## Resources
* [Redis serialization protocol specification](https://redis.io/docs/latest/develop/reference/protocol-spec/)
* [Offical Go Client](https://redis.io/docs/latest/develop/clients/go/)
* [Test Containers](https://testcontainers.com/)

## Implemented Commands
- GET - Get value of a key
- SET - Set a value of a key
- DEL - Delete a key
- EXISTS - Check if key exists
- EXPIRE - Sets a keys expiration 
- TTL - Get time to live of key
- RENAME - Rename a keys value 
- APPEND - Append value to a key 
- INCR - Increment value of key 
- DECR - Decrement value of key
- PING - PONG!
- INFO - Debug info about the server.

