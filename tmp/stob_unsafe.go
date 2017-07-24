package stob_unsafe

import (
	"unsafe"
)

func Marshal(i interface{}) []byte {
	return (*(*[1<<31 - 1]byte)(unsafe.Pointer(&i)))[:unsafe.Sizeof(i)]
}

func Unmarshal(data []byte) unsafe.Pointer {
	return unsafe.Pointer(&data[0])
}
