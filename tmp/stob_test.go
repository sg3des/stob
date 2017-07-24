package stob_unsafe

import (
	"encoding/hex"
	"fmt"
	"testing"
)

type Struct struct {
	Int   int64
	Uint  uint32
	Bytes [4]byte
}

var a Struct
var b Struct
var data []byte

func init() {
	a = Struct{
		Int:   10,
		Uint:  255,
		Bytes: [4]byte{1, 2, 3, 4},
	}
}

func TestMarshal(t *testing.T) {
	data = Marshal(a)
	fmt.Println(hex.Dump(data))
}

func TestUnmarshal(t *testing.T) {
	b2 := (*(*Struct)(Unmarshal(data)))
	fmt.Printf("%+v", b2)
}
