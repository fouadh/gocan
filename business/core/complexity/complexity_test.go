package complexity

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCountLinesIndentations(t *testing.T) {
	tests := []struct {
		name string
		line string
		want int
	}{
		{name: "empty string", line: "", want: 0,},
		{name: "one line with no indentation", line: "some string", want: 0,},
		{name: "one line with one space indentation", line: "  some string", want: 1,},
		{name: "one line with two spaces indentations", line: "    some string", want: 2,},
		{name: "one line with one tab indentation", line: "\tsome string", want: 1,},

	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CountLineIndentations(test.line, 2)
			if got != test.want {
				t.Errorf("Want %d, Got %d", test.want, got)
			}
		})
	}
}

func TestComplexityAnalysis(t *testing.T) {
	file, err := ioutil.TempFile("", "gocan")
	if err != nil {
		t.Log(err)
		t.Fatalf("Cannot create temp file")
	}

	const s = `
  line 1
    line 2
    line 3
  line 4
`
	file.WriteString(s)
	file.Close()

	got, err := AnalyzeComplexity(file.Name())
	if err != nil {
		t.Log(err)
		t.Fatalf("Cannot analyze complexity")
	}

	want := Complexity{
		Indentations: 6,
	}

	if got != want {
		t.Errorf("Want %v, Got %v", want, got)
	}
}

type Complexity struct {
	Indentations int
}

func AnalyzeComplexity(filename string) (Complexity, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Complexity{}, err
	}

	contents := string(bytes)
	indentations := 0

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		fmt.Println("read line", line)
		indentations += CountLineIndentations(line, 2)
	}

	return Complexity{
		Indentations: indentations,
	}, nil
}
