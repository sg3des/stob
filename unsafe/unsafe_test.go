package stob

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

type Struct struct {
	String  string `size:"10"`
	Int     int
	Uint    uint
	Float32 float32
	Bytes4  [4]byte
	Bytes   []byte `size:"6"`
	IP      net.IP `size:"16"`
	Time    time.Duration
	SubA    *A
}

type A struct {
	AString string
	// BString string
}

func init() {
	log.SetFlags(log.Lshortfile)
}

var a Struct
var u *UPrep
var p []byte

func TestRead(t *testing.T) {
	a = Struct{
		String:  "stringname",
		Int:     255,
		Uint:    65500,
		Float32: 0.12,
		Bytes4:  [4]byte{1, 2, 3, 4},
		Bytes:   []byte{10, 11, 12, 13, 14, 15},
		IP:      net.IPv4(192, 168, 0, 1),
		Time:    1e9,
		SubA:    &A{AString: "SubA"},
	}
	u = Prepare(&a)

	p = make([]byte, 196)
	n, err := u.Read(p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(n)
	fmt.Println(hex.Dump(p))
}

func TestWrite(t *testing.T) {
	var b Struct
	u := Prepare(&b)

	n, err := u.Write(p)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%02x\n", n)

	fmt.Println(net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 1})
	fmt.Println(b.String)
	fmt.Println(b.Bytes)
	fmt.Println(b.IP)
	fmt.Println(b.Time)
	fmt.Println(b.SubA.AString)
}

func BenchmarkRead(b *testing.B) {
	p := make([]byte, 256)

	for i := 0; i < b.N; i++ {
		u.Read(p)
	}

}

func BenchmarkWrite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u.Write(p)
	}
}
