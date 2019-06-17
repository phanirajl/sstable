package sstable

import (
	"fmt"
	"encoding/binary"
	"io/ioutil"
)

type IndexEntry struct {
	KeyLength uint16
	Key       []byte

	Position            uint64
	PromotedIndexLength uint32
	PromotedIndex       []byte
}

type PromotedIndex struct {
	PartitionHeaderLength    uint64
	DeletionTime             DeletionTime
	PromotedIndexBlocksCount uint32
	Block                    []PromotedIndexBlock
	Offsets                  []uint32
}

type PromotedIndexBlock struct {
	FirstName            ClusteringPrefix
	LastName             ClusteringPrefix
	Offset               int
	DeltaWidth           int
	EndOpenMarkerPresent byte
	EndOpenMarker        DeletionTime
}

type ClusteringPrefix struct {
	kind             byte
	size             uint16 // optional
	ClusteringBlocks []ClusteringBlock
}

func read_index(filepath string) ([]IndexEntry, error) {
	data, err := ioutil.ReadFile(filepath)
	// f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	index := make([]IndexEntry, 0)
	for {
		if len(data) == 0 {
			break
		}
		var entry IndexEntry
		entry.KeyLength = binary.BigEndian.Uint16(data[0:2])
		entry.Key = data[2 : 2+entry.KeyLength]
		data = data[2+entry.KeyLength:]
		position, n := decodeVarint(data)
		entry.Position = uint64(position)
		data = data[n:]
		prom_index_len, n2 := decodeVarint(data)
		data = data[n2:]

		entry.PromotedIndexLength = uint32(prom_index_len)
		// fmt.Println("KEY", string(entry.Key), position, "PROM INDEX LEN", prom_index_len)

		entry.PromotedIndex = data[:prom_index_len]
		data = data[prom_index_len:]
		index = append(index, entry)
	}
	fmt.Println(len(index))
	return nil, nil
}

// decodeVarint decodes an uint64 from buf and returns that value and the
// number of bytes read (> 0)
// the internal representation of varints is explained in:
// https://haaawk.github.io/2018/02/26/sstables-variant-integers.html
func decodeVarint(buf []byte) (x uint64, n int) {
	first := buf[0]
	// reading 0 byte
	if first&128 == 0 {
		return uint64(first), 1
	}

	// 9 bytes
	if first == 0xff {
		return binary.LittleEndian.Uint64(buf[1:9]), 9
	}

	// 2 bytes
	if first&192 == 128 { // first & 1100 000 == 1000 0000
		var b [2]byte
		b[0] = buf[0] & 63 // 0011 1111
		b[1] = buf[1]
		return uint64(binary.LittleEndian.Uint16(b[:])), 2
	}

	// 3 bytes
	if first&224 == 192 { //  first & 11100000 == 1100 0000
		// fmt.Printf("%x %x %x", buf[0], buf[1], buf[2])
		var b [4]byte
		b[0] = 0
		b[1] = buf[0] & 31 // 0001 1111
		b[2] = buf[1]
		b[3] = buf[2]
		return uint64(binary.LittleEndian.Uint32(b[:])), 3
	}

	// 4 bytes
	if first&240 == 224 { // 240: 11110000
		var b [4]byte
		b[0] = buf[0] & 15 // 0000 1111
		b[1] = buf[1]
		b[2] = buf[2]
		b[3] = buf[3]
		return uint64(binary.LittleEndian.Uint32(b[:])), 4
	}

	// 5 bytes
	if first&248 == 240 { // 248: 11111000
		var b [8]byte
		b[3] = buf[0] & 7 // 0000 0111
		b[4] = buf[1]
		b[5] = buf[2]
		b[6] = buf[3]
		b[7] = buf[4]
		return binary.LittleEndian.Uint64(b[:]), 5
	}

	// 6 bytes
	if first&252 == 248 { // 252: 1111 1100
		var b [8]byte
		b[2] = buf[0] & 3 // 0000 0011
		b[3] = buf[1]
		b[4] = buf[2]
		b[5] = buf[3]
		b[6] = buf[4]
		b[7] = buf[5]
		return binary.LittleEndian.Uint64(b[:]), 6
	}

	// 7 bytes
	if first&254 == 252 { // 254: 1111 1110
		var b [8]byte
		b[1] = buf[0] & 1 // 0000 0001
		b[2] = buf[1]
		b[3] = buf[2]
		b[4] = buf[3]
		b[5] = buf[4]
		b[6] = buf[5]
		b[7] = buf[6]
		return binary.LittleEndian.Uint64(buf[0:7]), 7
	}

	// 8 bytes
	if first&255 == 254 {
		return binary.LittleEndian.Uint64(buf[0:8]), 8
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}
