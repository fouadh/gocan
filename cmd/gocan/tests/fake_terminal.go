package tests

import (
	terminal2 "com.fha.gocan/internal/platform/terminal"
)

type FakeUI struct {
	FakeTable   *fakeTable
	Said        []string
	Errors      []string
	OksCount    int
	FailedCount int
}

type fakeTable struct {
	Rows               [][]string
	Headers            []string
	HasPrintBeenCalled bool
}

func (f *fakeTable) Add(row ...string) {
	f.Rows = append(f.Rows, row)
}

func (f *fakeTable) Print() {
	f.HasPrintBeenCalled = true
}

func (f *FakeUI) Say(message string) {
	f.Said = append(f.Said, message)
}

func (f *FakeUI) Ok() {
	f.OksCount++
}

func (f *FakeUI) Failed(message string) {
	f.FailedCount++
	f.Errors = append(f.Errors, message)
}

func (f *FakeUI) Table(headers []string) terminal2.UITable {
	t := &fakeTable{
		Headers: headers,
	}
	f.FakeTable = t
	return t
}
