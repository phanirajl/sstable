package sstable

import (
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
	// fill header
	i := 100
	for {
		i--
		if i == 0 {
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
		println("KEY", string(entry.Key), position, "PROM INDEX LEN", prom_index_len, "N", n, "N2", n2)

		entry.PromotedIndex = data[:entry.PromotedIndexLength]

		data = data[entry.PromotedIndexLength:]
		index = append(index, entry)
		if len(data) == 0 {
			println("KEY", string(entry.Key))
			break
		}
	}
	return nil, nil
}

// decodeVarint decodes an uint64 from buf and returns that value and the
// number of bytes read (> 0)
// the internal representation of varints is explained in:
// https://haaawk.github.io/2018/02/26/sstables-variant-integers.html
// Note: the first byte of buf will be modified for performance reason
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
	if first&128 == 128 {
		buf[0] &= 191
		return uint64(binary.LittleEndian.Uint16(buf[0:2])), 2
	}

	// 3 bytes
	if first&96 == 96 {
		buf[0] &= 223
		return uint64(binary.LittleEndian.Uint32(buf[0:3])), 3
	}

	// 4 bytes
	if first&224 == 224 {
		buf[0] &= 239
		return uint64(binary.LittleEndian.Uint32(buf[0:4])), 4
	}

	// 5 bytes
	if first&240 == 240 {
		buf[0] &= 247
		return binary.LittleEndian.Uint64(buf[0:5]), 5
	}

	// 6 bytes
	if first&248 == 248 {
		buf[0] &= 251
		return binary.LittleEndian.Uint64(buf[0:6]), 6
	}

	// 7 bytes
	if first&252 == 252 {
		buf[0] &= 253
		return binary.LittleEndian.Uint64(buf[0:7]), 7
	}

	// 8 bytes
	if first&254 == 254 {
		buf[0] &= 254
		return binary.LittleEndian.Uint64(buf[0:8]), 8
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}
