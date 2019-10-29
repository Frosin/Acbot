// Redis integration tests
package main

import (
	"acbot/types"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var testAct = types.Activation{
	Timestamp: time.Now(),
	User:      123456,
	Activator: 9876543,
	Complete:  false,
	Retry:     false,
}

func newTestRedis(t *testing.T, addr string) (Redis, error) {
	mr, err := miniredis.Run()
	if !assert.NoError(t, err, "Failed to run miniredis!") {
		assert.FailNow(t, "Can't connect to Miniredis, testing failed!")
	}
	rediska := new(Redis)
	if addr == "" {
		addr = mr.Addr()
	}
	err = rediska.Connect(&redis.Options{
		Addr: addr,
	})
	return *rediska, err
}

func TestCheckEmptyEnvs(t *testing.T) {
	testEnv1 := []string{"1", "2", "", "3"}
	testEnv2 := []string{"1", "2", "3", "4"}
	if checkNoEmptyEnvs(testEnv1) {
		t.Error("Empty check func fail!")
	}
	if !checkNoEmptyEnvs(testEnv2) {
		t.Error("Empty check func fail!")
	}
}

func TestRedisConnect(t *testing.T) {
	_, err := newTestRedis(t, "")
	assert.NoError(t, err, "Failed connect to Redis!")
}

func TestBadEnvs(t *testing.T) {
	_, err := newTestRedis(t, "Bab address")
	assert.Error(t, err, "Bad env check fail!")
}

func TestDropKey(t *testing.T) {
	rediska, err := newTestRedis(t, "")
	assert.NoError(t, err, "Failed connect to Redis!")
	err = rediska.DropKey(rediska.CurrentKey)
	assert.NoError(t, err, "Drop key fail!")
}

func TestPushToQueue(t *testing.T) {
	rediska, err := newTestRedis(t, "")
	assert.NoError(t, err, "Failed connect to Redis!")
	err = rediska.PushToQueue(rediska.GeneralKey, &testAct)
	assert.NoError(t, err, "Add to general queue failed!")
}

func TestPopFromQueue(t *testing.T) {
	rediska, err := newTestRedis(t, "")
	assert.NoError(t, err, "Failed connect to Redis!")
	err = rediska.PushToQueue(rediska.GeneralKey, &testAct)
	assert.NoError(t, err, "Add to general queue failed!")
	activation, err := rediska.PopFromQueue(rediska.GeneralKey)
	assert.Equal(t, 123456, int(activation.User), "Bad data in queue!")
	assert.NoError(t, err, "Getting from general queue fail!")
}

func TestSetToKey(t *testing.T) {
	rediska, err := newTestRedis(t, "")
	assert.NoError(t, err, "Failed connect to Redis!")
	err = rediska.SetToKey(rediska.CurrentKey, &testAct)
	assert.NoError(t, err, "Failed set value to key!")
}

func TestGetFromKey(t *testing.T) {
	rediska, err := newTestRedis(t, "")
	assert.NoError(t, err, "Failed connect to Redis!")
	err = rediska.SetToKey(rediska.CurrentKey, &testAct)
	assert.NoError(t, err, "Failed set value to key!")
	activation, err := rediska.GetFromKey(rediska.CurrentKey)
	assert.Equal(t, 123456, int(activation.User), "Bad data key!")
	assert.NoError(t, err, "Getting from key queue fail!")
}

func TestGetQueueLength(t *testing.T) {
	rediska, err := newTestRedis(t, "")
	assert.NoError(t, err, "Failed connect to Redis!")
	for i := 1; i <= 4; i++ {
		e := rediska.PushToQueue(rediska.GeneralKey, &testAct)
		assert.NoError(t, e, "Failed to push data!")
	}
	len, err := rediska.GetQueueLength()
	assert.NoError(t, err, "Failed get queue length!")
	assert.Equal(t, 4, int(len), "Entity count not valid!")
}
