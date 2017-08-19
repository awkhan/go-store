package store

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/awkhan/go-utility/utility"
	"github.com/stretchr/testify/assert"
)

var rs *Redis

func TestMain(m *testing.M) {
	maxIdle, _ := strconv.ParseInt(os.Getenv("REDIS_MAX_IDLE"), 10, 0)
	timeout, _ := strconv.ParseInt(os.Getenv("REDIS_IDLE_TIMEOUT"), 10, 0)
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	rs = NewRedisStore(int(maxIdle), int(timeout), host, port, password)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestRedisExpiry(t *testing.T) {
	rs.ClearDataStore()

	err := rs.Set("key", "value")
	assert.Nil(t, err, "Error setting value %v", err)

	err = rs.SetExpiry("key", 2)
	assert.Nil(t, err, "Error setting expiry %v", err)

	time.Sleep(3 * time.Second)

	_, err = rs.GetString("key")
	assert.NotNil(t, err, "Error not found for invalid key")
}

func TestRedisGetSetDelString(t *testing.T) {
	rs.ClearDataStore()

	_, err := rs.GetString("invalid_key")
	assert.NotNil(t, err, "Error should not be empty fetching invalid key")

	k := "key"
	v := "val"

	assert.Nil(t, rs.Set(k, v))

	fv, err := rs.GetString(k)
	assert.Nil(t, err, "Error fetching string %v", err)
	assert.Equal(t, v, fv, "Invalid fetched value")

	assert.Nil(t, rs.DeleteKey(k))
}

func TestRedisGetSetDelInt(t *testing.T) {
	rs.ClearDataStore()

	_, err := rs.GetInt64("invalid_key")
	assert.NotNil(t, err, "Error should not be empty fetching invalid key")

	k := "key"
	v := int64(123)

	assert.Nil(t, rs.Set(k, v))

	fv, err := rs.GetInt64(k)
	assert.Nil(t, err, "Error fetching int64 %v", err)
	assert.Equal(t, v, fv, "Invalid fetched value")

	assert.Nil(t, rs.DeleteKey(k))
}

func TestRedisIncrDecr(t *testing.T) {
	rs.ClearDataStore()

	k := "key"
	ev := int64(3)

	assert.Nil(t, rs.Increment(k), "Error incrementing key")
	assert.Nil(t, rs.Increment(k), "Error incrementing key")
	assert.Nil(t, rs.Increment(k), "Error incrementing key")

	v, err := rs.GetInt64(k)
	assert.Nil(t, err, "Error getting incremented key")
	assert.Equal(t, ev, v, "Increment count is invalid")

	ev = int64(1)

	assert.Nil(t, rs.Decrement(k), "Error incrementing key")
	assert.Nil(t, rs.Decrement(k), "Error incrementing key")

	v, err = rs.GetInt64(k)
	assert.Nil(t, err, "Error getting decremented key")
	assert.Equal(t, ev, v, "Decremented count is invalid")

}

func TestRedisHash(t *testing.T) {
	rs.ClearDataStore()

	k := "key"
	hk := "hash.key"
	v := "val"

	assert.Nil(t, rs.SetHash(k, hk, v), "Error setting hash key")

	fv, err := rs.GetHashString(k, hk)

	assert.Nil(t, err, "Error fetching hash value")
	assert.Equal(t, v, fv, "Incorrect fetched value")

	assert.Nil(t, rs.DeleteHash(k, hk), "Error deleting hash value")

	fv, _ = rs.GetHashString(k, hk)
	assert.Equal(t, "", fv, "Fetched value should be empty")

	key := "new_key"
	hkeys := []string{"a", "b", "c"}
	hvalues := []string{"11", "22", "33"}

	for idx, v := range hkeys {
		assert.Nil(t, rs.SetHash(key, v, hvalues[idx]), "Error setting hash key value")
	}

	rkeys, err := rs.GetAllHashKeys(key)
	assert.Nil(t, err, "Error getting all hash keys")
	assert.Equal(t, len(hkeys), len(rkeys), "Hash keys lenght not equal to retrieved key length")

	for _, v := range rkeys {
		assert.True(t, utility.SliceContainsString(rkeys, v), "Slice %s does not contain key %s", rkeys, v)
	}

	rvalues, err := rs.GetAllHashValues(key)
	assert.Nil(t, err, "Error getting all hash values")
	assert.Equal(t, len(hvalues), len(rvalues), "Hash value length not equal to retrieved value length")

	for _, v := range rvalues {
		assert.True(t, utility.SliceContainsString(rvalues, v), "Slice %s does not contain value %s", rkeys, v)
	}
}

func TestRedisListString(t *testing.T) {
	rs.ClearDataStore()

	items := []string{"abcdefg", "ajsdfjalsdfasdf", "asdfasdfasdf", "adsfasdfasfasdfa"}
	k := "lkey"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)

	}

	l, err := rs.LengthOfList(k)
	assert.Nil(t, err, "Error retrieving lenght of listt %v", err)

	assert.Equal(t, len(items), l, "Lenght of items incorrect")

	lastItem := items[len(items)-1]
	items = items[:len(items)-1]

	item, err := rs.PopItemFromList(k, DataTypeString, true)
	assert.Nil(t, err, "Error popping item from list %v", err)

	switch ty := item.(type) {
	case string:
		assert.Equal(t, lastItem, ty, "Popped item not equal to last item")
	default:
		assert.Fail(t, "Invalid type of popped item %T", ty)
	}

	firstItem := items[0]
	items = items[1:len(items)]

	item, err = rs.PopItemFromList(k, DataTypeString, false)
	assert.Nil(t, err, "error popping item from list %v", err)

	switch ty := item.(type) {
	case string:
		assert.Equal(t, firstItem, ty, "Popped item not equal to first item")
	default:
		assert.Fail(t, "Invalid type of popped item %T", ty)
	}
}

