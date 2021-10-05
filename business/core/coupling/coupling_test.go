package coupling

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
	got := CalculateCouplings(stats, 0, 0)

	if !isEqual(want, got) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func TestCouplingForManyFiles(t *testing.T) {
	stats := []stat.Stat{
		{CommitId: "123", File: "file1"},
		{CommitId: "123", File: "file2"},

		{CommitId: "456", File: "file1"},
		{CommitId: "456", File: "file3"},
		{CommitId: "456", File: "file4"},

		{CommitId: "789", File: "file4"},
		{CommitId: "789", File: "file2"},
		{CommitId: "789", File: "file1"},
	}

	want := []coupling.Coupling{
		{
			Entity:           "file1",
			Coupled:          "file2",
			Degree:           .8,
			AverageRevisions: 2.5,
		},

		{
			Entity:           "file1",
			Coupled:          "file3",
			Degree:           0.5,
			AverageRevisions: 2,
		},

		{
			Entity:           "file1",
			Coupled:          "file4",
			Degree:           .8,
			AverageRevisions: 2.5,
		},
		{
			Entity:           "file3",
			Coupled:          "file4",
			Degree:           0.6666666666666666,
			AverageRevisions: 1.5,
		},
		{
			Entity:           "file4",
			Coupled:          "file2",
			Degree:           .5,
			AverageRevisions: 2,
		},
	}
	got := CalculateCouplings(stats, 0, 0)

	if !isEqual(want, got) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func isEqual(aa, bb []coupling.Coupling) bool {
	eqCtr := 0
	for _, a := range aa {
		for _, b := range bb {
			if reflect.DeepEqual(a, b) {
				eqCtr++
			}
		}
	}
	if eqCtr != len(bb) || len(aa) != len(bb) {
		return false
	}
	return true
}