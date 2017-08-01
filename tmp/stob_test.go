package stob_unsafe

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"unsafe"
)

var a YourStruct
var b YourStruct
var data []byte

func init() {

}

func TestMarshal(t *testing.T) {
	a := YourStruct{
		Str:    "string",
		Int:    999,
		Byte:   255,
		Bytes:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Bytes4: [4]byte{10, 11, 12, 13},
		Bool:   rand.Intn(2) == 1,
		Float:  rand.Float32(),
		Uint16: 65500,
	}
	data := (*(*[1<<31 - 1]byte)(unsafe.Pointer(&a)))[:unsafe.Sizeof(YourStruct{})]
	fmt.Println(hex.Dump(data))
}

// func TestUnmarshal(t *testing.T) {
// 	b2 := (*YourStruct)(Unmarshal(data))
// 	fmt.Printf("%+v", b2)
// }

// func BenchmarkMarshal(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		Marshal(&a)
// 	}
// }

// func BenchmarkUnmarshal(b *testing.B) {
// 	var a *YourStruct
// 	for i := 0; i < b.N; i++ {
// 		a = (*YourStruct)(Unmarshal(data))
// 	}
// 	b.Log(a)
// }
