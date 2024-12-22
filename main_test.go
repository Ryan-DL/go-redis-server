package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NOTE:
// THIS TEST WILL ONLY WORK IF YOU HAVE BUILT THE IMAGE LOCALLY

var (
	redisContainer testcontainers.Container
	redisClient    *redis.Client
	ctx            = context.Background()
)

func TestMain(m *testing.M) {
	req := testcontainers.ContainerRequest{
		Image:        "go-redis-server:latest",
		ExposedPorts: []string{"6379/tcp"},
		Env: map[string]string{
			"REDIS_PASSWORD": "securepassword",
		},
		WaitingFor: wait.ForListeningPort("6379/tcp").WithStartupTimeout(30 * time.Second),
	}

	var err error
	redisContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start Redis container: %s", err)
	}

	host, err := redisContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %s", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatalf("Failed to get mapped port: %s", err)
	}

	address := fmt.Sprintf("%s:%s", host, port.Port())

	redisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "securepassword",
	})

	// verify connection
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis server: %s", err)
	}
	if pong != "PONG" {
		log.Fatalf("Unexpected ping response: %s", pong)
	}

	log.Printf("Successfully connected to Redis server at %s", address)

	code := m.Run()

	if err := redisContainer.Terminate(ctx); err != nil {
		log.Fatalf("Failed to terminate Redis container: %s", err)
	}

	os.Exit(code)
}

func TestPing(t *testing.T) {
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Failed to ping Redis server: %s", err)
	}
	if pong != "PONG" {
		t.Fatalf("Unexpected ping response: %s", pong)
	}
	t.Logf("Ping response: %s", pong)
}

func TestSetAndGet(t *testing.T) {
	key := "testKey"
	value := "testValue"

	// Set the key
	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key in Redis: %s", err)
	}

	retrieved, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get key from Redis: %s", err)
	}

	if retrieved != value {
		t.Fatalf("Expected value '%s', got '%s'", value, retrieved)
	}

	t.Logf("Successfully set and retrieved key '%s' with value '%s'", key, retrieved)
}

func TestKeyDoesNotExist(t *testing.T) {
	key := "nonExistentKey"

	_, err := redisClient.Get(ctx, key).Result()
	if err == nil {
		t.Fatalf("Expected error when getting non-existent key '%s', but got none", key)
	}

	if err != redis.Nil {
		t.Fatalf("Expected redis.Nil error, got: %s", err)
	}

	t.Logf("Correctly received redis.Nil for non-existent key '%s'", key)
}

func TestDelKey(t *testing.T) {
	key := "testDelKey"
	value := "testValue"

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key '%s': %s", key, err)
	}

	deleted, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		t.Fatalf("Unexpected error when deleting key '%s': %s", key, err)
	}

	if deleted == 0 {
		t.Fatalf("Expected at least one key to be deleted for key '%s', but none were deleted", key)
	}

	t.Logf("Successfully deleted key '%s'", key)

	_, err = redisClient.Get(ctx, key).Result()
	if err == nil {
		t.Fatalf("Expected error when getting deleted key '%s', but got none", key)
	}

	if err != redis.Nil {
		t.Fatalf("Expected redis.Nil for deleted key '%s', got: %s", key, err)
	}

	t.Logf("Verified key '%s' no longer exists", key)
}

func TestExists(t *testing.T) {
	key := "testExistsKey"
	value := "testValue"

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key in Redis: %s", err)
	}

	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to check existence of key '%s': %s", key, err)
	}

	if exists != 1 {
		t.Fatalf("Expected key '%s' to exist, but it does not", key)
	}

	t.Logf("Verified key '%s' exists", key)
}

func TestExpireAndTTL(t *testing.T) {
	key := "testExpireKey"
	value := "testValue"

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key '%s': %s", key, err)
	}

	err = redisClient.Expire(ctx, key, 5*time.Second).Err()
	if err != nil {
		t.Fatalf("Failed to set expiration for key '%s': %s", key, err)
	}

	ttl, err := redisClient.TTL(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get TTL for key '%s': %s", key, err)
	}

	if ttl <= 0 {
		t.Fatalf("Expected positive TTL for key '%s', got %s", key, ttl)
	}

	t.Logf("TTL for key '%s' is %s", key, ttl)
}

func TestRename(t *testing.T) {
	key := "testRenameKey"
	newKey := "renamedKey"
	value := "testValue"

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key '%s': %s", key, err)
	}

	err = redisClient.Rename(ctx, key, newKey).Err()
	if err != nil {
		t.Fatalf("Failed to rename key '%s' to '%s': %s", key, newKey, err)
	}

	retrieved, err := redisClient.Get(ctx, newKey).Result()
	if err != nil {
		t.Fatalf("Failed to get renamed key '%s': %s", newKey, err)
	}

	if retrieved != value {
		t.Fatalf("Expected value '%s', got '%s'", value, retrieved)
	}

	t.Logf("Successfully renamed key '%s' to '%s' with value '%s'", key, newKey, retrieved)
}

func TestAppend(t *testing.T) {
	key := "testAppendKey"
	value := "initial"
	appendValue := "Appended"

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key '%s': %s", key, err)
	}

	newLength, err := redisClient.Append(ctx, key, appendValue).Result()
	if err != nil {
		t.Fatalf("Failed to append to key '%s': %s", key, err)
	}

	retrieved, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get key '%s': %s", key, err)
	}

	if retrieved != value+appendValue {
		t.Fatalf("Expected value '%s', got '%s'", value+appendValue, retrieved)
	}

	t.Logf("Successfully appended to key '%s'. New length: %d, value: '%s'", key, newLength, retrieved)
}

func TestIncr(t *testing.T) {
	key := "testIncrKey"

	err := redisClient.Set(ctx, key, 0, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key '%s': %s", key, err)
	}

	// increment the key
	newValue, err := redisClient.Incr(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to increment key '%s': %s", key, err)
	}

	if newValue != 1 {
		t.Fatalf("Expected value 1 after increment, got %d", newValue)
	}

	t.Logf("Successfully incremented key '%s'. New value: %d", key, newValue)
}

func TestDecr(t *testing.T) {
	key := "testDecrKey"

	err := redisClient.Set(ctx, key, 0, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set key '%s': %s", key, err)
	}

	// decrement the key
	newValue, err := redisClient.Decr(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to decrement key '%s': %s", key, err)
	}

	if newValue != -1 {
		t.Fatalf("Expected value -1 after decrement, got %d", newValue)
	}

	t.Logf("Successfully decremented key '%s'. New value: %d", key, newValue)
}
