package stob

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"testing"
)

type YourStruct struct {
	Struct     SubStruct
	PtrStruct  *SubStruct
	CustomType CustomType
	Str        string
	SliceStr   []string  `num:"2" size:"6"`
	ArrayStr   [2]string `size:"8"`
	Int        int       `bo:"le"`
	Byte       byte
	Bytes      []byte `num:"6"`
	Bytes4     [4]byte
	Bool       bool
	Float32    float32 `bo:"be"`
	Uint16     uint16
}

type SubStruct struct {
	Addr net.HardwareAddr `num:"6"`
	IP   net.IP           `num:"4"`
}

type CustomType [4]byte

func init() {
	log.SetFlags(log.Lshortfile)
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
		Struct:     SubStruct{Addr: net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, IP: []byte{127, 0, 0, 1}},
		PtrStruct:  &SubStruct{Addr: net.HardwareAddr{0x55, 0x55, 0x55, 0x55, 0x55, 0x55}},
		CustomType: [4]byte{127, 0, 0, 1},
		Str:        "string",
		Int:        999,
		Byte:       255,
		Bytes:      []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Bytes4:     [4]byte{10, 11, 12, 13},
		Bool:       rand.Intn(2) == 1,
		Float32:    0.98765,
		Uint16:     65500,
	}

	s, err := NewStruct(&a)
	if err != nil {
		t.Fatal(err)
	}

	p := make([]byte, 128)
	nr, err := s.Read(p)
	if err != nil {
		t.Error(err, nr)
	}

	fmt.Println(hex.Dump(p))

	var b YourStruct
	sb, err := NewStruct(&b)
	if err != nil {
		t.Fatal(err)
	}

	nw, err := sb.Write(p)
	if err != nil {
		t.Fatal(err, nw)
	}
	if nr != nw {
		t.Error("count readed bytes and writed bytes are different", nw, nr)
	}

	if !bytes.Equal(b.Struct.Addr, a.Struct.Addr) {
		t.Error("hw addr not equal", b.Struct.Addr, a.Struct.Addr)
	}
	if b.Str != a.Str {
		t.Error("failed read string", b.Str, a.Str)
	}
	if b.Int != a.Int {
		t.Error("failed read int", b.Int, a.Int)
	}
	if b.Byte != a.Byte {
		t.Error("failed read byte", b.Byte, a.Byte)
	}
	if !bytes.Equal(b.Bytes, a.Bytes[:6]) {
		t.Error("failed read []byte", b.Bytes, a.Bytes[:6])
	}
	if b.Bytes4 != a.Bytes4 {
		t.Error("failed read [4]byte", b.Bytes4, a.Bytes4)
	}
	if b.Bool != a.Bool {
		t.Error("failed read bool", b.Bool, a.Bool)
	}
	if b.Float32 != a.Float32 {
		t.Error("failed read float32", b.Float32, a.Float32)
	}
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

	float := 1.32
	Itob(p[8:16], int64(math.Float64bits(float)), BigEndian)

	n := Btoi(p[8:16], BigEndian)
	if rf := math.Float64frombits(uint64(n)); float != rf {
		t.Error("Failed restore float", rf, float)
	}

	// fmt.Println(hex.Dump(p))
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
