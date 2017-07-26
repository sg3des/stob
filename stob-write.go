package stob

import (
	"fmt"
	"io"
	"math"
	"reflect"
)

func (s *Struct) Write(p []byte) (n int, err error) {
	for _, f := range s.fields {
		if n+f.l >= len(p) {
			return n, io.ErrUnexpectedEOF
		}

		nw := f.write(p[n:])
		if f.err != nil {
			return n, f.err
		}

		// log.Println(f.rsf.Name, n, n+nw)
		n += nw
	}

	return
}

type fieldWriter func(p []byte) int

func (f *field) setWriter() (err error) {
	switch f.rk {
	case reflect.String:
		f.write = f.SetString

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.write = f.SetInt

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f.write = f.SetUint

	case reflect.Uint8:
		f.write = f.SetByte

	case reflect.Bool:
		f.write = f.SetBool

	case reflect.Float32:
		f.write = f.SetFloat32
	case reflect.Float64:
		f.write = f.SetFloat64

	case reflect.Slice:

		switch f.rv.Interface().(type) {
		// case []string:
		// 	f.write = f.writeSliceString
		// case []int, []int8, []int16, []int32, []int64:
		// 	f.write = f.writeSliceInt
		// case []uint, []uint16, []uint32, []uint64:
		// 	f.write = f.writeSliceUint
		case []byte:
			f.write = f.SetSliceByte
		// case []bool:
		// 	f.write = f.writeSliceBool
		default:
			err = fmt.Errorf("Unknown field type, %s:%T", f.rsf.Name, f.rv.Interface())
		}

	case reflect.Array:

		switch f.rv.Index(0).Interface().(type) {
		// case string:
		// 	f.write = f.writeArrayString
		// case int, int8, int16, int32, int64:
		// 	f.write = f.writeSliceInt
		// case uint, uint16, uint32, uint64:
		// 	f.write = f.writeSliceUint
		case byte:
			f.write = f.SetArrayByte
		// case bool:
		// 	f.write = f.writeSliceBool
		default:
			err = fmt.Errorf("Unknown field type, %s:%T", f.rsf.Name, f.rv.Interface())
		}

	case reflect.Struct:
		f.s, err = newStruct(f.rv)
		f.write = f.SetStruct

	case reflect.Ptr:
		f.s, err = newStruct(f.rv.Elem())
		f.write = f.SetStruct

	default:
		err = fmt.Errorf("Unknown field type, %s:%T", f.rsf.Name, f.rv.Interface())
	}

	return
}

//
//string

func Btos(p []byte) (string, int) {
	var s []byte

	for _, b := range p {
		s = append(s, b)
		if b == 0x00 {
			break
		}
	}

	return string(s), len(s)
}

func (f *field) SetString(p []byte) int {
	if f.l != 0 {
		f.rv.SetString(string(p[:f.l]))
		return f.l
	}

	s, n := Btos(p)
	f.rv.SetString(s)

	return n
}

// func (f *field) writeArrayString(p []byte, n int) (int, error) {
// 	var err error
// 	for i := 0; i < f.rv.Len(); i++ {
// 		n, err = writeString(p, n, f.rv.Index(i))
// 		if err != nil {
// 			return n, err
// 		}
// 	}

// 	return n, nil
// }

//
//int

func (f *field) SetInt(p []byte) int {
	f.rv.SetInt(Btoi(p[:f.l], f.e))
	return f.l
}

func (f *field) SetArrayInt(p []byte) int {
	var count = f.rv.Len()
	var limit = count * f.l

	var xx []int64
	for i := 0; i < f.rv.Len(); i++ {
		xx = append(xx, Btoi(p[:f.l], f.e))
	}

	f.rv.Set(reflect.ValueOf(xx))

	return limit
}

//
//uint

func (f *field) SetUint(p []byte) int {
	f.rv.SetUint(uint64(Btoi(p[:f.l], f.e)))
	return f.l
}

//
//byte

func (f *field) SetByte(p []byte) int {
	f.rv.SetUint(uint64(p[0]))
	return 1
}

func (f *field) SetSliceByte(p []byte) int {
	f.rv.SetBytes(p[:f.l])
	return f.l
}

func (f *field) SetArrayByte(p []byte) int {
	for i := 0; i < f.l; i++ {
		f.rv.Index(i).Set(reflect.ValueOf(p[i]))
	}
	return f.l
}

//
//bool

func (f *field) SetBool(p []byte) int {
	if p[0] != 0x00 {
		f.rv.SetBool(true)
	}
	return 1
}

//
//float32

func (f *field) SetFloat32(p []byte) int {
	x := Btoi(p[:f.l], f.e)
	float := math.Float32frombits(uint32(x))
	f.rv.SetFloat(float64(float))

	return f.l
}

//
//float64

func (f *field) SetFloat64(p []byte) int {
	x := Btoi(p[:f.l], f.e)
	float := math.Float64frombits(uint64(x))
	f.rv.SetFloat(float)

	return f.l
}

//
//struct
func (f *field) SetStruct(p []byte) (n int) {
	for _, subf := range f.s.fields {
		n += subf.write(p[n:])
	}
	return n
}

//

func Btoi(p []byte, e Endian) (x int64) {
	l := len(p)
	switch e {
	case BigEndian:
		for i := range p {
			x |= int64(p[i]) << uint((l-i-1)*8)
		}
	case LittleEndian:
		for i := range p {
			x |= int64(p[i]) << uint(i*8)
		}
	}

	return
}
