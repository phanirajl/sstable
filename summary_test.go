package sstable

import (
	"encoding/json"
	"testing"
)

func  TestSummary(t *testing.T) {
	t.Skip()
	sum, err := read_summary("./test_data/mc-99-big-Summary.db")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.Marshal(sum)
	println(string(b))
}
