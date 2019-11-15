package io

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

// BinReader is a convenient wrapper around a io.Reader and err object.
// Used to simplify error handling when reading into a struct with many fields.
type BinReader struct {
	r   io.Reader
	Err error
}

// NewBinReaderFromIO makes a BinReader from io.Reader.
func NewBinReaderFromIO(ior io.Reader) *BinReader {
	return &BinReader{r: ior}
}

// NewBinReaderFromBuf makes a BinReader from byte buffer.
func NewBinReaderFromBuf(b []byte) *BinReader {
	r := bytes.NewReader(b)
	return NewBinReaderFromIO(r)
}

// ReadLE reads from the underlying io.Reader
// into the interface v in little-endian format.
func (r *BinReader) ReadLE(v interface{}) {
	if r.Err != nil {
		return
	}
	r.Err = binary.Read(r.r, binary.LittleEndian, v)
}

// ReadArray reads array into value which must be
// a pointer to a slice.
func (r *BinReader) ReadArray(t interface{}) {
	value := reflect.ValueOf(t)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Slice {
		panic(value.Type().String() + " is not a pointer to a slice")
	}

	sliceType := value.Elem().Type()
	elemType := sliceType.Elem()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		checkHasDecodeBinary(elemType)
	} else {
		checkHasDecodeBinary(reflect.PtrTo(elemType))
	}

	if r.Err != nil {
		return
	}

	l := int(r.ReadVarUint())
	arr := reflect.MakeSlice(sliceType, l, l)

	for i := 0; i < l; i++ {
		var elem reflect.Value
		if isPtr {
			elem = reflect.New(elemType.Elem())
			arr.Index(i).Set(elem)
		} else {
			elem = arr.Index(i).Addr()
		}
		method := elem.MethodByName("DecodeBinary")
		method.Call([]reflect.Value{reflect.ValueOf(r)})
	}

	value.Elem().Set(arr)
}

func checkHasDecodeBinary(v reflect.Type) {
	method, ok := v.MethodByName("DecodeBinary")
	if !ok || !isDecodeBinaryMethod(method) {
		panic(v.String() + " does not have DecodeBinary(*io.BinReader)")
	}
}

func isDecodeBinaryMethod(method reflect.Method) bool {
	t := method.Type
	return t != nil &&
		t.NumIn() == 2 && t.In(1) == reflect.TypeOf((*BinReader)(nil)) &&
		t.NumOut() == 0
}

// ReadBE reads from the underlying io.Reader
// into the interface v in big-endian format.
func (r *BinReader) ReadBE(v interface{}) {
	if r.Err != nil {
		return
	}
	r.Err = binary.Read(r.r, binary.BigEndian, v)
}

// ReadVarUint reads a variable-length-encoded integer from the
// underlying reader.
func (r *BinReader) ReadVarUint() uint64 {
	if r.Err != nil {
		return 0
	}

	var b uint8
	r.Err = binary.Read(r.r, binary.LittleEndian, &b)

	if b == 0xfd {
		var v uint16
		r.Err = binary.Read(r.r, binary.LittleEndian, &v)
		return uint64(v)
	}
	if b == 0xfe {
		var v uint32
		r.Err = binary.Read(r.r, binary.LittleEndian, &v)
		return uint64(v)
	}
	if b == 0xff {
		var v uint64
		r.Err = binary.Read(r.r, binary.LittleEndian, &v)
		return v
	}

	return uint64(b)
}

// ReadBytes reads the next set of bytes from the underlying reader.
// ReadVarUInt() is used to determine how large that slice is
func (r *BinReader) ReadBytes() []byte {
	n := r.ReadVarUint()
	b := make([]byte, n)
	r.ReadLE(b)
	return b
}

// ReadString calls ReadBytes and casts the results as a string.
func (r *BinReader) ReadString() string {
	b := r.ReadBytes()
	return string(b)
}