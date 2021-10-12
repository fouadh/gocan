package terminal

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"io"
	"regexp"
	"strings"
)

const defaultFgColor = 38

func ColorizeBold(message string, textColor color.Attribute) string {
	return colorize(message, textColor, 1)
}

func HeaderColor(header string) string {
	return ColorizeBold(header, defaultFgColor)
}

func TableContentHeaderColor(message string) string {
	return ColorizeBold(message, color.FgCyan)
}

var decolorizerRegex = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`)

func Decolorize(message string) string {
	return string(decolorizerRegex.ReplaceAll([]byte(message), []byte("")))
}

type Transformer func(s string) string

func nop(s string) string {
	return s
}

func colorize(message string, textColor color.Attribute, bold int) string {
	colorPrinter := color.New(textColor)
	if bold == 1 {
		colorPrinter = colorPrinter.Add(color.Bold)
	}
	f := colorPrinter.SprintFunc()
	return f(message)
}

type UITable interface {
	Add(row ...string)
	Print()
}

type Table struct {
	ui          UI
	headers     []string
	rows        [][]string
	columnWidth []int
	colSpacing  string
	transformer []Transformer
	csv         bool
}

type rowTransformer interface {
	Transform(column int, s string) string
}

type transformHeader struct{}

var transHeader = &transformHeader{}

func (th *transformHeader) Transform(column int, s string) string {
	return HeaderColor(s)
}

func (t *Table) Transform(column int, s string) string {
	return t.transformer[column](s)
}

func (t *terminal) Table(headers []string, csv bool) UITable {
	pt := &Table{
		ui:          t,
		headers:     headers,
		columnWidth: make([]int, len(headers)),
		colSpacing:  "    ",
		transformer: make([]Transformer, len(headers)),
		csv: csv,
	}

	for i := range pt.transformer {
		pt.transformer[i] = nop
	}
	if len(headers) > 0 {
		pt.transformer[0] = TableContentHeaderColor
	}
	return pt
}

func (t *Table) Add(row ...string) {
	t.rows = append(t.rows, row)
}

func (t *Table) CalculateMaxWidth(transformer rowTransformer, cols []string) {
	for columnIndex, col := range cols {
		width := visibleSize(Decolorize(transformer.Transform(columnIndex, col)))
		if t.columnWidth[columnIndex] < width {
			t.columnWidth[columnIndex] = width
		}
	}
}

func visibleSize(s string) int {
	r := strings.NewReader(s)
	var size int
	for range s {
		_, runeSize, _ := r.ReadRune()
		if runeSize == 3 {
			size += 2 // Kanji and Katakana characters appear as double-width
		} else {
			size++
		}
	}
	return size
}

func (t *Table) Print() {
	if t.csv {
		t.ui.Say(strings.Join(t.headers, ","))
		for _, row := range t.rows {
			t.ui.Say(strings.Join(row, ","))
		}
	} else {
		rowIndex := 0
		t.CalculateMaxWidth(transHeader, t.headers)
		rowIndex++

		for _, row := range t.rows {
			t.CalculateMaxWidth(t, row)
			rowIndex++
		}

		rowIndex = 0
		t.printRow(transHeader, rowIndex, t.headers)
		rowIndex++
		for _, row := range t.rows {
			t.printRow(t, rowIndex, row)
		}
	}
}

func (t *Table) printRow(transformer rowTransformer, rowIndex int, row []string) {
	last := len(t.headers) - 1
	line := &bytes.Buffer{}
	for columnIndex, col := range row {
		t.printCellValue(line, transformer, columnIndex, last, col)
	}
	t.ui.Say(strings.TrimSpace(string(line.Bytes())))
	rowIndex++
}

func (t *Table) printCellValue(line *bytes.Buffer, transformer rowTransformer, col, last int, value string) {
	value = transformer.Transform(col, value)
	fmt.Fprintf(line, value)
	if col < last {
		width := visibleSize(Decolorize(value))
		padlen := t.columnWidth[col] - width
		padding := strings.Repeat(" ", padlen)
		fmt.Fprintf(line, padding)
		fmt.Fprintf(line, t.colSpacing)
	}
}

type UI interface {
	Say(message string)
	Ok()
	Failed(message string)
	Table(headers []string, csv bool) UITable
}

func NewUI(stdout io.Writer, stderr io.Writer) UI {
	return &terminal{
		stderr: stderr,
		stdout: stdout,
	}
}

type terminal struct {
	stderr io.Writer
	stdout io.Writer
}

func (t *terminal) Failed(message string) {
	t.SayError(ColorizeBold("\nFAILED", color.FgRed))
	t.SayError(message)
}

func (t *terminal) SayError(message string) {
	fmt.Fprintln(t.stderr, message)
}

func (t *terminal) Ok() {
	t.Say(ColorizeBold("OK\n", color.FgGreen))
}

func (t *terminal) Say(message string) {
	fmt.Fprintln(t.stdout, message)
}
