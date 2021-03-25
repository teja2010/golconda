package golconda

import (
	d "github.com/teja2010/golconda/src/debug"
	"regexp"
	"strings"
)

const (
	_NEWLINE = "\n"
)

func Regex2Func(rule string) func(string) bool {
	exp, err := regexp.Compile(rule)
	if d.DebugCheck(err != nil) {
		d.Bug("Invalid regexp", rule)
	}

	f := func(l string) bool {
		return exp.MatchString(l)
	}

	return f
}

func Filter(lines []string, f func(string) bool) []string {

	matched_len := 0
	matched_lines := make([]string, len(lines))

	for _, line := range lines {
		if f(line) {
			matched_lines[matched_len] = line
			matched_len++
		}
	}

	return matched_lines[:matched_len]
}

func TakeWhile(lines []string, f func(string) bool) []string {
	matched_lines := make([]string, len(lines))

	for i, line := range lines {
		if f(line) {
			matched_lines[i] = line
		} else {
			return matched_lines[:i]
		}
	}

	return matched_lines
}

// matches
func FindLine(lines []string, f func(string) bool) string {
	for _, line := range lines {
		if f(line) {
			return line
		}
	}

	d.Bug("Unable to find a line matching", f)
	return "BUG!! BUG!!"
}

// clean-up after using generics
func FmapSS(lines []string, f func(string) string) []string {
	res := make([]string, len(lines))

	for i, l := range lines {
		res[i] = f(l)
	}

	return res
}

func FmapSI(lines []string, f func(string) int) []int {
	res := make([]int, len(lines))

	for i, l := range lines {
		res[i] = f(l)
	}

	return res
}

func FmapSI64(lines []string, f func(string) int64) []int64 {
	res := make([]int64, len(lines))

	for i, l := range lines {
		res[i] = f(l)
	}

	return res
}

func FmapSCpu_stat(lines []string, f func(string) cpu_stat_data) []cpu_stat_data {
	res := make([]cpu_stat_data, len(lines))

	for i, l := range lines {
		res[i] = f(l)
	}

	return res
}

func Words(line string) []string {
	_words := strings.Split(line, " ")
	not_empty := func (w string) bool { return w != "" }
	return Filter(_words, not_empty)
}
