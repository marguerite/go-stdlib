package extglob

import (
	"bytes"
	"testing"

	"github.com/marguerite/go-stdlib/internal"
)

func TestMakeSliceWithHyphen(t *testing.T) {
	s := makeSlice(internal.Str2bytes("-a-z-!-"))
	if cap(s) != 1 {
		t.Errorf("makeSlice test failed, expected capacity of 1, got %d", cap(s))
	}
}

func TestMakeSliceWithEquation(t *testing.T) {
	s := makeSlice(internal.Str2bytes("=c="))
	if cap(s) != 1 {
		t.Errorf("makeSlice test failed, expected capacity of 1, got %d", cap(s))
	}
}

func TestMakeSliceWithDot(t *testing.T) {
	s := makeSlice(internal.Str2bytes(".!."))
	if cap(s) != 1 {
		t.Errorf("makeSlice test failed, expected capacity of 1, got %d", cap(s))
	}
}

func TestMakeSlice(t *testing.T) {
	s := makeSlice(internal.Str2bytes("-a-z-=c=.!..-"))
	if cap(s) != 3 {
		t.Errorf("makeSlice test failed, expected capacity of 1, got %d", cap(s))
	}
}

func TestSplitextglobpatterns(t *testing.T) {
	buf := bytes.NewBufferString("g|h[|z]|")
	bufs := splitextglobpattern(buf)
	if len(bufs) == 0 || bufs[0].String() != "g" || bufs[1].String() != "h[|z]|" {
		t.Errorf("getShellPattern test failed, epxected [g, h[|z]], got %v", bufs)
	}
}

func TestJoinbytes(t *testing.T) {
	b1 := internal.Str2bytes("abc")
	b2 := internal.Str2bytes("def")
	b3 := internal.Str2bytes("ghi")
	b4 := joinbytes(b1, b2, b3)
	if internal.Bytes2str(b4) != "abcdefghi" {
		t.Errorf("joinbytes failed, expected abcdef, got %s", string(b3))
	}
}
