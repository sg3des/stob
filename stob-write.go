package stob

import (
	"io"
	"math"
	"reflect"
)

func (s *Struct) Write(p []byte) (n int, err error) {
	for _, f := range s.fields {

		// log.Println(f.rsf.Name, f.len, n, len(p))

		if f.len+n > len(p) {
			return n, io.ErrUnexpectedEOF
		}

		nw, err := f.write(p[n:])
		if err != nil {
			return n, err
		}

		n += nw
	}

	return
}

type fieldWriter func(p []byte) (int, error)

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
		case []string:
			f.write = f.SetSliceString
		// case []int, []int8, []int16, []int32, []int64:
		// 	f.write = f.writeSliceInt
		// case []uint, []uint16, []uint32, []uint64:
		// 	f.write = f.writeSliceUint
		case []byte:
			// if f.len == 0 {
			// 	return fmt.Errorf("Field %s type []byte should have count nums in tags: `num:\"#\"`", f.rsf.Name)
			// }
			f.write = f.SetSliceByte
		// case []bool:
		// 	f.write = f.writeSliceBool
		default:
			f.write = f.SetCustom
			// f.rv.Set(reflect.New(f.rv.Type()).Elem())
			// log.Printf("%T %s\n", f.rv.Interface(), f.rv.Interface())
			// err = fmt.Errorf("Unknown field type, %s:%T", f.rsf.Name, f.rv.Interface())
		}

	case reflect.Array:

		switch f.rv.Index(0).Interface().(type) {
		case string:
			f.write = f.SetArrayString
		// case int, int8, int16, int32, int64:
		// 	f.write = f.writeSliceInt
		// case uint, uint16, uint32, uint64:
		// 	f.write = f.writeSliceUint
		case byte:
			f.write = f.SetArrayByte
			// case bool:
			// 	f.write = f.writeSliceBool
			// default:
			// err = fmt.Errorf("Unknown field type, %s:%T", f.rsf.Name, f.rv.Interface())
		default:
			f.write = f.SetCustom
		}

	case reflect.Struct:
		// f.s, err = newStruct(f.rv) //already in prepare readers
		f.write = f.SetStruct

	case reflect.Ptr:
		// f.s, err = newStruct(f.rv.Elem()) // aready on preapre readers
		f.write = f.SetStruct

	default:
		f.write = f.SetCustom
		// err = fmt.Errorf("Unknown field type, %s:%T", f.rsf.Name, f.rv.Interface())
	}

	return
}

//
//string

//Btos = bytes to string, read byte array to first 0x00 byte, then return string and count of readed bytes.
func Btos(p []byte) (_ string, n int) {
	var s []byte
	var b byte

	for n, b = range p {
		if b == 0x00 {
			n++
			break
		}
		s = append(s, b)
	}

	return string(s), n
}

func (f *field) SetString(p []byte) (n int, _ error) {
	var s string

	if f.size != 0 {
		s, _ = Btos(p[:f.size])
		n = f.size
	} else {
		s, n = Btos(p)
	}
	f.rv.SetString(s)

	return
}

func (f *field) SetSliceString(p []byte) (n int, err error) {
	var ss []string
	var s string
	var ns int

	for {
		if f.num != 0 && len(ss) >= f.num {
			break
		}

		if n >= len(p) {
			break
		}

		if f.size != 0 {
			if n+f.size >= len(p) {
				return n, io.ErrUnexpectedEOF
			}

			s, _ = Btos(p[n : n+f.size])
			ns = f.size
		} else {
			s, ns = Btos(p[n:])
		}

		ss = append(ss, s)
		n += ns
	}

	f.rv.Set(reflect.ValueOf(ss))

	return
}

func (f *field) SetArrayString(p []byte) (n int, err error) {
	var s string
	var ns int

	for i := 0; i < f.rv.Len(); i++ {
		if f.size != 0 {
			s, _ = Btos(p[n : n+f.size])
			n += f.size
		} else {
			s, ns = Btos(p[n:])
			n += ns
		}

		f.rv.Index(i).SetString(s)
	}

	return
}

//
//int

func (f *field) SetInt(p []byte) (int, error) {
	f.rv.SetInt(Btoi(p[:f.size], f.e))
	return f.size, nil
}

func (f *field) SetArrayInt(p []byte) (n int, err error) {
	var xx []int64
	for i := 0; i < f.rv.Len(); i++ {
		xx = append(xx, Btoi(p[n:n+f.size], f.e))
		n += f.size
	}

	f.rv.Set(reflect.ValueOf(xx))

	return
}

//
//uint

func (f *field) SetUint(p []byte) (int, error) {
	f.rv.SetUint(uint64(Btoi(p[:f.size], f.e)))
	return f.size, nil
}

//
//byte

func (f *field) SetByte(p []byte) (int, error) {
	f.rv.SetUint(uint64(p[0]))
	return 1, nil
}

func (f *field) SetSliceByte(p []byte) (int, error) {
	if f.len == 0 {
		f.len = len(p)
	}
	f.rv.SetBytes(p[:f.len])
	return f.len, nil
}

func (f *field) SetArrayByte(p []byte) (int, error) {
	for i := 0; i < f.rv.Len(); i++ {
		f.rv.Index(i).Set(reflect.ValueOf(p[i]))
	}
	return f.rv.Len(), nil
}

//
//bool

func (f *field) SetBool(p []byte) (int, error) {
	if p[0] != 0x00 {
		f.rv.SetBool(true)
	}
	return 1, nil
}

//
//float32

func (f *field) SetFloat32(p []byte) (int, error) {
	x := Btoi(p[:f.size], f.e)
	float := math.Float32frombits(uint32(x))
	f.rv.SetFloat(float64(float))

	return f.size, nil
}

//
//float64

func (f *field) SetFloat64(p []byte) (int, error) {
	x := Btoi(p[:f.size], f.e)
	float := math.Float64frombits(uint64(x))
	f.rv.SetFloat(float)

	return f.size, nil
}

//
//struct
func (f *field) SetStruct(p []byte) (n int, _ error) {
	for _, subf := range f.s.fields {
		nw, err := subf.write(p[n:])
		if err != nil {
			return n, err
		}
		n += nw
	}
	return n, nil
}

//
//custom types

func (f *field) SetCustom(p []byte) (n int, err error) {
	count := f.num
	if count == 0 {
		count = int(f.rv.Type().Size())
	}

	f.rv.Set(reflect.ValueOf(p[:count]))
	return count, nil
}

func Btoi(p []byte, e ByteOrder) (x int64) {
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
