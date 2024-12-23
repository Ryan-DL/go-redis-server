# Go Redis Server

A weekend project implementing a partial, in-memory RESP2-compliant Redis server in Go. Implementation tested against the official `go-redis` client. 

## Build and Execute

```bash
# Run without building
go run main.go 
# Run Test
go test -v 
# Build Docker Container
docker build -t go-redis-server:latest . 
# Run the container 
docker run -d -p 6379:6379 --name go-redis-server \ 
  -e REDIS_PASSWORD=your_redis_password \
  -e REDIS_PORT=6379 \
  go-redis-server
```

## Resources & Libraries Used
* [Redis serialization protocol specification](https://redis.io/docs/latest/develop/reference/protocol-spec/)
* [List of Redis Commands](https://redis.io/docs/latest/commands/)
* [Offical Redis Go Client](https://redis.io/docs/latest/develop/clients/go/)
* [Test Containers](https://testcontainers.com/)

## Implemented Protocol Commands
- AUTH - AUTH command on the client
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

## Caveats 

1. Theres no enforcement on integer overflows.
2. No timeouts, very little validation of input. Theoretically commands can hang forever.
3. There are several sub-operations on the commands that still need implementing.
