package stob

import (
	"encoding/binary"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	DefaultEndian Endian = "le"
	LittleEndian  Endian = "le"
	BigEndian     Endian = "be"

	TRUE  = []byte{0x01}
	FALSE = []byte{0x00}
)

type Endian string

func Write(w io.Writer, i interface{}) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		l, e := readTag(t.Field(i))
		if l == 0 && (f.Kind() == reflect.Array || f.Kind() == reflect.Slice) {
			l = f.Len()
		}

		switch f.Interface().(type) {
		case int:
			switch e {
			case BigEndian:
				binary.Write(w, binary.BigEndian, f.Int())
			case LittleEndian:
				binary.Write(w, binary.LittleEndian, f.Int())
			}

		case float32:
			switch e {
			case BigEndian:
				binary.Write(w, binary.BigEndian, float32(f.Float()))
			case LittleEndian:
				binary.Write(w, binary.LittleEndian, float32(f.Float()))
			}

		case float64:
			switch e {
			case BigEndian:
				binary.Write(w, binary.BigEndian, f.Float())
			case LittleEndian:
				binary.Write(w, binary.LittleEndian, f.Float())
			}

		case string:
			w.Write([]byte(f.String()))
			w.Write(FALSE)

		case byte:
			w.Write([]byte{byte(f.Uint())})

		case []byte:
			w.Write(getBytes(f, l))

		case bool:
			if v.Field(i).Bool() {
				w.Write(TRUE)
			} else {
				w.Write(FALSE)
			}

		case time.Time:
			b, _ := f.Interface().(time.Time).MarshalBinary()
			w.Write(b)

		default:

			if l == 0 {
				log.Println("WARNING! for field %s сould not determine length", f.Type().Field(i).Name)
				continue
			}

			if f.Kind() == reflect.Array {
				switch f.Index(0).Interface().(type) {
				case byte:
					w.Write(getBytes(f, l))
				}
			}

		}

	}
}

func getBytes(f reflect.Value, l int) (data []byte) {
	data = make([]byte, l)

	for i := 0; i < l; i++ {
		data[i] = f.Index(i).Interface().(byte)
	}

	return
}

func readTag(sf reflect.StructField) (length int, endian Endian) {
	ss := strings.Split(sf.Tag.Get("stob"), ",")
	if len(ss) > 0 {
		length, _ = strconv.Atoi(ss[0])
	}
	if len(ss) > 1 {
		endian = Endian(ss[1])
	} else {
		endian = DefaultEndian
	}
	return
}

func Read(r io.Reader, i interface{}) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		l, e := readTag(t.Field(i))
		if l == 0 && (f.Kind() == reflect.Array || f.Kind() == reflect.Slice) {
			l = f.Len()
		}

		switch f.Interface().(type) {
		case int:
			b := make([]byte, unsafe.Sizeof(i))
			r.Read(b)
			f.SetInt(int64(Btoi(b, e)))

		case int32:
			b := make([]byte, 4)
			r.Read(b)
			f.Set(reflect.ValueOf(int32(Btoi(b, e))))

		case int64:
			b := make([]byte, 8)
			r.Read(b)
			f.Set(reflect.ValueOf(int64(Btoi(b, e))))

		case string:
			f.SetString(readString(r))

		case byte:
			b := make([]byte, 1)
			r.Read(b)
			f.Set(reflect.ValueOf(b[0]))

		case []byte:
			l, _ := readTag(v.Type().Field(i))
			if l == 0 {
				log.Println("WARNING! for field %s сould not determine length", f.Type().Field(i).Name)
				continue
			}
			b := make([]byte, l)
			r.Read(b)
			f.SetBytes(b)

		case float32:
			var float float32
			switch e {
			case BigEndian:
				binary.Read(r, binary.BigEndian, &float)
			case LittleEndian:
				binary.Read(r, binary.LittleEndian, &float)
			}
			f.SetFloat(float64(float))

		case float64:
			var float float64
			switch e {
			case BigEndian:
				binary.Read(r, binary.BigEndian, &float)
			case LittleEndian:
				binary.Read(r, binary.LittleEndian, &float)
			}
			f.SetFloat(float)

		case bool:
			var b = make([]byte, 1)
			r.Read(b)
			if b[0] == 0x01 {
				f.SetBool(true)
			} else {
				f.SetBool(false)
			}
		case time.Time:
			var t time.Time
			var b = make([]byte, 15)
			r.Read(b)
			t.UnmarshalBinary(b)
			f.Set(reflect.ValueOf(t))

		default:

			if l == 0 {
				log.Println("WARNING! for field %s сould not determine length", f.Type().Field(i).Name)
				continue
			}

			if f.Kind() == reflect.Array {
				switch f.Index(0).Interface().(type) {
				case byte:
					bs := make([]byte, l)
					r.Read(bs)

					for i, b := range bs {
						f.Index(i).Set(reflect.ValueOf(b))
					}

				}
			}
		}

	}
}

func Btoi(b []byte, endian Endian) (a int) {
	l := len(b)
	switch endian {
	case BigEndian:
		for i := range b {
			a |= int(b[i]) << uint((l-i-1)*8)
		}
	case LittleEndian:
		for i := range b {
			a |= int(b[i]) << uint(i*8)
		}
	}

	return
}

func readString(r io.Reader) string {
	var s []byte
	for {
		b := make([]byte, 1)
		_, err := r.Read(b)
		if err != nil {
			break
		}
		if b[0] == 0x00 {
			break
		}
		s = append(s, b[0])
	}
	return string(s)
}
