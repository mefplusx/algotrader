package helper

import (
	"os"
	"strconv"
)

func ExistsArgs(bullets ...string) bool {
	for _, arg := range os.Args {
		for _, bullet := range bullets {
			if arg == bullet {
				return true
			}
		}
	}

	return false
}

func GetInt64FromArgs(id int, defaultVal int64) int64 {
	v := defaultVal
	if len(os.Args) >= id+1 {
		v, _ = strconv.ParseInt(os.Args[id], 10, 64)
	}
	return v
}

func GetStringFromArgs(id int, defaultVal string) string {
	v := defaultVal
	if len(os.Args) >= id+1 {
		return os.Args[id]
	}
	return v
}
