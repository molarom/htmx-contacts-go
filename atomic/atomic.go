package atomic

import (
	"math"
	"sync/atomic"
	"unsafe"
)

// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://github.com/golang/go/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock() {}

type String struct {
	_ noCopy
	v unsafe.Pointer
}

type stringHeader struct {
	data unsafe.Pointer
	sz   int
}

func (s *String) Set(v string) {
	p := stringHeader{
		data: unsafe.Pointer(unsafe.StringData(v)),
		sz:   len(v),
	}
	atomic.StorePointer(&s.v, unsafe.Pointer(&p))
}

func (s *String) Value() string {
	p := (*stringHeader)(atomic.LoadPointer(&s.v))
	if p != nil {
		return unsafe.String((*byte)(p.data), p.sz)
	}
	return ""
}

type Float64 struct {
	_     noCopy
	value uint64
}

func (f *Float64) Set(value float64) {
	atomic.StoreUint64(&f.value, math.Float64bits(value))
}

func (f *Float64) Value() (value float64) {
	return math.Float64frombits(atomic.LoadUint64(&f.value))
}
