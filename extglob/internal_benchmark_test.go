package extglob

import (
	"path/filepath"
	"testing"

	"github.com/marguerite/go-stdlib/internal"
)

func BenchmarkBasenamebytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		basenamebytes("/home/marguerite")
	}
}

func BenchmarkFilepathBase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		internal.Str2bytes(filepath.Base("/home/marguerite"))
	}
}

func BenchmarkJoinBytes(b *testing.B) {
	b1 := []byte("12345")
	b2 := []byte("67890")
	for i := 0; i < b.N; i++ {
		joinbytes(b1, b2)
	}
}
