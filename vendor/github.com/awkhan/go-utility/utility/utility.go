package utility

import "encoding/json"

//JSONBytes creates JSON bytes from the provided interface
func JSONBytes(data interface{}) []byte {
	res, _ := json.Marshal(data)
	return res
}

//SliceContainsString returns true if the given slice of strings contains a given string.
func SliceContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
