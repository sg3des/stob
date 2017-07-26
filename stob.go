package stob

import (
	"reflect"
	"strconv"
	"strings"
)

type Endian string
type Size uint8

var (
	DefaultEndian Endian = "le"
	LittleEndian  Endian = "le"
	BigEndian     Endian = "be"

	TRUE  = []byte{0x01}
	FALSE = []byte{0x00}

	Size8  Size = 1
	Size16 Size = 2
	Size32 Size = 4
	Size64 Size = 8
)

type Struct struct {
	rv reflect.Value
	rt reflect.Type

	fields []*field
}

func NewStruct(x interface{}) (*Struct, error) {
	return newStruct(reflect.ValueOf(x).Elem())
}

func newStruct(rv reflect.Value) (*Struct, error) {
	s := new(Struct)
	s.rv = rv
	s.rt = rv.Type()

	for i := 0; i < s.rv.NumField(); i++ {
		f, ok, err := newField(s.rv.Field(i), s.rt.Field(i))
		if err != nil {
			return s, err
		}
		if ok {
			s.fields = append(s.fields, f)
		}
	}

	return s, nil
}

type field struct {
	rv  reflect.Value
	rsf reflect.StructField
	rk  reflect.Kind

	l int
	e Endian

	read  fieldReader
	write fieldWriter

	s   *Struct
	err error
}

func newField(rv reflect.Value, rsf reflect.StructField) (f *field, ok bool, err error) {
	if !rv.CanSet() {
		return nil, false, nil
	}

	f = new(field)
	f.rv = rv
	f.rsf = rsf
	f.rk = rv.Kind()

	f.l, f.e, ok = f.readTag(rsf.Tag.Get("stob"))

	if f.l == 0 {
		if f.rk != reflect.String {
			f.l = int(rv.Type().Size())
		}
	}

	if err = f.setReader(); err != nil {
		return
	}

	if err = f.setWriter(); err != nil {
		return
	}

	return
}

func (f *field) readTag(tag string) (l int, e Endian, ok bool) {
	if tag == "-" {
		return
	}

	ss := strings.Split(tag, ",")
	if len(ss) > 0 {
		l, _ = strconv.Atoi(ss[0])
	}

	e = DefaultEndian
	if len(ss) > 1 {
		e = Endian(ss[1])
	}

	ok = true
	return
}

// func writeInt(w io.Writer, x interface{}, e Endian) error {
// 	switch e {
// 	case BigEndian:
// 		return binary.Write(w, binary.BigEndian, x)
// 	case LittleEndian:
// 		return binary.Write(w, binary.LittleEndian, x)
// 	}
// 	return errors.New("unknown endian")
// }

// func Write(w io.Writer, i interface{}) (err error) {
// 	v := reflect.ValueOf(i)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	t := v.Type()

// 	for i := 0; i < v.NumField(); i++ {
// 		f := v.Field(i)

// 		l, e := readTag(t.Field(i))
// 		if l == 0 {
// 			if f.Kind() == reflect.Array || f.Kind() == reflect.Slice {
// 				l = f.Len()
// 			}
// 		}

// 		switch f.Interface().(type) {
// 		case int:
// 			switch unsafe.Sizeof(i) {
// 			case 8:
// 				err = writeInt(w, int64(f.Int()), e)
// 			case 4:
// 				err = writeInt(w, int32(f.Int()), e)
// 			}

// 		case int8:
// 			err = writeInt(w, int8(f.Int()), e)
// 		case int16:
// 			err = writeInt(w, int16(f.Int()), e)
// 		case int32:
// 			err = writeInt(w, int32(f.Int()), e)
// 		case int64:
// 			err = writeInt(w, int64(f.Int()), e)

// 		case byte: //same as uint8
// 			w.Write([]byte{byte(f.Uint())})
// 		case uint16:
// 			writeInt(w, uint16(f.Uint()), e)
// 		case uint32:
// 			writeInt(w, uint32(f.Uint()), e)
// 		case uint64:
// 			writeInt(w, uint64(f.Uint()), e)
// 		case uint:
// 			writeInt(w, uint(f.Uint()), e)

// 		case float32:
// 			writeInt(w, float32(f.Float()), e)
// 		case float64:
// 			writeInt(w, f.Float(), e)

// 		case string:
// 			w.Write([]byte(f.String()))
// 			w.Write(FALSE)

// 		case []byte:
// 			w.Write(getBytes(f, l))

// 		case bool:
// 			if v.Field(i).Bool() {
// 				w.Write(TRUE)
// 			} else {
// 				w.Write(FALSE)
// 			}

