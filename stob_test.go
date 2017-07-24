package stob

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

type YourStruct struct {
	Str    string
	Int    int `stob:",le"`
	Byte   byte
	Bytes  []byte `stob:"8"`
	Bytes4 [4]byte
	Bool   bool
	Float  float32 `stob:",be"`
}

type TestStruct struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int `stob:",le"`
	Byte     byte
	Bytes    []byte `stob:"8"`
	Bytes4   [4]byte
	Spouse   bool
	Money    float32 `stob:",be"`
}

var a TestStruct
var b TestStruct

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestWriteRead(t *testing.T) {
	a := YourStruct{
		Str:    "string",
		Int:    999,
		Byte:   255,
		Bytes:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Bytes4: [4]byte{10, 11, 12, 13},
		Bool:   rand.Intn(2) == 1,
		Float:  rand.Float32(),
	}

	// a = TestStruct{
	// 	Name:     "mynameismynameis",
	// 	BirthDay: time.Now(),
	// 	Phone:    "9991112233",
	// 	Siblings: 999,
	// 	Byte:     255,
	// 	Bytes:    []byte{1, 2, 3, 4, 5, 6, 7, 8},
	// 	Bytes4:   [4]byte{10, 11, 12, 13},
	// 	Spouse:   rand.Intn(2) == 1,
	// 	Money:    rand.Float32(),
	// }

	buf := bytes.NewBuffer([]byte{})
	Write(buf, a)

	fmt.Println(hex.Dump(buf.Bytes()))

	var b YourStruct
	Read(buf, &b)

	fmt.Printf("%+v", b)
}

func BenchmarkWriteRead(b *testing.B) {
	b.StopTimer()
	buf := bytes.NewBuffer([]byte{})
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Write(buf, a)
		Read(buf, &b)

		b.StopTimer()
		buf.Reset()
		b.StartTimer()
	}
}
