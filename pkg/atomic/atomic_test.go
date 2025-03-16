package atomic_test

import (
	"testing"

	"gitlab.com/romalor/htmx-contacts/pkg/atomic"
)

func Test_String(t *testing.T) {
	var s atomic.String
	if s.Value() != "" {
		t.Fatal("unknown value returned from String")
	}

	in := "values"
	s.Set(in)
	if v := s.Value(); v != in {
		if v != "" {
			t.Fatalf("expected: [%v]; got: [%v]", in, v)
		}
		t.Fatal("string value was not set")
	}

	s.Set("")
	if v := s.Value(); v != "" {
		t.Fatalf("value was not changed; got [%v]", v)
	}
}

func Test_Float64(t *testing.T) {
	var f atomic.Float64
	if f.Value() != 0 {
		t.Fatal("unknown value returned from Float64")
	}

	in := 3.14
	f.Set(in)
	if v := f.Value(); v != in {
		if v != 0 {
			t.Fatalf("expected: [%v]; got: [%v]", in, v)
		}
		t.Fatal("float value was not set")
	}

	f.Set(0.0)
	if v := f.Value(); v != 0.0 {
		t.Fatalf("value was not changed; got [%v]", v)
	}
}
