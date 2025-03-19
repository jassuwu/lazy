package utils

import "strings"

func ExtractArg(args []string, key string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == key {
			if i+1 < len(args) {
				return args[i+1]
			} else {
				return "" // Key found, but no value provided
			}
		} else if strings.HasPrefix(args[i], key+"=") {
			return strings.SplitN(args[i], "=", 2)[1]
		}
	}
	return "" // Key not found
}