func TestRedisListRange(t *testing.T) {
	rs.ClearDataStore()

	items := []string{"abcdefg", "ajsdfjalsdfasdf", "asdfasdfasdf", "adsfasdfasfasdfa"}
	k := "lkey"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)
	}

	fi, err := rs.ItemsFromList(k, DataTypeString, 0, 2)
	assert.Nil(t, err, "Error retrieving items from list %v", err)

	switch fetchedItems := fi.(type) {
	case []string:
		assert.Equal(t, 3, len(fetchedItems), "Invalid number of fetched items")
		assert.Equal(t, items[0], fetchedItems[0], "Invalid item at index 0")
		assert.Equal(t, items[1], fetchedItems[1], "Invalid item at index 1")
		assert.Equal(t, items[2], fetchedItems[2], "Invalid item at index 2")

	default:
		assert.Fail(t, "Invalid type of fetched items %T", fetchedItems)
	}
}

func TestRedisListRemoveValue(t *testing.T) {
	rs.ClearDataStore()

	items := []string{"abcdefg", "abcdefg", "asdfasdfasdf", "adsfasdfasfasdfa"}
	k := "lkey"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)
	}

	err := rs.RemoveItemFromList(k, 0, "abcdefg")
	assert.Nil(t, err, "Error removing item from list %v", err)

	fi, err := rs.ItemsFromList(k, DataTypeString, 0, 2)
	assert.Nil(t, err, "Error retrieving items from list %v", err)

	items = []string{items[2], items[3]}
	switch fetchedItems := fi.(type) {
	case []string:
		assert.Equal(t, 2, len(fetchedItems), "Invalid number of fetched items")
		assert.Equal(t, items[0], fetchedItems[0], "Invalid item at index 0")
		assert.Equal(t, items[1], fetchedItems[1], "Invalid item at index 1")

	default:
		assert.Fail(t, "Invalid type of fetched items %T", fetchedItems)
	}
}

