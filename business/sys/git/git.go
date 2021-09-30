package git

import (
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/data/store/stat"
	"com.fha.gocan/foundation/date"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetCommits(path string, before time.Time, after time.Time) ([]commit.Commit, error) {
	fmt.Println("looking for commits between " + date.FormatDay(after) + " and " + date.FormatDay(before))
	cmd := exec.Command("git", "log", "--after", date.FormatDay(after), "--before", date.FormatDay(before),
		"--date=iso", "--pretty=format:{%n  \"Id\": \"%H\",%n  \"Author\": \"%aN\",%n  \"Date\": \"%ad\",%n  \"Message\": \"%f\"%n},")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "cannot run git log command")
	}

	outStr := string(out)
	if len(outStr) == 0 {
		return nil, errors.New("no output returned by the git command")
	}

	outStr = outStr[:len(outStr)-1]
	data := fmt.Sprintf("[%s]", outStr)
	data = strings.ReplaceAll(data, "\\", "\\\\")
	gitCommits := []gitCommit{}
	if err := json.Unmarshal([]byte(data), &gitCommits); err != nil {
		return nil, err
	}

	if err != nil {
		return []commit.Commit{}, fmt.Errorf("Cannot get commits: %s", err)
	}

	commits := []commit.Commit{}
	for _, gc := range gitCommits {
		date, _ := time.Parse("2006-01-02 15:04:05 -0700", gc.Date)
		if date.After(after) && date.Before(before) {
			commits = append(commits, commit.Commit{
				Id:      gc.Id,
				Author:  gc.Author,
				Date:    date,
				Message: gc.Message,
			})
		}
	}

	return commits, nil
}

func GetStats(path string, before time.Time, after time.Time, commitsMap map[string]commit.Commit) ([]stat.Stat, error) {
	cmd := exec.Command("git", "log", "--after", date.FormatDay(after), "--before", date.FormatDay(before), "--numstat", "--format=%H")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get source code stats")
	}
	outStr := string(out)
	rows := strings.Split(outStr, "\n")
	stats := []stat.Stat{}
	var currentCommit string
	for _, row := range rows {
		if len(strings.Split(row, "\t")) > 1 {
			if _, ok := commitsMap[currentCommit]; ok {
				stats = append(stats, buildStat(currentCommit, row))
			}
		} else if row != "" {
			currentCommit = row
		}
	}
	return stats, nil
}

func buildStat(commitId string, line string) stat.Stat {
	cols := strings.Split(line, "\t")
	insertions, _ := strconv.Atoi(cols[0])
	deletions, _ := strconv.Atoi(cols[1])
	return stat.Stat{
		CommitId:   commitId,
		File:       cols[2],
		Insertions: insertions,
		Deletions:  deletions,
	}
}

type gitCommit struct {
	Id      string
	Author  string
	Date    string
	Message string
}

