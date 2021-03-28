package golconda

import (
	"regexp"
	"strings"

	d "github.com/teja2010/golconda/src/debug"
)

const (
	_NEWLINE = "\n"
)

// Regex2Func converts a regular expression into a function to match a string
func Regex2Func(rule string) func(string) bool {
	exp, err := regexp.Compile(rule)
	if d.Unlikely(err != nil) {
		d.Bug("Invalid regexp", rule)
	}

	f := func(l string) bool {
		return exp.MatchString(l)
	}

	return f
}

// Filter filters a list of strings based on a strings that match a condition
func Filter(lines []string, f func(string) bool) []string {

	matchedLen := 0
	matchedLines := make([]string, len(lines))

	for _, line := range lines {
		if f(line) {
			matchedLines[matchedLen] = line
			matchedLen++
		}
	}

	return matchedLines[:matchedLen]
}

// TakeWhile takes elements from a list till a condition is true
func TakeWhile(lines []string, f func(string) bool) []string {

	matchedLines := make([]string, len(lines))

	for i, line := range lines {
		if f(line) {
			matchedLines[i] = line
		} else {
			return matchedLines[:i]
		}
	}

	return matchedLines
}

// FindLine finds the first line that matches the condition
func FindLine(lines []string, f func(string) bool) string {
	for _, line := range lines {
		if f(line) {
			return line
		}
	}

	d.Bug("Unable to find a line matching", f)
	return "BUG!! BUG!!"
}

// clean-up the functions below later using generics

// FmapSS fmaps string to string
func FmapSS(lines []string, f func(string) string) []string {
	res := make([]string, len(lines))

	for i, l := range lines {
		res[i] = f(l)
	}

	return res
}

// FmapSI fmaps string to int
func FmapSI(lines []string, f func(string) int) []int {
	res := make([]int, len(lines))

	for i, l := range lines {
		res[i] = f(l)
	}

	return res
}

// FmapSI64 fmaps string to int64
func FmapSI64(lines []string, f func(string) int64) []int64 {
	res := make([]int64, len(lines))

	for i, l := range lines {
		res[i] = f(l)
	}

	return res
}

// FmapSCpuStat fmaps string to cpuStatData
func FmapSCpuStat(lines []string, f func(string) cpuStatData) []cpuStatData {
	res := make([]cpuStatData, len(lines))

	for i, l := range lines {
		res[i] = f(l)
	}

	return res
}

// Words converts a string block into lines
func Words(line string) []string {
	_words := strings.Split(line, " ")
	notEmpty := func(w string) bool { return w != "" }
	return Filter(_words, notEmpty)
}
