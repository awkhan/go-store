# go-store

[![CircleCI](https://circleci.com/gh/awkhan/go-store/tree/master.svg?style=svg)](https://circleci.com/gh/awkhan/go-store/tree/master)
[![codecov](https://codecov.io/gh/awkhan/go-store/branch/master/graph/badge.svg)](https://codecov.io/gh/awkhan/go-store)

A storage interface used to do key value data storage. Gives you the freedom to swap out storage implementations without any hassel. Also allows DI and mocking.

## Usage

The store intereface defines the following methods

```
DeleteKey(key string) error
GetString(key string) (string, error)
GetInt64(key string) (int64, error)
Set(key string, value interface{}) error
SetHash(key string, hash string, value interface{}) error
DeleteHash(key string, hash string) error
GetHashString(key string, hash string) (string, error)
GetAllHashValues(key string) ([]string, error)
GetAllHashKeys(key string) ([]string, error)
SetExpiry(key string, seconds int) error
Increment(key string) error
Decrement(key string) error
SetAdd(key string, value interface{}) error
GetSetStringMembers(key string) ([]string, error)
SetRemove(key string, value interface{}) error
SetIsMember(key string, value interface{}) (bool, error)
PushItemToList(key string, value interface{}, atEnd bool) error
PopItemFromList(key string, dataType int, atEnd bool) (interface{}, error)
ItemsFromList(key string, dataType int, start, end int) (interface{}, error)
RemoveItemFromList(key string, count int, value interface{}) error
LengthOfList(key string) (int, error)
ClearDataStore()
```
