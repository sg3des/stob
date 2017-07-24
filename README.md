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
	Int      int `stob:",le"`
	Byte     byte
	Bytes    []byte `stob:"8"`
	Bytes4   [4]byte
	Bool     bool
	Float    float32 `stob:",be"`
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
