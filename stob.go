package stob

import (
	"io/ioutil"
	"reflect"
	"strconv"
)

type ByteOrder string

var (
	DefaultEndian ByteOrder = "le"
	LittleEndian  ByteOrder = "le"
	BigEndian     ByteOrder = "be"
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

	num  int
	size int
	len  int
	e    ByteOrder

	Read  fieldReader
	Write fieldWriter

	s *Struct
}

func newField(rv reflect.Value, rsf reflect.StructField) (f *field, ok bool, err error) {
	if !rv.CanSet() {
		return nil, false, nil
	}

	f = new(field)
	f.rv = rv
	f.rsf = rsf
	f.rk = rv.Kind()

	if ok = f.readTag(rsf.Tag); !ok {
		return
	}

	f.lookupSizes()

	if err = f.setReader(); err != nil {
		return
	}

	if err = f.setWriter(); err != nil {
		return
	}

	if f.rk == reflect.Struct || f.rk == reflect.Ptr {
		f.lookupStructSizes()
	}

	return
}

func (f *field) readTag(tag reflect.StructTag) bool {
	if tag.Get("stob") == "-" {
		return false
	}

	f.e = DefaultEndian
	f.e = ByteOrder(tag.Get("bo"))

	f.size, _ = strconv.Atoi(tag.Get("size"))
	f.num, _ = strconv.Atoi(tag.Get("num"))

	return true
}

func (f *field) lookupSizes() {
	if f.size == 0 {
		if f.rk != reflect.String && f.rk != reflect.Slice && f.rk != reflect.Array {
			f.size = int(f.rv.Type().Size())
		}
	}

	if f.num == 0 && f.rk == reflect.Array {
		f.num = f.rv.Len()
	}

	f.len = f.size
	if f.len == 0 {
		f.len = f.num
	}

	if f.num != 0 && f.size != 0 {
		f.len = f.num * f.size
	}
}

func (f *field) lookupStructSizes() (err error) {
	if s, ok := f.rv.Interface().(Size); ok {
		f.len = s.Size()
		return
	}

	if f.s == nil {
		f.s, err = newStruct(f.rv)
		if err != nil {
			return err
		}
	}

	var structLen int

	for _, subf := range f.s.fields {
		structLen += subf.len
	}

	f.len = structLen

	return nil
}

//
//
//

func Marshal(x interface{}) ([]byte, error) {
	s, err := NewStruct(x)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(s)
}

func Unmarshal(data []byte, x interface{}) error {
	s, err := NewStruct(x)
	if err != nil {
		return err
	}

	_, err = s.Write(data)
	return err
}
