package jsonc

import (
	"encoding/json"
	"errors"
	"strings"

	d "github.com/teja2010/golconda/src/debug"
)

// All the chars which are of importance
const (
	NEWLINE         = "\n"
	SINGLE_LINE_CMT = `//`
	BLOCK_CMT_START = `/*`
	BLOCK_CMT_END   = `*/`
)

// Unmarshal a commented json
func Unmarshal(data []byte, v interface{}) error {
	strContent := string(data)

	err := blkCommentCheck(strContent)
	if err != nil {
		d.Error("blkCommentCheck failed", err)
		return err
	}

	lines := strings.Split(strContent, BLOCK_CMT_END)
	lines, err = fmapM(lines, removeBlkComments)
	if err != nil {
		d.Error("removeBlkComments failed", err)
		return err
	}
	strContent = strings.Join(lines, "")

	lines = strings.Split(strContent, NEWLINE)
	lines, err = fmapM(lines, removeSingleLineComment)
	if err != nil {
		d.Error("removeSingleLineComment failed", err)
		return err
	}

	lines = filter(lines, nonEmptyLines)

	data = []byte(strings.Join(lines, NEWLINE))
	d.DebugLog("Uncommented config \n", string(data))

	err = json.Unmarshal(data, v)
	if err != nil {
		d.Error("json Unmarshal failed", err)
	}

	return err

}

type fmapFunc func(string) (string, error)

// fmapM is not exactly fmap (structure is not preserved),
// it is a Monadic version of fmap
func fmapM(lines []string, f fmapFunc) ([]string, error) {
	res := make([]string, len(lines))

	for i, l := range lines {
		var err error
		res[i], err = f(l)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func filter(lines []string, f func(string) bool) []string {
	res := make([]string, len(lines))

	i := 0
	for _, l := range lines {
		if f(l) {
			res[i] = l
			i++
		}
	}

	return res
}

func removeBlkComments(l string) (string, error) {
	return trimSuffixStartsWith(l, BLOCK_CMT_START)
}

func removeSingleLineComment(l string) (string, error) {
	return trimSuffixStartsWith(l, SINGLE_LINE_CMT)
}

func trimSuffixStartsWith(l string, suffix string) (
	string, error) {

	splits := strings.SplitN(l, suffix, 2)
	if len(splits) == 1 {
		return splits[0], nil
	} else if len(splits) == 2 {
		l2 := splits[0]
		return l2, nil
	}

	return l,
		errors.New("trimSuffixStartsWith: SplitAfterN (" + suffix +
			") returned more than 2 elements :" + l)
}

func nonEmptyLines(l string) bool {
	return l != ""
}

func blkCommentCheck(str string) error {
	startCount := strings.Count(str, BLOCK_CMT_START)
	endCount := strings.Count(str, BLOCK_CMT_END)

	if endCount > startCount {
		return errors.New("More '*/'s found that matching '/*'s")
	}
	return nil
}
