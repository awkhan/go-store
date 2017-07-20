package store

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

//Redis provides an interface to redis.
type Redis struct {
	redis *redis.Pool
}

//NewRedisStore creates a new redis store with the supplied pool.
func NewRedisStore(pool *redis.Pool) *Redis {
	return &Redis{
		redis: pool,
	}
}

//DeleteKey deletes the key from redis.
func (r *Redis) DeleteKey(key string) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("DEL", key)
	return e
}

//GetString retrieves the string data stored in redis.
func (r *Redis) GetString(key string) (string, error) {
	conn := r.redis.Get()
	defer conn.Close()
	res, e := redis.String(conn.Do("GET", key))
	return res, e
}

//GetInt64 retrieves the int64 data stored in redis.
func (r *Redis) GetInt64(key string) (int64, error) {
	conn := r.redis.Get()
	defer conn.Close()
	res, e := redis.Int64(conn.Do("GET", key))
	return res, e
}

//Set sets the value for the specified key.
func (r *Redis) Set(key string, value interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("SET", key, value)
	return e
}

//SetHash sets the value for the specific hash key.
func (r *Redis) SetHash(key string, hash string, value interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("HSET", key, hash, value)
	return e
}

//DeleteHash deletes the hash value for the specific key.
func (r *Redis) DeleteHash(key string, hash string) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("HDEL", key, hash)
	return e
}

//GetHashString returns the string value of the hash.
func (r *Redis) GetHashString(key string, hash string) (string, error) {
	conn := r.redis.Get()
	defer conn.Close()
	val, e := redis.String(conn.Do("HGET", key, hash))
	return val, e
}

//GetAllHashValues returns all the hash values for the key.
func (r *Redis) GetAllHashValues(key string) ([]string, error) {
	conn := r.redis.Get()
	defer conn.Close()
	val, e := redis.Strings(conn.Do("HVALS", key))
	return val, e
}

//GetAllHashKeys returns all the hash keys for the key.
func (r *Redis) GetAllHashKeys(key string) ([]string, error) {
	conn := r.redis.Get()
	defer conn.Close()
	val, e := redis.Strings(conn.Do("HKEYS", key))
	return val, e
}

//SetExpiry sets the expiry for the specified key.
func (r *Redis) SetExpiry(key string, seconds int) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("EXPIRE", key, seconds)
	return e
}

//Increment increments the value of key by 1.
func (r *Redis) Increment(key string) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("INCR", key)
	return e
}

//Decrement decrements the value of key by 1.
func (r *Redis) Decrement(key string) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("DECR", key)
	return e
}

//SetAdd adds a the value to a set.
func (r *Redis) SetAdd(key string, value interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("SADD", key, value)
	return e
}

//GetSetStringMembers returns the string members of a set.
func (r *Redis) GetSetStringMembers(key string) ([]string, error) {
	conn := r.redis.Get()
	defer conn.Close()
	val, e := redis.Strings(conn.Do("SMEMBERS", key))
	return val, e
}

//SetRemove removes the value from the set.
func (r *Redis) SetRemove(key string, value interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()
	_, e := conn.Do("SREM", key, value)
	return e
}

//SetIsMember returns true if the value is a member of the set.
func (r *Redis) SetIsMember(key string, value interface{}) (bool, error) {
	conn := r.redis.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("SISMEMBER", key, value))
}

//PushItemToList pushes an item to the list. Use inFront to specifiy if the item should go in front of at the end of the list.
func (r *Redis) PushItemToList(key string, value interface{}, atEnd bool) error {
	conn := r.redis.Get()
	defer conn.Close()
	cmd := "LPUSH"
	if atEnd {
		cmd = "RPUSH"
	}

	_, err := conn.Do(cmd, key, value)
	return err
}

//PopItemFromList pops an item from the front or the back of the list.
func (r *Redis) PopItemFromList(key string, dataType int, atEnd bool) (interface{}, error) {
	conn := r.redis.Get()
	defer conn.Close()

	cmd := "LPOP"
	if atEnd {
		cmd = "RPOP"
	}

	val, err := conn.Do(cmd, key)

	switch dataType {
	case DataTypeString:
		return redis.String(val, err)
	case DataTypeBool:
		return redis.Bool(val, err)
	case DataTypeInt:
		return redis.Int(val, err)
	case DataTypeInt64:
		return redis.Int64(val, err)
	default:
		return nil, errors.New("Invalid data type")
	}
}

//ItemsFromList returns a list of items from the list from the start to end.
func (r *Redis) ItemsFromList(key string, dataType int, start, end int) (interface{}, error) {
	conn := r.redis.Get()
	defer conn.Close()

	val, err := conn.Do("LRANGE", key, start, end)

	switch dataType {
	case DataTypeString:
		return redis.Strings(val, err)
	case DataTypeInt:
		return redis.Ints(val, err)
	default:
		return nil, errors.New("Invalid data type")
	}
}

//RemoveItemFromList removes the item from the list with the count occurances.
func (r *Redis) RemoveItemFromList(key string, count int, value interface{}) error {
	conn := r.redis.Get()
	defer conn.Close()

	_, err := conn.Do("LREM", key, count, value)
	return err
}

//LengthOfList returns the lenght of the list.
func (r *Redis) LengthOfList(key string) (int, error) {
	conn := r.redis.Get()
	defer conn.Close()

	return redis.Int(conn.Do("LLEN", key))
}

//ClearDataStore clears up all the keys in the redis datastore.
func (r *Redis) ClearDataStore() {
	conn := r.redis.Get()
	defer conn.Close()
	conn.Do("FLUSHDB")
}
