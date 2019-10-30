package main

import (
	"acbot/types"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

type Redis struct {
	Client     *redis.Client
	Connected  bool
	GeneralKey string
	DBKey      string
	CurrentKey string
}

// Validate envs for empties
func checkNoEmptyEnvs(envs []string) bool {
	for _, value := range envs {
		if value == "" {
			return false
		}
	}
	return true
}

// Get redis options from .env file
func getConnectOptions(envFile string) (*redis.Options, error) {
	var err error
	if envFile == "" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(envFile)
	}
	if err != nil {
		return &redis.Options{}, err
	}
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	if !checkNoEmptyEnvs([]string{
		addr,
		dbStr,
	}) {
		return nil, errors.New("Your .env file have empty connect variables!")
	}
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, err
	}
	return &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}, nil
}

// Get key value names from .env file
func (r *Redis) setDatabaseKeys() error {
	generalKey := os.Getenv("REDIS_GENERAL_KEY")
	dbKey := os.Getenv("REDIS_DB_KEY")
	currentKey := os.Getenv("REDIS_CURRENT_KEY")
	if !checkNoEmptyEnvs([]string{
		generalKey,
		dbKey,
		currentKey,
	}) {
		return errors.New("Your .env file have empty key value variables!")
	}
	r.GeneralKey = generalKey
	r.DBKey = dbKey
	r.CurrentKey = currentKey
	return nil
}

// Check all envs and connect to Redis
// use fn getConnectOptions for get connectOptions
func (r *Redis) Connect(connectOptions *redis.Options) (err error) {
	r.Client = redis.NewClient(connectOptions)
	err = r.Client.Ping().Err()
	if err == nil {
		r.Connected = true
	}
	return
}

func (r *Redis) DropKey(key string) error {
	err := r.Client.Del(key).Err()
	return err
	// result, err := r.Client.LTrim(key, 100, 0).Result()
	// if result != "OK" {
	// 	errMessage := strings.Join([]string{
	// 		"Drop result not OK: ",
	// 		result,
	// 	}, "")
	// 	return errors.New(errMessage)
	// }
	// return err
}

// Push activation to general queue
func (r *Redis) PushToQueue(key string, activation *types.Activation) (err error) {
	err = r.Client.RPush(key, activation).Err()
	return
}

// Get activation from general queue
func (r *Redis) PopFromQueue(key string) (*types.Activation, error) {
	var activation types.Activation
	err := r.Client.LPop(key).Scan(&activation)
	return &activation, err
}

// Set value to key
func (r *Redis) SetToKey(key string, activation *types.Activation) (err error) {
	err = r.Client.Set(key, activation, time.Hour).Err()
	return
}

// Get value from key
func (r *Redis) GetFromKey(key string) (*types.Activation, error) {
	var activation types.Activation
	err := r.Client.Get(key).Scan(&activation)
	return &activation, err
}

// Get queue length
func (r *Redis) GetQueueLength() (int64, error) {
	len, err := r.Client.LLen(r.GeneralKey).Result()
	return len, err
}
