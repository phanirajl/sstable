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
	i := 10
	for {
		i--
		if i == 0 {
			break
		}
		var entry IndexEntry
		entry.KeyLength = binary.BigEndian.Uint16(data[0:2])
		entry.Key = data[2 : 2+entry.KeyLength]
		data = data[2+entry.KeyLength:]
		position, n := binary.Uvarint(data)
		entry.Position = uint64(position)
		for u := 0; u < n; u++ {
			println("SSS", data[u])
		}
		data = data[n:]

		prom_index_len, n2 := binary.Uvarint(data)

		data = data[n2:]

		entry.PromotedIndexLength = uint32(prom_index_len)
		println("KEY", entry.KeyLength, position, "PROM INDEX LEN", prom_index_len, "N", n, "N2", n2)

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
