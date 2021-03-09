package internal

import "unsafe"

// Str2bytes convert string to bytes
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2str convert []byte to string
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