func TestRedisListInt(t *testing.T) {
	rs.ClearDataStore()

	items := []int{1, 2, 4}
	k := "key"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)
	}

	fi, err := rs.ItemsFromList(k, DataTypeInt, 0, 3)
	assert.Nil(t, err, "Error retrieving items from list %v", err)

	switch fetchedItems := fi.(type) {
	case []int:
		assert.Equal(t, 3, len(fetchedItems), "Invalid number of fetched items")

		assert.Equal(t, items[0], fetchedItems[0], "Invalid item at index 0")

		assert.Equal(t, items[1], fetchedItems[1], "Invalid item at index 1")

	default:
		assert.Fail(t, "Invalid type of fetched items %T", fetchedItems)
	}
}

func TestRedisPopBool(t *testing.T) {
	rs.ClearDataStore()

	items := []bool{true, false, true}
	k := "key"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)
	}

	i, err := rs.PopItemFromList(k, DataTypeBool, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[2], i, "Last item should be true")

	i, err = rs.PopItemFromList(k, DataTypeBool, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[1], i, "Middle item should be false")

	i, err = rs.PopItemFromList(k, DataTypeBool, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[0], i, "First item should be true")
}

func TestRedisPopInt(t *testing.T) {
	rs.ClearDataStore()

	items := []int{1, 2, 3}
	k := "key"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)
	}

	i, err := rs.PopItemFromList(k, DataTypeInt, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[2], i, "Last item should be 3")

	i, err = rs.PopItemFromList(k, DataTypeInt, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[1], i, "Middle item should be 2")

	i, err = rs.PopItemFromList(k, DataTypeInt, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[0], i, "First item should be 1")
}

func TestRedisPopInt64(t *testing.T) {
	rs.ClearDataStore()

	items := []int64{1, 2, 3}
	k := "key"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)
	}

	i, err := rs.PopItemFromList(k, DataTypeInt64, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[2], i, "Last item should be 3")

	i, err = rs.PopItemFromList(k, DataTypeInt64, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[1], i, "Middle item should be 2")

	i, err = rs.PopItemFromList(k, DataTypeInt64, true)
	assert.Nil(t, err, "Error popping item from list %v", err)
	assert.Equal(t, items[0], i, "First item should be 1")
}

func TestRedisListInvalidType(t *testing.T) {
	rs.ClearDataStore()

	items := []int{1, 2, 4}
	k := "key"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)
	}

	_, err := rs.ItemsFromList(k, 23123123, 0, 3)
	assert.NotNil(t, err, "There should have been an error retrieving an invalid data type")
}

func TestRedisPopInvalidType(t *testing.T) {
	rs.ClearDataStore()

	items := []int{1}
	k := "key"

	for _, i := range items {
		err := rs.PushItemToList(k, i, true)
		assert.Nil(t, err, "Error pushing item to list %v", err)
	}

	_, err := rs.PopItemFromList(k, 23123123, true)
	assert.NotNil(t, err, "There should have been an error retrieving an invalid data type")
}

func TestRedisSet(t *testing.T) {
	rs.ClearDataStore()

	k := "key"
	sv := []string{"a", "b", "c"}

	for _, v := range sv {
		assert.Nil(t, rs.SetAdd(k, v), "Error adding item to set")
	}

	b, err := rs.SetIsMember(k, sv[0])
	assert.Nil(t, err, "Error getting set value")
	assert.True(t, b, "Value should be a member of set %s", sv)

	b, err = rs.SetIsMember(k, "123123123123123")
	assert.Nil(t, err, "Error getting set value")
	assert.False(t, b, "Value should not be a member of set %s", sv)

	members, err := rs.GetSetStringMembers(k)
	assert.Nil(t, err, "Error getting set members")
	assert.Equal(t, len(sv), len(members), "Retrieved set member lenght not equal to expected set member length")

	for _, v := range members {
		assert.True(t, utility.SliceContainsString(sv, v), "Set %s does not contain member %s", sv, v)
	}

	assert.Nil(t, rs.SetRemove(k, sv[0]), "Error removing item from set")

	sv = []string{"b", "c"}

	members, err = rs.GetSetStringMembers(k)
	assert.Nil(t, err, "Error getting set members")
	assert.Equal(t, len(sv), len(members), "Retrieved set member lenght not equal to expected set member length")

	for _, v := range members {
		assert.True(t, utility.SliceContainsString(sv, v), "Set %s does not contain member %s", sv, v)
	}
}
