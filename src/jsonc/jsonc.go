package jsonc

import (
	"strings"
	"encoding/json"
	"errors"
	 d "github.com/teja2010/golconda/src/debug"
)

const (
	NEWLINE = "\n"
	SINGLE_LINE_CMT = `//`
	BLOCK_CMT_START = `/*`
	BLOCK_CMT_END = `*/`
)

func Unmarshal(data []byte, v interface{}) error {
	str_content := string(data)
	
	err := blkCommentCheck(str_content)
	if err != nil {
		d.Error("blkCommentCheck failed", err)
		return err
	}

	lines := strings.Split(str_content, BLOCK_CMT_END)
	lines, err = fmapM(lines, remove_blk_comments)
	if err != nil {
		d.Error("remove_blk_comments failed", err)
		return err
	}
	str_content = strings.Join(lines, "")

	lines = strings.Split(str_content, NEWLINE)
	lines, err = fmapM(lines, remove_single_line_comment)
	if err != nil {
		d.Error("remove_single_line_comment failed", err)
		return err
	}

	lines = filter(lines, non_empty_lines)

	data = []byte(strings.Join(lines, NEWLINE))
	//d.DebugLog("Uncommented config \n", string(data))

	err = json.Unmarshal(data, v)
	if err != nil {
		d.Error("json Unmarshal failed", err)
	}

	return err

}

type fmapFunc func(string) (string, error)

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

func remove_blk_comments(l string) (string, error) {
	return trimSuffixStartsWith(l, BLOCK_CMT_START)

}

func remove_single_line_comment(l string) (string, error) {
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

	return l, errors.New("trimSuffixStartsWith: SplitAfterN (" + suffix +
			     ") returned more than 2 elements :" + l)
}

func non_empty_lines (l string) bool {
	return l != ""
}

func blkCommentCheck(str string) error {
	start_count := strings.Count(str, BLOCK_CMT_START)
	end_count := strings.Count(str, BLOCK_CMT_END)

	if end_count > start_count {
		return errors.New("More '*/'s found that matching '/*'s")
	}
	return nil
}
