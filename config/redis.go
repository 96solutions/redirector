package config

import (
	"fmt"
)

// RedisConf holds Redis connection configuration.
type RedisConf struct {
	// Host is the Redis server hostname
	Host string `mapstructure:"redis_host"`
	// Port is the Redis server port
	Port string `mapstructure:"redis_port"`
	// Password is the Redis server password
	Password string `mapstructure:"redis_pass"`
	// DB is the Redis database number
	DB int `mapstructure:"redis_db"`
	// MaxRetries is the maximum number of retries before giving up
	MaxRetries int `mapstructure:"redis_max_retries"`
	// MinRetryBackoff is the minimum backoff between each retry
	MinRetryBackoff int `mapstructure:"redis_min_retry_backoff"`
	// MaxRetryBackoff is the maximum backoff between each retry
	MaxRetryBackoff int `mapstructure:"redis_max_retry_backoff"`
	// DialTimeout is the timeout for establishing new connections
	DialTimeout int `mapstructure:"redis_dial_timeout"`
	// ReadTimeout is the timeout for socket reads
	ReadTimeout int `mapstructure:"redis_read_timeout"`
	// WriteTimeout is the timeout for socket writes
	WriteTimeout int `mapstructure:"redis_write_timeout"`
	// PoolSize is the maximum number of socket connections
	PoolSize int `mapstructure:"redis_pool_size"`
	// PoolTimeout is the amount of time client waits for connection if all connections are busy
	PoolTimeout int `mapstructure:"redis_pool_timeout"`
}

// Addr returns the Redis server address in host:port format
func (c *RedisConf) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
