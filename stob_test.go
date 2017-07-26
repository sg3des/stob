package stob

import (
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/rand"
	"testing"
)

type YourStruct struct {
	DstHwAddr HwAddr
	// SrcHwAddr *HwAddr
	Str     string
	Int     int `stob:",le"`
	Byte    byte
	Bytes   []byte `stob:"6"`
	Bytes4  [4]byte
	Bool    bool
	Float32 float32 `stob:",be"`
	Uint16  uint16
}

type HwAddr struct {
	Addr [6]byte
}

func init() {
	log.SetFlags(log.Lshortfile)

	p := []byte{1, 2, 3, 4, 5, 6}
	arr(p[2:4])
	fmt.Println(p)
}

func arr(p []byte) {
	p[0] = 10
}

// func TestBuffer(t *testing.T) {
// 	p := make([]byte, 16)

// 	buf := bytes.NewBuffer([]byte{})
// 	buf.Write([]byte("asdasdasdasdasd01234"))

// 	n, err := buf.Read(p)
// 	log.Println(n, err)

// 	fmt.Println(hex.Dump(p))
// }

func TestWriteRead(t *testing.T) {
	a := YourStruct{
		DstHwAddr: HwAddr{[6]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		// SrcHwAddr: &HwAddr{[6]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55}},
		Str:     "string",
		Int:     999,
		Byte:    255,
		Bytes:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Bytes4:  [4]byte{10, 11, 12, 13},
		Bool:    rand.Intn(2) == 1,
		Float32: 0.98765,
		Uint16:  65500,
	}

	s, err := NewStruct(&a)
	if err != nil {
		t.Fatal(err)
	}

	p := make([]byte, 128)
	n, err := s.Read(p)
	log.Println(n, err)

	fmt.Println(hex.Dump(p))

	var b YourStruct
	sb, err := NewStruct(&b)
	if err != nil {
		t.Fatal(err)
	}

	n, err = sb.Write(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("writed:", n)
	t.Logf("%+v", b)

	// if err := Write(buf, a); err != nil {
	// 	t.Fatal(err)
	// }

	// // fmt.Println(hex.Dump(buf.Bytes()))

	// var b YourStruct
	// Read(buf, &b)

	// if b.Str != a.Str {
	// 	t.Error("failed read string", b.Str, a.Str)
	// }
	// if b.Int != a.Int {
	// 	t.Error("failed read int", b.Int, a.Int)
	// }
	// if b.Byte != a.Byte {
	// 	t.Error("failed read byte", b.Byte, a.Byte)
	// }
	// if !bytes.Equal(b.Bytes, a.Bytes) {
	// 	t.Error("failed read []byte", b.Bytes, a.Bytes)
	// }
	// if b.Bytes4 != a.Bytes4 {
	// 	t.Error("failed read [4]byte", b.Bytes4, a.Bytes4)
	// }
	// if b.Bool != a.Bool {
	// 	t.Error("failed read bool", b.Bool, a.Bool)
	// }
	// if b.Float != a.Float {
	// 	t.Error("failed read float32", b.Float, a.Float)
	// }

	// fmt.Printf("%+v", b)
}

func TestItob(t *testing.T) {
	var p = make([]byte, 32)
	Itob(p[0:4], 255, BigEndian)
	if p[3] != 0xff {
		t.Error("BigEndian failed")
	}

	Itob(p[4:8], 255, LittleEndian)
	if p[4] != 0xff {
		t.Error("LittleEndian failed")
	}

	Itob(p[8:12], int64(math.Float64bits(1.32)), BigEndian)

	n := Btoi(p[8:12], BigEndian)
	log.Println(math.Float32frombits(uint32(n)))

	fmt.Println(hex.Dump(p))
}

func BenchmarkRead(b *testing.B) {
	b.StopTimer()
	a := YourStruct{
		Str:     "string",
		Int:     999,
		Byte:    255,
		Bytes:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Bytes4:  [4]byte{10, 11, 12, 13},
		Bool:    rand.Intn(2) == 1,
		Float32: rand.Float32(),
		Uint16:  65500,
	}
	buf := make([]byte, 128)

	s, err := NewStruct(&a)
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		s.Read(buf)
	}
}

func BenchmarkWrite(b *testing.B) {
	b.StopTimer()
	a := YourStruct{
		Str:     "string",
		Int:     999,
		Byte:    255,
		Bytes:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Bytes4:  [4]byte{10, 11, 12, 13},
		Bool:    rand.Intn(2) == 1,
		Float32: rand.Float32(),
		Uint16:  65500,
	}
	buf := make([]byte, 128)

	s, err := NewStruct(&a)
	if err != nil {
		b.Fatal(err)
	}
	s.Read(buf)

	sw, err := NewStruct(&a)
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sw.Write(buf)
	}
}

// func BenchmarkWrite(b *testing.B) {
// 	b.StopTimer()
// 	a := YourStruct{
// 		Str:    "string",
// 		Int:    999,
// 		Byte:   255,
// 		Bytes:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
// 		Bytes4: [4]byte{10, 11, 12, 13},
// 		Bool:   rand.Intn(2) == 1,
// 		Float:  rand.Float32(),
// 		Uint16: 65500,
// 	}
// 	buf := bytes.NewBuffer([]byte{})
// 	Write(buf, a)
// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		Read(buf, &a)

// 		// b.StopTimer()
// 		// buf.Reset()
// 		// b.StartTimer()
// 	}
// }

// func BenchmarkWrite(b *testing.B) {
// 	b.StopTimer()
// 	a := YourStruct{
// 		Str:    "string",
// 		Int:    999,
// 		Byte:   255,
// 		Bytes:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
// 		Bytes4: [4]byte{10, 11, 12, 13},
// 		Bool:   rand.Intn(2) == 1,
// 		Float:  rand.Float32(),
// 		Uint16: 65500,
// 	}

// 	buf := bytes.NewBuffer([]byte{})
// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		Write(buf, a)

// 		b.StopTimer()
// 		buf.Reset()
// 		b.StartTimer()
// 	}
// }

// func BenchmarkRead(b *testing.B) {
// 	b.StopTimer()
// 	a := YourStruct{
// 		Str:    "string",
// 		Int:    999,
// 		Byte:   255,
// 		Bytes:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
// 		Bytes4: [4]byte{10, 11, 12, 13},
// 		Bool:   rand.Intn(2) == 1,
// 		Float:  rand.Float32(),
// 		Uint16: 65500,
// 	}
// 	buf := bytes.NewBuffer([]byte{})
// 	Write(buf, a)
// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		Read(buf, &a)

// 		// b.StopTimer()
// 		// buf.Reset()
// 		// b.StartTimer()
// 	}
// }
