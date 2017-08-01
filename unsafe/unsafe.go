package stob

import (
	"reflect"
	"unsafe"
)

const w = int(unsafe.Sizeof(1))

const MaxBuffer = 4096

type UPrep struct {
	s interface{}
	v reflect.Value
}

func Prepare(i interface{}) *UPrep {
	return &UPrep{
		s: i,
		v: reflect.ValueOf(i).Elem(),
	}
}

func (u *UPrep) Read(p []byte) (n int, err error) {
	n, _, err = StructRead(p, u.v)
	return
}

func StructRead(p []byte, v reflect.Value) (po int, do int, err error) {
	ptr := v.UnsafeAddr()

	d := (*[MaxBuffer]byte)(unsafe.Pointer(ptr))[:]

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.Ptr:
			spo, _, err := StructRead(p[po:], f.Elem())

			po += spo
			do += w * 2

			if err != nil {
				return po, do, err
			}

		case reflect.String, reflect.Slice:
			n := *(*int)(unsafe.Pointer(&d[do+w]))
			s := *(*[]byte)(unsafe.Pointer(&d[do]))
			copy(p[po:], d[do+w:do+w+2])
			copy(p[po+2:], s)

			po += n + 2 // 2 is size of string or slice
			do += int(f.Type().Size())

		default:
			size := int(f.Type().Size())
			copy(p[po:], d[do:do+size])

			do += size
			po += size
		}

	}
	return
}

func (u *UPrep) Write(p []byte) (n int, err error) {
	n, _, err = StructWrite(p, u.v)
	return
}

func StructWrite(p []byte, v reflect.Value) (po int, do int, err error) {
	ptr := v.UnsafeAddr()
	d := (*[MaxBuffer]byte)(unsafe.Pointer(ptr))[:]

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.Ptr:
			sv := reflect.New(f.Type().Elem()).Elem()

			subptr := sv.UnsafeAddr()
			copy(d[do:], (*[8]byte)(unsafe.Pointer(&subptr))[:])

			spo, _, err := StructWrite(p[po:], sv)
			po += spo
			do += w * 2
			if err != nil {
				return po, do, err
			}

		case reflect.String, reflect.Slice:
			n := int(*(*uint16)(unsafe.Pointer(&p[po])))
			s := []byte(p[po+2 : po+2+n])

			strptr := (*[w]byte)(unsafe.Pointer(&s))[:]
			copy(d[do:], strptr)
			copy(d[do+w:], p[po:po+2])

			if f.Kind() == reflect.Slice {
				copy(d[do+w+w:], p[po:po+2])
				do += w * 3
			} else {
				do += w * 2
			}
			po += 2 + n

		default:
			size := int(f.Type().Size())

			copy(d[do:], p[po:po+size])
			po += size
			do += size
		}
	}
	// log.Print("\n", hex.Dump(d))

	return
}
