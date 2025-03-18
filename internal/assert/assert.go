package assert

import "fmt"

func NotNil(item any, message string) {
	if item == nil {
		panic(fmt.Sprintln("skill issue'd by: ", message, item))
	}
}

func Nil(item any, message string) {
	if item != nil {
		panic(fmt.Sprintln("skill issue'd by: ", message, item))
	}
}
