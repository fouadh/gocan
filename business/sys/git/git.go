package git

import (
	"com.fha.gocan/business/data/store/commit"
	"com.fha.gocan/business/data/store/stat"
	"com.fha.gocan/foundation"
	"com.fha.gocan/foundation/date"
	"com.fha.gocan/foundation/shell"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetCommits(path string, beforeDate string, afterDate string, beforeCommit string, afterCommit string, ctx foundation.Context) ([]commit.Commit, error) {
	ctx.Ui.Log("looking for commits between " + afterDate + " and " + beforeDate)

	args := []string{
		"log",
		"--before",
		beforeDate,
		"--date=iso",
		"--pretty=format:%H|GOCAN|%aN|GOCAN|%ad|GOCAN|%f",
	}

	if afterDate != "" {
		args = append(args, "--after", afterDate)
	}

	if afterCommit != "" {
		args = append(args, afterCommit)
		if beforeCommit != "" {
			args = append(args, ".."+beforeCommit)
		} else {
			args = append(args, "..HEAD")
		}
	} else if beforeCommit != "" {
		args = append(args, beforeCommit)
	}

	gitLogs := []gitLog{}

	cmd := exec.Command("git", args...)
	cmd.Dir = path
	ctx.Ui.Log("git command: " + cmd.String())
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "cannot run git log command")
	}
	outStr := string(out)
	if len(outStr) == 0 {
		return nil, errors.New("no output returned by the git command")
	}

	data := fmt.Sprintf("%s", outStr)
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		cols := strings.Split(line, "|GOCAN|")
		gitLogs = append(gitLogs, gitLog{
			Id:      cols[0],
			Author:  cols[1],
			Date:    cols[2],
			Message: cols[3],
		})
	}

	commits := []commit.Commit{}

	beforeTime, err := date.ParseDay(beforeDate)
	if err != nil {
		return []commit.Commit{}, errors.Wrap(err, "Unable to parse before date")
	}

	var afterTime time.Time

	if afterDate != "" {
		afterTime, err = date.ParseDay(afterDate)
		if err != nil {
			return []commit.Commit{}, errors.Wrap(err, "Unable to parse after date")
		}
	}

	ctx.Ui.Log("Building commit objects from logs...")
	for _, gc := range gitLogs {
		date, _ := time.Parse("2006-01-02 15:04:05 -0700", gc.Date)
		if afterDate != "" {
			if date.After(afterTime) && date.Before(beforeTime) {
				commits = append(commits, commit.Commit{
					Id:      gc.Id,
					Author:  gc.Author,
					Date:    date,
					Message: gc.Message,
				})
			}
		} else {
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

func GetStats(path string, before string, after string, commitsMap map[string]commit.Commit, ctx foundation.Context) ([]stat.Stat, error) {
	ctx.Ui.Log("Getting the git logs to extract stats")
	args := []string{
		"log",
		"--before",
		before,
		"--numstat",
		"--format=%H",
	}
	if after != "" {
		args = append(args, "--after", after)
	}
	cmd := exec.Command("git", args...)
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

func Checkout(commitId string, directory string) error {
	_, err := shell.ExecuteCommand("git", []string{"checkout", commitId}, directory)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("git checkout %s failed", commitId))
	}
	fmt.Println("git checkout", commitId)
	return nil
}

func GetCurrentBranch(directory string) (string, error) {
	out, err := shell.ExecuteCommand("git", []string{"rev-parse", "--abbrev-ref", "HEAD"}, directory)
	if err != nil {
		return "", errors.Wrap(err, "Cannot get git info")
	}

	initialBranch := strings.TrimRight(string(out), "\n")
	return initialBranch, nil
}

func CheckIfAllCommited(directory string) (bool, error) {
	out, err := shell.ExecuteCommand("git", []string{"status", "--porcelain"}, directory)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(out)) == "", nil
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

type gitLog struct {
	Id      string
	Author  string
	Date    string
	Message string
}
