package stob_unsafe

import (
	"unsafe"
)

type YourStruct struct {
	Str    string
	Int    int `stob:",le"`
	Byte   byte
	Bytes  []byte `stob:"8"`
	Bytes4 [4]byte
	Bool   bool
	Float  float32 `stob:",be"`
	Uint16 uint16
}

// func main() {
// 	a := Struct{
// 		Int:   12,
// 		Uint:  54000,
// 		Bytes: [4]byte{1, 2, 3, 4},
// 	}
// 	data := Marshal(a)
// 	fmt.Println(hex.Dump(data))

// 	b := (*Struct)(Unmarshal(data))
// 	fmt.Printf("%#v", b)
// }

func Marshal(p *YourStruct) []byte {
	return (*(*[1<<31 - 1]byte)(unsafe.Pointer(&p)))[:unsafe.Sizeof(&p)]
}

func Unmarshal(data []byte) unsafe.Pointer {
	return unsafe.Pointer(&data[0])
}
