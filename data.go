package sstable

type DeletionTime struct {
	local_deletion_time  uint32
	marked_for_delete_at uint64
}

type ClusteringBlock struct {
	ClusteringBlockHeader int
	clustering_cells      []SimpleCell
}

type SimpleCell struct {
}
