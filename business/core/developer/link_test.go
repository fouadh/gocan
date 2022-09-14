package developer

import (
	"reflect"
	"testing"
)

func TestCombination(t *testing.T) {
	want := []Link{
		{
			Source: "A",
			Target: "B",
		},
		{
			Source: "A",
			Target: "C",
		},
		{
			Source: "B",
			Target: "C",
		},
	}

	var got []Link

	entities := make(map[string][]string)

	entities["file1"] = []string{"A", "B", "C"}

	links := buildLinks(entities)

	got = links

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Want %q, Got %q", want, got)
	}

}
