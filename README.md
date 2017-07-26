# stob - convert struct to bytes and back 

stob is package for convert structs to bytes or fill structs from bytes, as it does on `C`, i.e convert raw bytes to struct, and restore struct from raw bytes. 

Unfortunately `go` does not allow do it also as simple as it can be done on `C`. There are 3 solutions:

1. use `unsafe` package - this is very fast, but not flexibly.

2. generatable code - it`s also fast, but it is not convenient.

3. reflection - it`s relatively slow...

stob use reflection for read and write struct to bytes.


# Install

```
go get github.com/sg3des/stob
```

# Example

```go
type YourStruct struct {
	Str      string 
	Int      int `bo:"le"`
	Byte     byte
	Bytes    []byte `num:"8"`
	Bytes4   [4]byte
	Bool     bool
	Float    float32 `bo:"be"`
}

a := YourStruct{
	Str:      "string",
	Int:      999,
	Byte:     255,
	Bytes:    []byte{1, 2, 3, 4, 5, 6, 7, 8},
	Bytes4:   [4]byte{10, 11, 12, 13},
	Bool:     rand.Intn(2) == 1,
	Float:    rand.Float32(),
}

stob.Write(w, a)
```

it will be write how:
```
73 74 72 69 6e 67 00 e7  03 00 00 00 00 00 00 ff  |string..........|
01 02 03 04 05 06 07 08  0a 0b 0c 0d 01 3f 70 c5  |.............?p.|
34                                                |4|
```

to restore struct from bytes:

```go
stob.Read(r, &a)
```

## Tags

stob knows 3 tags:

 * `bo:"le"` or `bo:"be"` - it`s byte order little or big endian
 * `num:"8"` - count of elements in slice
 * `size:"4"` - size of element, example size of string, but it also allows read\write big integers to small number of bytes.

**WARNING:** if `[]byte` slice does not have *num* tag, then all next bytes will be writed to this field!


# Benchmark

	BenchmarkRead-8    2000000    762 ns/op    10 B/op    10 allocs/op
	BenchmarkWrite-8   2000000    750 ns/op    24 B/op    6 allocs/op

# TODO

* types
* tests
* reader and writer for custom types
