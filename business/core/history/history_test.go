package history

import (
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/stat"
	"reflect"
	"testing"
)

func TestCouplingForTwoFiles(t *testing.T) {
	stats := []stat.Stat{
		{CommitId: "123", File: "file1"},
		{CommitId: "123", File: "file2"},
		{CommitId: "456", File: "file1"},
	}

	want := []coupling.Coupling{
		{
			Entity:           "file1",
			Coupled:          "file2",
			Degree:           0.6666666666666666,
			AverageRevisions: 1.5,
		},
	}
	got := BuildCoupling(stats)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}


