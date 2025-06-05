package redis_utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/shahfaizansr/models"
)

type RedisClientModel struct {
	RedisClient *redis.Client
	logger      *logharbour.Logger
}
type RedisTypes int

const (
	INT RedisTypes = iota
	STRING
	JSON
	BOOL
)

var (
	ErrKeyDoesNotExist  = errors.New("key does not exist")
	ErrEmptyValueForKey = errors.New("empty value for key")
)

type RedisClientInterface interface {
	GetValueFromCache(ctx context.Context, key string, valueType RedisTypes) (any, error)
	SetValueInCache(context.Context, SetKeyValue) error
	IsKeyExists(ctx context.Context, key string) (bool, error)
	DeleteKeys(ctx context.Context, keys ...string) error
	PerFormRedisTx(ctx context.Context, redisTxFunc func(tx *redis.Tx) error) error
	PerformRedisPipeLine(ctx context.Context, redisPipelineFunc func(redis.Pipeliner) error) ([]redis.Cmder, error)
}

func (r *RedisClientModel) PerFormRedisTx(ctx context.Context, redisTxFunc func(tx *redis.Tx) error) error {

	if err := r.RedisClient.Watch(ctx, redisTxFunc); err != nil {
		return fmt.Errorf("redis_tx cancelled %v", err)
	}
	return nil
}

func (r *RedisClientModel) PerformRedisPipeLine(ctx context.Context, redisPipelineFunc func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	return r.RedisClient.Pipelined(ctx, redisPipelineFunc)
}

// GetValueFromCache retrieves the value associated with the given key from the Redis cache.
//
// Parameters:
// - ctx: The context.Context object for the request.
// - key: The key for which the value needs to be retrieved.
// - valueType: The expected type of the value.
//
// Returns:
// - any: The value associated with the key, in the specified type.
// - error: An error if the key does not exist or if there is an error retrieving the value.
func (r *RedisClientModel) GetValueFromCache(ctx context.Context, key string, valueType RedisTypes) (any, error) {

	val, err := r.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key does not exist: %v", key)
		}
		return nil, err
	}
	if val == "" {
		return nil, errors.New("empty value for key: " + key)

	}

	switch valueType {
	case INT:
		intValue, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.New("couldn't convert value into integer format: " + val)
		}
		return intValue, nil

	case STRING:
		return val, nil

	case JSON:
		if !json.Valid([]byte(val)) {
			return nil, errors.New("invalid json")
		}
		return val, nil

	}
	return val, nil
}

type SetKeyValue struct {
	Key            string
	Value          any
	ExpirationTime time.Duration
	SetAsAJson     bool
}

// SetValueInCache sets a key-value pair in the Redis cache.
//
// Parameters:
// - ctx: The context.Context object for the request.
// - setValue: The SetKeyValue struct containing the key, value, and expiration time.
//
// Returns:
// - error: An error if there was an issue setting the value in the Redis cache.
func (r *RedisClientModel) SetValueInCache(ctx context.Context, setValue SetKeyValue) error {
	if setValue.SetAsAJson {
		jsonValue, err := json.Marshal(setValue.Value)
		if err != nil {
			return err
		}
		setValue.Value = string(jsonValue)
	}
	redisSetStauts := r.RedisClient.Set(ctx, setValue.Key, setValue.Value, setValue.ExpirationTime)
	return redisSetStauts.Err()
}

// IsKeyExists checks if the given key exists in the Redis cache.
//
// Parameters:
// - ctx: The context.Context object for the request.
// - key: The key to check for existence in the Redis cache.
func (r *RedisClientModel) IsKeyExists(ctx context.Context, key string) (bool, error) {

	result, err := r.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisClientModel) DeleteKeys(ctx context.Context, keys ...string) error {

	_, err := r.RedisClient.Del(ctx, keys...).Result()
	if err != nil {
		return err
	}
	return nil
}

var redisC *redis.Client

// NewRedisClient creates a new Redis client with the given context, logger, and app configuration.
//
// Parameters:
// - ctx: The context.Context object for the request.
// - logger: The *logharbour.Logger object for logging.
// - appConfig: The models.AppConfig object containing the Redis configuration.
//
// Returns:
// - RedisClientInterface: The Redis client interface.
func NewRedisClient(ctx context.Context, logger *logharbour.Logger, appConfig models.AppConfig) RedisClientInterface {

	redisClientModel := RedisClientModel{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     appConfig.Redis.Address,
			Password: appConfig.Redis.Password,
			DB:       0,
		}),
		logger: logger,
	}
	status := redisClientModel.RedisClient.Ping(ctx)
	if status.Err() != nil {
		logger.Err().Log(status.Err().Error())
	}
	redisC = redisClientModel.RedisClient
	return &redisClientModel
}

func GetRedisClient() *redis.Client {
	return redisC

}
