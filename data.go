package sstable

type DeletionTime struct {
	LocalDeletionTime uint32
	MarkedForDeleteAt uint64
}

type ClusteringBlock struct {
	ClusteringBlockHeader int
	ClusteringCells       []SimpleCell
}

type SimpleCell struct {
	Cell
	Flags                   byte
	DeltaTimestamp          uint64
	DeltaLocalDeletion_time uint64
	DeltaTtl                uint64
	Path                    CellPath // only in cells nested into complex_cells
	Value                   CellValue
}

type Unfiltered struct{}

type CellPath struct {
	Length uint64
	Value  []byte
}

type CellValue struct {
	Length uint64
	Value  []byte
}

type Partition struct {
	Header      PartitionHeader
	StaticRow   Row // Has IS_STATIC flag set
	Unfiltereds []Unfiltered
}

type Row struct {
	Flags         byte
	ExtendedFlags byte // optional

	// only present for non-static rows
	ClusteringBlocks   []ClusteringBlock
	RowBodySize        uint64
	PrevUnfilteredSize uint64 // for backward traversing
	LivenessInfo       LivenessInfo
	DeletionTime       DeltaDeletionTime
	MissingColumns     []uint64
	Cells              []Cell
}

type Cell struct{}

type ComplexCell struct {
	Cell
	ComplexDeletionTime DeltaDeletionTime
	ItemsCount          uint64
	Items               []SimpleCell
}

type DeltaDeletionTime struct {
	DeltaMarkedForDeleteAt uint64
	DeltaLocalDeletionTime uint64
}

type LivenessInfo struct {
	DeltaTimestamp         uint64
	DeltaTtl               uint64
	DeltaLocalDeletionTime uint64
}

type PartitionHeader struct {
	KeyLength    uint16
	Key          []byte
	DeletionTime DeletionTime
}
