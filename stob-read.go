package stob

import (
	"io"
	"log"
	"math"
	"reflect"
	"unsafe"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type Reader interface {
	io.Reader
	Size
}

type Size interface {
	Size() int
}

func (s *Struct) Read(p []byte) (n int, err error) {
	for _, f := range s.fields {

		// log.Println(f.rsf.Name, f.len, n, len(p))

		if f.len+n > len(p) {
			return n, io.ErrUnexpectedEOF
		}

		n += f.Read(p[n:])
	}

	return n, io.EOF
}

type fieldReader func(p []byte) int

func (f *field) setReader() (err error) {
	rt := f.rv.Type()
	if f.rv.Kind() == reflect.Ptr {
		rt = f.rv.Type().Elem()
	}

	rv := reflect.New(rt)

	if _, ok := rv.Interface().(Reader); ok {
		f.rk = f.rv.Kind()

		f.Read = func(p []byte) int {
			n, err := f.rv.Interface().(Reader).Read(p)
			if err != nil {
				return n
			}

			return n
		}

		return nil
	}

	// if _, ok := f.rv.Interface().(Reader); ok {
	// 	return nil
	// }

	switch f.rk {
	case reflect.String:
		f.Read = f.String

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.Read = f.Int

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f.Read = f.Uint

	case reflect.Uint8:
		f.Read = f.Byte

	case reflect.Bool:
		f.Read = f.Bool

	case reflect.Float32:
		f.Read = f.Float32
	case reflect.Float64:
		f.Read = f.Float64

	case reflect.Slice:

		switch f.rv.Interface().(type) {
		case []string:
			f.Read = f.SliceString
		case []int, []int8, []int16, []int32, []int64:
			f.Read = f.SliceInt
		case []uint, []uint16, []uint32, []uint64:
			f.Read = f.SliceUint
		case []byte:
			f.Read = f.Bytes
		case []bool:
			f.Read = f.SliceBool
		default:
			f.Read = f.Custom
		}

	case reflect.Array:

		switch f.rv.Index(0).Interface().(type) {
		case string:
			f.Read = f.SliceString
		case int, int8, int16, int32, int64:
			f.Read = f.SliceInt
		case uint, uint16, uint32, uint64:
			f.Read = f.SliceUint
		case byte:
			f.Read = f.Bytes
		case bool:
			f.Read = f.SliceBool
		default:
			f.Read = f.Custom
		}

	case reflect.Struct:
		f.s, err = newStruct(f.rv)
		f.Read = f.Struct

	case reflect.Ptr:
		if f.rv.IsNil() {
			f.rv.Set(reflect.New(f.rv.Type().Elem()))
		}

		f.s, err = newStruct(f.rv.Elem())
		f.Read = f.Struct

	default:
		f.Read = f.Custom
		// log.Printf("%T\n", f.rv.Interface())
		// err = fmt.Errorf("Unknown field type, %s:%T", f.rsf.Name, f.rv.Interface())
	}

	return
}

//
// string

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
	return putString(p, []byte(f.rv.String()), f.size)
}

func (f *field) SliceString(p []byte) (n int) {
	count := f.num
	if count == 0 {
		count = f.rv.Len()
	}

	for i := 0; i < count; i++ {
		if i < f.rv.Len() {
			n += putString(p[n:], []byte(f.rv.Index(i).String()), f.size)
		} else {
			n += putString(p[n:n+f.size], nil, f.size)
		}
	}
	return n
}

//
// int

func (f *field) Int(p []byte) int {
	Itob(p[:f.size], f.rv.Int(), f.e)
	return f.size
}

func (f *field) SliceInt(p []byte) (n int) {
	count := f.num
	if count == 0 {
		count = f.rv.Len()
	}

	for i := 0; i < count; i++ {
		Itob(p[n:n+f.size], f.rv.Index(i).Int(), f.e)
		n += f.size
	}

	return n
}

//
// uint

func (f *field) Uint(p []byte) int {
	Itob(p[:f.size], int64(f.rv.Uint()), f.e)
	return f.size
}

func (f *field) SliceUint(p []byte) (n int) {
	count := f.num
	if count == 0 {
		count = f.rv.Len()
	}

	for i := 0; i < count; i++ {
		Itob(p[n:n+f.size], int64(f.rv.Index(i).Uint()), f.e)
		n += f.size
	}
	return n
}

//
// byte

func (f *field) Byte(p []byte) int {
	p[0] = byte(f.rv.Uint())
	return 1
}

func (f *field) Bytes(p []byte) int {
	count := f.num
	if count == 0 {
		count = f.rv.Len()
	}

	for i := 0; i < count; i++ {
		p[i] = f.rv.Index(i).Interface().(byte)
	}

	return count
}

//
// bool

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
	Itob(p[:f.size], int64(uf), f.e)
	return f.size
}

func (f *field) SliceFloat32(p []byte) (n int) {
	count := f.num
	if count == 0 {
		count = f.rv.Len()
	}

	for i := 0; i < count; i++ {
		uf := math.Float32bits(float32(f.rv.Index(i).Float()))
		Itob(p[n:n+f.size], int64(uf), f.e)
		n += f.size
	}

	return
}

//
// float64

func (f *field) Float64(p []byte) int {
	uf := math.Float64bits(f.rv.Float())
	Itob(p[:f.size], int64(uf), f.e)
	return f.size
}

func (f *field) SliceFloat64(p []byte) (n int) {
	count := f.num
	if count == 0 {
		count = f.rv.Len()
	}

	for i := 0; i < count; i++ {
		uf := math.Float64bits(f.rv.Float())
		Itob(p[n:n+f.size], int64(uf), f.e)
		n += f.size
	}

	return
}

//
// struct

func (f *field) Struct(p []byte) (n int) {
	for _, subf := range f.s.fields {
		n += subf.Read(p[n:])
	}

	return n
}

//
// custom types

// Custom use unsafe pointer
func (f *field) Custom(p []byte) int {
	count := f.num
	if count == 0 {
		count = int(f.rv.Type().Size())
	}

	ptr := f.rv.Pointer()
	if ptr != 0 {
		for i := 0; i < count; i++ {
			p[i] = *((*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i))))
		}
	} else {
		for i := 0; i < count; i++ {
			p[i] = 0x00
		}
	}

	return count
}

// Itob convert int to bytes
func Itob(p []byte, x int64, e ByteOrder) {
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
