package sstable

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

type Summary struct {
	Header         SummaryHeader
	SummaryEntries SummaryEntriesBlock
	First, Last    SerializedKey
}

type SummaryHeader struct {
	MinIndexInterval   uint32
	EntriesCount       uint32
	SummaryEntriesSize uint64
	SamplingLevel      uint32
	SizeAtFullSampling uint32
}

type SummaryEntriesBlock struct {
	Offsets []uint32
	Entries []SummaryEntry
}

type SummaryEntry struct {
	Key      []byte
	Position uint64
}

type SerializedKey struct {
	Size uint32
	Key  []byte
}

func read_summary(filepath string) (Summary, error) {
	var sum Summary

	data, err := ioutil.ReadFile(filepath)
	// f, err := os.Open(filepath)
	if err != nil {
		return sum, err
	}

	// fill header
	sum.Header.MinIndexInterval = binary.BigEndian.Uint32(data[0:4])
	sum.Header.EntriesCount = binary.BigEndian.Uint32(data[4:8])
	sum.Header.SummaryEntriesSize = binary.BigEndian.Uint64(data[8:16])
	sum.Header.SamplingLevel = binary.BigEndian.Uint32(data[16:20])
	sum.Header.SizeAtFullSampling = binary.BigEndian.Uint32(data[20:24])
	data = data[24:]

	var entries_count = int(sum.Header.EntriesCount)
	offsets := make([]uint32, entries_count, entries_count)
	for i, off := 0, 0; i < entries_count; i++ {
		offsets[i] = binary.LittleEndian.Uint32(data[off : off+4])
		off += 4
	}
	sum.SummaryEntries.Offsets = offsets

	// reading entries
	var SUM_ENT_SIZE = uint32(sum.Header.SummaryEntriesSize)
	entries := make([]SummaryEntry, entries_count, entries_count)
	for i := 0; i < entries_count; i++ {
		var nextoffset = SUM_ENT_SIZE
		if i < entries_count-1 {
			nextoffset = offsets[i+1]
		}
		entries[i].Key = data[offsets[i] : nextoffset-8]
		entries[i].Position = binary.BigEndian.Uint64(data[nextoffset-8 : nextoffset])
	}
	sum.SummaryEntries.Entries = entries
	data = data[sum.Header.SummaryEntriesSize:]

	// read first and last key
	sum.First.Size = binary.BigEndian.Uint32(data[0:4])
	sum.First.Key = data[4:sum.First.Size+4]
	sum.Last.Size = binary.BigEndian.Uint32(data[sum.First.Size+4:sum.First.Size+8])
	sum.Last.Key = data[sum.First.Size+8:sum.First.Size+8+sum.Last.Size]
	data = data[sum.First.Size+8+sum.Last.Size:]

	if len(data) != 0 {
		return sum, fmt.Errorf("corrupted summary, trailing bytes")
	}
	return sum, nil
}
