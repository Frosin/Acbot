package main

import (
	"acbot/types"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func getTestQueue(t *testing.T) (*Queue, error) {
	q := &Queue{
		Redis: &Redis{},
	}
	mr, err := miniredis.Run()
	if !assert.NoError(t, err, "Failed to run miniRedis!") {
		assert.FailNow(t, "Failed to run tests!")
	}
	q.Redis = new(Redis)
	err = q.Redis.Connect(&redis.Options{
		Addr: mr.Addr(),
	})
	return q, err
}

func getTestCurAct(t *testing.T) (*CurrentActivation, error) {
	c := &CurrentActivation{
		Redis: &Redis{},
	}
	mr, err := miniredis.Run()
	assert.NoError(t, err, "Failed to run miniRedis!")
	c.Redis = new(Redis)
	err = c.Redis.Connect(&redis.Options{
		Addr: mr.Addr(),
	})
	return c, err
}

func TestAddNew(t *testing.T) {
	g, err := getTestQueue(t)
	assert.NoError(t, err, "Error by getting queue!")
	err = g.AddNew(g.Redis.GeneralKey, &testAct)
	assert.NoError(t, err, "Failed to add new to queuet!")
	len, err := g.Redis.GetQueueLength()
	assert.NoError(t, err, "Failed to get length of queue!")
	testAct, err := g.Redis.PopFromQueue(g.Redis.GeneralKey)
	assert.NoError(t, err, "Failed to get value from queue!")
	assert.Equal(t, 123456, int(testAct.User), "Got bad data from queue!")
	assert.Equal(t, 1, int(len), "Got bad data from queue!")
}

func TestGetNext(t *testing.T) {
	g, err := getTestQueue(t)
	assert.NoError(t, err, "Error by getting queue!")
	var testAct1 = &types.Activation{
		Timestamp: time.Now(),
		User:      111111,
		Activator: 9876543,
		Complete:  false,
		Retry:     false,
	}
	var testAct2 = &types.Activation{
		Timestamp: time.Now(),
		User:      222222,
		Activator: 9876543,
		Complete:  false,
		Retry:     false,
	}
	err = g.AddNew(g.Redis.GeneralKey, testAct1)
	assert.NoError(t, err, "Failed to add value!")
	err = g.AddNew(g.Redis.GeneralKey, testAct2)
	assert.NoError(t, err, "Failed to add value!")
	getAct, err := g.GetNext(g.Redis.GeneralKey)
	assert.NoError(t, err, "Failed to get value from queue!")
	len, err := g.Redis.GetQueueLength()
	assert.NoError(t, err, "Failed to get length of queue!")
	assert.Equal(t, 111111, int(getAct.User), "Got bad data from queue!")
	assert.Equal(t, 1, int(len), "Got bad data from queue!")
}

func TestSetNew(t *testing.T) {
	c, err := getTestCurAct(t)
	assert.NoError(t, err, "Error by getting queue!")
	err = c.SetNew(&testAct)
	assert.NoError(t, err, "Failed to set value!")
	current, err := c.Redis.GetFromKey(c.Redis.CurrentKey)
	assert.Equal(t, 123456, int(current.User), "Got bad data from queue!")
	assert.NoError(t, err, "Got bad data from queue!")
}

func TestSetActivator(t *testing.T) {
	c, err := getTestCurAct(t)
	assert.NoError(t, err, "Error by getting queue!")
	err = c.SetNew(&testAct)
	assert.NoError(t, err, "Failed to set value!")
	err = c.SetActivator(100500)
	assert.NoError(t, err, "Failed to set `Activator` field!")
	curActivation, err := c.Redis.GetFromKey(c.Redis.CurrentKey)
	assert.Equal(t, 100500, int(curActivation.Activator), "Bad data in `Activator` field!")
}

func TestSetComplete(t *testing.T) {
	c, err := getTestCurAct(t)
	assert.NoError(t, err, "Error by getting queue!")
	err = c.SetNew(&testAct)
	assert.NoError(t, err, "Failed to set value!")
	err = c.SetComplete()
	assert.NoError(t, err, "Failed to set `Complete` field!")
	curActivation, err := c.Redis.GetFromKey(c.Redis.CurrentKey)
	assert.Equal(t, true, curActivation.Complete, "Bad data in `Complete` field!")
}

func TestSetRetry(t *testing.T) {
	c, err := getTestCurAct(t)
	assert.NoError(t, err, "Error by getting queue!")
	err = c.SetNew(&testAct)
	assert.NoError(t, err, "Failed to set value!")
	err = c.SetRetry()
	assert.NoError(t, err, "Failed to set `Retry` field!")
	curActivation, err := c.Redis.GetFromKey(c.Redis.CurrentKey)
	assert.Equal(t, true, curActivation.Retry, "Bad data in `Retry` field!")
}
