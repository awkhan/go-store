package store

const (
	//DataTypeString is a string data type.
	DataTypeString = iota
	//DataTypeBool is a bool data type.
	DataTypeBool
	//DataTypeInt is a int data type.
	DataTypeInt
	//DataTypeInt64 is a int64 data type.
	DataTypeInt64
)

//Store represents an interface associated with NO SQL databases
type Store interface {
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
}
