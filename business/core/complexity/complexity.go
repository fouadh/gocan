package complexity

import (
	"com.fha.gocan/business/data/store/complexity"
	"com.fha.gocan/business/sys/git"
	"com.fha.gocan/foundation/date"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Core struct {
	complexity complexity.Store
}

func NewCore(connection *sqlx.DB) Core {
	return Core{
		complexity: complexity.NewStore(connection),
	}
}

func (c Core) CountLineIndentations(line string, size int) int {
	indentation := ""
	for i := 0; i < size; i++ {
		indentation = indentation + " "
	}
	line = strings.ReplaceAll(line, "\t", indentation)
	tline := strings.TrimLeft(line, indentation)
	return (len(line) - len(tline)) / size
}

func (c Core) AnalyzeComplexity(complexityId string, filename string, date time.Time, spaces int) (complexity.ComplexityEntry, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return complexity.ComplexityEntry{}, err
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
			lineIndentations := c.CountLineIndentations(line, spaces)
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
		stdev += (float64(i) - mean) * (float64(i) - mean)
	}
	stdev = math.Sqrt(stdev / float64(len(indentations)))

	return complexity.ComplexityEntry{
		ComplexityId: complexityId,
		Indentations: indentationsCounter,
		Lines:        linesCounter,
		Mean:         mean,
		Max:          max,
		Stdev:        stdev,
		Date:         date,
	}, nil
}

func (c Core) CreateComplexityAnalysis(analysisName string, appId string, before time.Time, after time.Time, filename string, directory string, spaces int) (complexity.Complexity, error) {
	cmd := exec.Command("git", "log", "--oneline", "--pretty=format:%h;%ad", "--after", date.FormatDay(after), "--before", date.FormatDay(before), "--date=iso")
	cmd.Dir = directory
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return complexity.Complexity{}, errors.Wrap(err, "Unable to get the revisions")
	}

	outStr := string(out)
	if len(outStr) == 0 {
		return complexity.Complexity{}, errors.New("No output returned by the git command")
	}

	complexities := []complexity.ComplexityEntry{}
	complexityId := uuid.New()
	lines := strings.Split(outStr, "\n")

	initialBranch, err := git.GetCurrentBranch(directory)
	if err != nil {
		return complexity.Complexity{}, errors.Wrap(err, "Unable to get current branch info")
	}

	if !strings.HasSuffix(directory, "/") && !strings.HasPrefix(filename, "/") {
		directory += "/"
	}
	filePath := directory + filename

	for _, line := range lines {
		cols := strings.Split(line, ";")
		rev := cols[0]
		revDate, err := time.Parse("2006-01-02 15:04:05 -0700", cols[1])
		if err != nil {
			return complexity.Complexity{}, errors.Wrap(err, "Unable to parse date")
		}

		if revDate.After(after) && revDate.Before(before) {
			cmd := exec.Command("git", "checkout", rev)
			cmd.Dir = directory
			cmd.Stderr = os.Stderr
			out, err = cmd.Output()
			if err != nil {
				return complexity.Complexity{}, errors.Wrap(err, "Fail to checkout revision "+rev)
			}

			c, err := c.AnalyzeComplexity(complexityId, filePath, revDate, spaces)
			if err != nil {
				fmt.Println("WARNING: File cannot be analyzed for revision " + rev)
			} else {
				complexities = append(complexities, c)
			}
		}
	}

	err = git.Checkout(initialBranch, directory)
	if err != nil {
		return complexity.Complexity{}, errors.Wrap(err, "Unable to reinitialize initial branch")
	}

	result := complexity.Complexity{
		Id:      complexityId,
		Name:    analysisName,
		Entity:  filename,
		AppId:   appId,
		Entries: complexities,
	}

	return c.complexity.Create(result)
}

func (c Core) QueryAnalyses(appId string) ([]complexity.ComplexityAnalysisSummary, error) {
	return c.complexity.QueryAnalyses(appId)
}

func (c Core) QueryAnalysisEntriesById(complexityId string) ([]complexity.ComplexityEntry, error) {
	return c.complexity.QueryAnalysisEntriesById(complexityId)
}

func (c Core) DeleteAnalysisByName(appId string, analysisName string) error {
	return c.complexity.DeleteAnalysisByName(appId, analysisName)
}