// 		case time.Time:
// 			b, _ := f.Interface().(time.Time).MarshalBinary()
// 			w.Write(b)

// 		default:

// 			if l == 0 {
// 				log.Printf("WARNING! for field %s сould not determine length\n", f.Type().Field(i).Name)
// 				continue
// 			}

// 			if f.Kind() == reflect.Array {
// 				switch f.Index(0).Interface().(type) {
// 				case byte:
// 					w.Write(getBytes(f, l))
// 				}
// 			}

// 		}

// 		if err != nil {
// 			return
// 		}

// 	}

// 	return
// }

// func getBytes(f reflect.Value, l int) (data []byte) {
// 	data = make([]byte, l)

// 	for i := 0; i < l; i++ {
// 		data[i] = f.Index(i).Interface().(byte)
// 	}

// 	return
// }

// func readInt(r io.Reader, s Size, e Endian) int64 {
// 	b := make([]byte, s)
// 	r.Read(b)
// 	return Btoi(b, e)
// }

// func Read(r io.Reader, i interface{}) {
// 	v := reflect.ValueOf(i)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	t := v.Type()

// 	for i := 0; i < v.NumField(); i++ {
// 		f := v.Field(i)

// 		l, e := readTag(t.Field(i))
// 		if l == 0 {
// 			if f.Kind() == reflect.Array || f.Kind() == reflect.Slice {
// 				l = f.Len()
// 			}
// 		}

// 		switch f.Interface().(type) {
// 		case int:
// 			f.SetInt(readInt(r, Size(unsafe.Sizeof(i)), e))
// 		case int8:
// 			f.SetInt(readInt(r, Size8, e))
// 		case int16:
// 			f.SetInt(readInt(r, Size16, e))
// 		case int32:
// 			f.SetInt(readInt(r, Size32, e))
// 		case int64:
// 			f.SetInt(readInt(r, Size64, e))

// 		case uint8: //same as byte
// 			f.SetUint(uint64(readInt(r, Size8, e)))
// 		case uint16:
// 			f.SetUint(uint64(readInt(r, Size16, e)))
// 		case uint32:
// 			f.SetUint(uint64(readInt(r, Size32, e)))
// 		case uint64:
// 			f.SetUint(uint64(readInt(r, Size64, e)))

// 		case string:
// 			f.SetString(readString(r))

// 		case []byte:
// 			l, _ := readTag(v.Type().Field(i))
// 			if l == 0 {
// 				log.Printf("WARNING! for field %s сould not determine length\n", f.Type().Field(i).Name)
// 				continue
// 			}
// 			b := make([]byte, l)
// 			r.Read(b)
// 			f.SetBytes(b)

// 		case float32:
// 			var float float32
// 			switch e {
// 			case BigEndian:
// 				binary.Read(r, binary.BigEndian, &float)
// 			case LittleEndian:
// 				binary.Read(r, binary.LittleEndian, &float)
// 			}
// 			f.SetFloat(float64(float))

// 		case float64:
// 			var float float64
// 			switch e {
// 			case BigEndian:
// 				binary.Read(r, binary.BigEndian, &float)
// 			case LittleEndian:
// 				binary.Read(r, binary.LittleEndian, &float)
// 			}
// 			f.SetFloat(float)

// 		case bool:
// 			var b = make([]byte, 1)
// 			r.Read(b)
// 			if b[0] == 0x01 {
// 				f.SetBool(true)
// 			} else {
// 				f.SetBool(false)
// 			}
// 		case time.Time:
// 			var t time.Time
// 			var b = make([]byte, 15)
// 			r.Read(b)
// 			t.UnmarshalBinary(b)
// 			f.Set(reflect.ValueOf(t))

// 		default:

// 			if l == 0 {
// 				log.Printf("WARNING! for field %s сould not determine length", t.Field(i).Name)
// 				continue
// 			}

// 			if f.Kind() == reflect.Array {
// 				switch f.Index(0).Interface().(type) {
// 				case byte:
// 					bs := make([]byte, l)
// 					r.Read(bs)

// 					for i, b := range bs {
// 						f.Index(i).Set(reflect.ValueOf(b))
// 					}

// 				}
// 			}
// 		}

// 	}
// }

// func readString(r io.Reader) string {
// 	var s []byte
// 	for {
// 		b := make([]byte, 1)
// 		_, err := r.Read(b)
// 		if err != nil {
// 			break
// 		}
// 		if b[0] == 0x00 {
// 			break
// 		}
// 		s = append(s, b[0])
// 	}
// 	return string(s)
// }
