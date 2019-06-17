package sstable

import (
	"testing"
)

func TestIndex(t *testing.T) {
	_, err := read_index("../sample/mc-4360-big-Index.db")
	if err != nil {
		t.Fatal(err)
	}

}

func BenchmarkDecodeVarint(b *testing.B) {
	data := []byte{0xda, 0xf0, 0x41}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		decodeVarint(data)
	}
}
