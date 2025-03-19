package assert

import (
	"log"
)

func NotNil(item any, message string) {
	if item == nil {
		log.Fatalln("skill issue'd by: ", message, item)
	}
}

func Nil(item any, message string) {
	if item != nil {
		log.Fatalln("skill issue'd by: ", message, item)
	}
}
