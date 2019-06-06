package helpers

import (
	"log"
	"strconv"
)

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func ConvertStringToInteger(value string) int {
	integerValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return integerValue
}
