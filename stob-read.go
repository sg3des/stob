package stob

import (
	"fmt"
	"io"
	"math"
	"reflect"
)

func (s *Struct) Read(p []byte) (n int, err error) {
	for _, f := range s.fields {
		if f.l+n >= len(p) {
			return n, io.ErrUnexpectedEOF
		}

		nr := f.read(p[n:])
		if f.err != nil {
			return n, f.err
		}

		// log.Println(f.rsf.Name, n, n+nr, fmt.Sprintf(" : % 02x", p[n:n+nr]))
		n += nr
	}

	return
}

type fieldReader func(p []byte) int

func (f *field) setReader() (err error) {
	switch f.rk {
	case reflect.String:
		f.read = f.String

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.read = f.Int

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f.read = f.Uint

	case reflect.Uint8:
		f.read = f.Byte

	case reflect.Bool:
		f.read = f.Bool

	case reflect.Float32:
		f.read = f.Float32
	case reflect.Float64:
		f.read = f.Float64

	case reflect.Slice:

		switch f.rv.Interface().(type) {
		case []string:
			f.read = f.SliceString
		case []int, []int8, []int16, []int32, []int64:
			f.read = f.SliceInt
		case []uint, []uint16, []uint32, []uint64:
			f.read = f.SliceUint
		case []byte:
			f.read = f.Bytes
		case []bool:
			f.read = f.SliceBool
		}

	case reflect.Array:

		switch f.rv.Index(0).Interface().(type) {
		case string:
			f.read = f.SliceString
		case int, int8, int16, int32, int64:
			f.read = f.SliceInt
		case uint, uint16, uint32, uint64:
			f.read = f.SliceUint
		case byte:
			f.read = f.Bytes
		case bool:
			f.read = f.SliceBool
		}

	case reflect.Struct:
		f.s, err = newStruct(f.rv)
		f.read = f.Struct

	case reflect.Ptr:
		f.s, err = newStruct(f.rv.Elem())
		f.read = f.Struct

	default:
		err = fmt.Errorf("Unknown field type, %s:%T", f.rsf.Name, f.rv.Interface())
	}

	return
}

//
//string

func putString(p []byte, s []byte, l int) int {
	if l == 0 {
		s = append(s, 0x00)
		l = len(s)
	} else if len(s) < l {
		s = append(s, make([]byte, l-len(s))...)
	}

	for i := 0; i < l; i++ {
		p[i] = s[i]
	}

	return l
}

func (f *field) String(p []byte) int {
	return putString(p, []byte(f.rv.String()), f.l)
}

func (f *field) SliceString(p []byte) (n int) {
	for i := 0; i < f.rv.Len(); i++ {
		n += putString(p[n:], []byte(f.rv.Index(i).String()), f.l)
	}
	return n
}

//
//int

func (f *field) Int(p []byte) int {
	Itob(p[:f.l], f.rv.Int(), f.e)
	return f.l
}

func (f *field) SliceInt(p []byte) (n int) {
	for i := 0; i < f.rv.Len(); i++ {
		Itob(p[n:n+f.l], f.rv.Index(i).Int(), f.e)
		n += f.l
	}
	return n
}

//
//uint

func (f *field) Uint(p []byte) int {
	Itob(p[:f.l], int64(f.rv.Uint()), f.e)
	return f.l
}

func (f *field) SliceUint(p []byte) (n int) {
	for i := 0; i < f.rv.Len(); i++ {
		Itob(p[n:n+f.l], int64(f.rv.Index(i).Uint()), f.e)
		n += f.l
	}
	return n
}

//
//byte

func (f *field) Byte(p []byte) int {
	p[0] = byte(f.rv.Uint())
	return 1
}

func (f *field) Bytes(p []byte) int {
	var i int
	for i = 0; i < f.l; i++ {
		p[i] = f.rv.Index(i).Interface().(byte)
	}
	return i
}

//
//bool

func (f *field) Bool(p []byte) int {
	if f.rv.Bool() {
		p[0] = 0x01
	} else {
		p[0] = 0x00
	}

	return 1
}

func (f *field) SliceBool(p []byte) int {
	count := f.rv.Len()

	for i := 0; i < count; i++ {
		if f.rv.Index(i).Bool() {
			p[i] = 0x01
		} else {
			p[i] = 0x00
		}
	}

	return count
}

// float32

func (f *field) Float32(p []byte) int {
	uf := math.Float32bits(float32(f.rv.Float()))
	Itob(p[:f.l], int64(uf), f.e)
	return f.l
}

func (f *field) SliceFloat32(p []byte) (n int) {
	for i := 0; i < f.rv.Len(); i++ {
		uf := math.Float32bits(float32(f.rv.Index(i).Float()))
		Itob(p[n:f.l], int64(uf), f.e)
		n += f.l
	}

	return
}

//
//float64

func (f *field) Float64(p []byte) int {
	uf := math.Float64bits(f.rv.Float())
	Itob(p[:f.l], int64(uf), f.e)
	return f.l
}

func (f *field) SliceFloat64(p []byte) (n int) {
	for i := 0; i < f.rv.Len(); i++ {
		uf := math.Float64bits(f.rv.Float())
		Itob(p[n:f.l], int64(uf), f.e)
		n += f.l
	}

	return
}

//
//struct

func (f *field) Struct(p []byte) (n int) {
	for _, subf := range f.s.fields {
		n += subf.read(p[n:])
	}

	return n
}

func Itob(p []byte, x int64, e Endian) {
	l := len(p)

	switch e {
	case BigEndian:
		for i := range p {
			p[i] = byte(x >> (uint(l-i-1) * 8))
		}
	case LittleEndian:
		for i := range p {
			p[i] = byte(x >> uint(i*8))
		}
	}
}