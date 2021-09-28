package complexity

import (
	"com.fha.gocan/foundation/date"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
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

func (c Core) AnalyzeComplexity(filename string, date time.Time) (Complexity, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Complexity{}, err
	}

	contents := string(bytes)
	indentations := []int{}
	indentationsCounter := 0
	linesCounter := 0
	max := 0

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		if line != "" {
			linesCounter++
			lineIndentations := c.CountLineIndentations(line, 2)
			indentations = append(indentations, lineIndentations)
			indentationsCounter += lineIndentations
			if max < lineIndentations {
				max = lineIndentations
			}
		}
	}

	mean := float64(indentationsCounter) / float64(linesCounter)
	stdev := 0.
	for _, i := range indentations {
		stdev += (float64(i)-mean) * (float64(i)-mean)
	}
	stdev = math.Sqrt(stdev / float64(len(indentations)))

	return Complexity{
		Indentations: indentationsCounter,
		Lines:        linesCounter,
		Mean:         mean,
		Max:          max,
		Stdev:        stdev,
		Date:         date,
	}, nil
}

func (c Core) CreateComplexityAnalysis(analysisName string, appId string, before time.Time, after time.Time, filename string, directory string) ([]Complexity, error) {
	cmd := exec.Command("git", "log", "--oneline", "--pretty=format:%h;%ad", "--after", date.FormatDay(after), "--before", date.FormatDay(before), "--date=iso")
	cmd.Dir = directory
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return []Complexity{}, errors.Wrap(err, "Unable to get the revisions")
	}

	outStr := string(out)
	if len(outStr) == 0 {
		return []Complexity{}, errors.New("No output returned by the git command")
	}

	complexities := []Complexity{}
	lines := strings.Split(outStr, "\n")
	for _, line := range lines {
		cols := strings.Split(line, ";")
		rev := cols[0]
		revDate, err := time.Parse("2006-01-02 15:04:05 -0700", cols[1])
		if err != nil {
			return []Complexity{}, errors.Wrap(err, "Unable to parse date")
		}

		cmd := exec.Command("git", "checkout", rev)
		cmd.Dir = directory
		cmd.Stderr = os.Stderr
		out, err = cmd.Output()
		if err != nil {
			return []Complexity{}, errors.Wrap(err, "Fail to checkout revision "+rev)
		}

		c, err := c.AnalyzeComplexity(directory+filename, revDate)
		if err != nil {
			// in case of error, we consider that the file's complexity is 0 for every field
			fmt.Println("File cannot be analyzed for revision " + rev)
			return []Complexity{}, nil
		}

		complexities = append(complexities, c)
	}

	return complexities, nil
}

type Complexity struct {
	Lines        int
	Indentations int
	Mean         float64
	Max          int
	Stdev        float64
	Date         time.Time
}
