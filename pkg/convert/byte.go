package convert

import (
	"reflect"
	"strconv"
	"unsafe"
)

func ByteAtoInt(b []byte) (int, error) {
	return strconv.Atoi(BytesToString(b))
}

func ByteToInt(b []byte, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(BytesToString(b), base, bitSize)
}

func ByteToUint(b []byte, base int, bitSize int) (uint64, error) {
	return strconv.ParseUint(BytesToString(b), base, bitSize)
}

func ByteToFloat(b []byte, bitSize int) (float64, error) {
	return strconv.ParseFloat(BytesToString(b), bitSize)
}

// BytesToString converts byte slice to string.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}