package main

import (
	"acbot/types"
)

type QueueInterface interface {
	AddNew(key string, a *types.Activation) error
	GetNext(key string) (*types.Activation, error)
}

type CurrentActivationInterface interface {
	SetNew(a *types.Activation) error
	SetActivator(ActivatorId int64) error
	SetComplete() error
	SetRetry() error
}

type Queue struct {
	Redis *Redis
}

type CurrentActivation struct {
	Redis *Redis
}

func (g *Queue) AddNew(key string, a *types.Activation) (e error) {
	e = g.Redis.PushToQueue(key, a)
	return
}

func (g *Queue) GetNext(key string) (a *types.Activation, e error) {
	a, e = g.Redis.PopFromQueue(key)
	return
}

// Set new current activation for updating
func (c *CurrentActivation) SetNew(a *types.Activation) (e error) {
	e = c.Redis.SetToKey(c.Redis.CurrentKey, a)
	return
}

// Set activator to current activation
func (c *CurrentActivation) SetActivator(ActivatorId int64) error {
	curActivator, err := c.Redis.GetFromKey(c.Redis.CurrentKey)
	if err != nil {
		return err
	}
	curActivator.Activator = ActivatorId
	err = c.Redis.SetToKey(c.Redis.CurrentKey, curActivator)
	return err
}

// Set complete flag to current activation
func (c *CurrentActivation) SetComplete() (e error) {
	curActivator, err := c.Redis.GetFromKey(c.Redis.CurrentKey)
	if err != nil {
		return err
	}
	curActivator.Complete = true
	err = c.Redis.SetToKey(c.Redis.CurrentKey, curActivator)
	return err
}

// Set retry flag for current activation
func (c *CurrentActivation) SetRetry() (e error) {
	curActivator, err := c.Redis.GetFromKey(c.Redis.CurrentKey)
	if err != nil {
		return err
	}
	curActivator.Retry = true
	err = c.Redis.SetToKey(c.Redis.CurrentKey, curActivator)
	return err
}
