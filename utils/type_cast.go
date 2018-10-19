package utils

import (
	"strconv"
)

func IntToByte(data int) []byte {
	return []byte(IntToString(data))
}

func ByteToInt(data []byte) (int, error) {
	return StringToInt(string(data))
}

func StringToInt(data string) (int, error) {
	return strconv.Atoi(data)
}

func StringToInt64(data string) (int64, error) {
	return strconv.ParseInt(data, 10, 64)
}

func IntToString(data int) string {
	return strconv.Itoa(data)
}

func Int64ToString(data int64) string {
	return strconv.FormatInt(data, 10)
}
