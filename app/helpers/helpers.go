package helpers

import (
	"encoding/json"
)

func Bytes(n int, args []string) [][]byte {
	bs := [][]byte{}
	if len(args) > n {
		for _, arg := range args[n:] {
			bs = append(bs, []byte(arg))
		}
	}
	return bs
}

func MustMarshal(v interface{}) []byte {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bs
}

func MustMarshalPretty(v interface{}) []byte {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return bs
}
