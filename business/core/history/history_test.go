package history

import (
	"com.fha.gocan/business/data/store/coupling"
	"com.fha.gocan/business/data/store/stat"
	"reflect"
	"testing"
)

func TestCouplingForOnePair(t *testing.T) {
	stats := []stat.Stat{
		{CommitId: "123", File: "file1"},
		{CommitId: "123", File: "file2"},
	}

	want := []coupling.Coupling{
		{
			Entity:           "file1",
			Coupled:          "file2",
			Degree:           1,
			AverageRevisions: 1,
		},
	}
	got := BuildCoupling(stats)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %v, got %v", want, got)
	}

}
