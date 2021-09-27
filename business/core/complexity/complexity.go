package complexity

import (
	"io/ioutil"
	"strings"
	"time"
)

type Core struct {

}

func NewCore() Core {
	return Core{}
}

func (c Core) CountLineIndentations(line string, size int) int {
	indentation := ""
	for i := 0; i < size; i++ {
		indentation = indentation + " "
	}
	line = strings.ReplaceAll(line, "\t", indentation)
	tline := strings.TrimLeft(line, indentation)
	return (len(line) - len(tline)) / 2
}

func (c Core) AnalyzeComplexity(filename string) (Complexity, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Complexity{}, err
	}

	contents := string(bytes)
	indentations := 0

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		indentations += c.CountLineIndentations(line, 2)
	}

	return Complexity{
		Indentations: indentations,
	}, nil
}

func (c Core) Analyze(appId string, before time.Time, after time.Time, filename string) (Complexity, error) {
	return c.AnalyzeComplexity(filename)
}

type Complexity struct {
	Indentations int
}

