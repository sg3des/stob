package stob

import (
	"bytes"
	"log"
	"math/rand"
	"testing"
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

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestWriteRead(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

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

	if err := Write(buf, a); err != nil {
		t.Fatal(err)
	}

	// fmt.Println(hex.Dump(buf.Bytes()))

	var b YourStruct
	Read(buf, &b)

	if b.Str != a.Str {
		t.Error("failed read string", b.Str, a.Str)
	}
	if b.Int != a.Int {
		t.Error("failed read int", b.Int, a.Int)
	}
	if b.Byte != a.Byte {
		t.Error("failed read byte", b.Byte, a.Byte)
	}
	if !bytes.Equal(b.Bytes, a.Bytes) {
		t.Error("failed read []byte", b.Bytes, a.Bytes)
	}
	if b.Bytes4 != a.Bytes4 {
		t.Error("failed read [4]byte", b.Bytes4, a.Bytes4)
	}
	if b.Bool != a.Bool {
		t.Error("failed read bool", b.Bool, a.Bool)
	}
	if b.Float != a.Float {
		t.Error("failed read float32", b.Float, a.Float)
	}

	// fmt.Printf("%+v", b)
}

func BenchmarkWrite(b *testing.B) {
	b.StopTimer()
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

	buf := bytes.NewBuffer([]byte{})
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Write(buf, a)

		b.StopTimer()
		buf.Reset()
		b.StartTimer()
	}
}

func BenchmarkRead(b *testing.B) {
	b.StopTimer()
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
	buf := bytes.NewBuffer([]byte{})
	Write(buf, a)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Read(buf, &a)

		// b.StopTimer()
		// buf.Reset()
		// b.StartTimer()
	}
}
