package main

import (
	"testing"
)

func TestKvToEntry(t *testing.T) {
	arrayKeys := []string{"arrayKey1", "arrayKey2"}
	arrayVals := []value{
		{flag: set, val: []byte{1, 2, 3, 4}},
		{flag: set, val: []byte{5, 6, 7, 8}},
	}
	arrayEntries := make([][]byte, len(arrayKeys))
	for i := range arrayKeys {
		arrayEntries[i] = kvToEntry(arrayKeys[i], arrayVals[i])
		if len(arrayEntries[i]) == 0 {
			t.Errorf("kvToEntry produced an empty entry for array entry %d", i)
		}
	}
}
