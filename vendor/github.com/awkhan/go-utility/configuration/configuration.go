package configuration

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"os"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// GetString returns string value of the environment variable key.
func GetString(key string) string {
	return os.Getenv(key)
}

// GetInt returns the integer value of environment variable key. 0 if the conversion cannot occur.
func GetInt(key string) int {
	v := GetString(key)
	iv, err := strconv.Atoi(v)

	if err != nil {
		log.Errorf("Unable to convert %s to int for key %s", v, key)
		iv = 0
	}

	return iv
}

// GetDuration returns the duration value of environment variable key. 0 if the conversion cannot occur.
func GetDuration(key string) time.Duration {
	v := GetInt(key)
	return time.Duration(v)
}

// GetBool returns the boolean value of the environment variable key. Defaults to false if the value does not exist.
func GetBool(key string) bool {
	v := GetString(key)
	bv, err := strconv.ParseBool(v)

	if err != nil {
		log.Errorf("Unable to convert %s to bool for key %s", v, key)
		bv = false
	}

	return bv
}

// GetRedisPool creates an instance of the redis pool.
func GetRedisPool(maxIdle int, idleTimeout int, host, port, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeout),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				log.Error("Unable to get redis connection")
				return nil, err
			}
			if password != "" {
				c.Do("AUTH", password)
			}
			return c, err
		},
	}
}

// GetSQLDatabase returns an instance of the sql.DB.
func GetSQLDatabase(user, password, host, port, databasename string) *sql.DB {
	sqlDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, databasename))
	if err != nil {
		log.Error("Unable to get sql connection. Please make sure the correct import is used")
		return nil
	}
	return sqlDB
}
