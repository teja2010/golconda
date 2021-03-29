package meta

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	d "github.com/teja2010/golconda/src/debug"
)

const (
	_SPACE        = " "
	_TOKEN_REGEXP = `<[A-Za-z0-9_]*( %[A-Za-z0-9#\-+ .]*)?>`
	// The variable name followed by a formatting verb
	// i.e. <Var> or <Var %v> or <Real %6.3f>
)

// Format - given a string, print the data into a string.
// e.g.
//	type ComplexNum struct {
//		Real, Img int
//	}
//
// then "<Real> + <Img %03d>i"  should format ComplexNum{1,2} as "1 + 2i"
// For Now, it only can handle a struct without multiple levels
func Format(s string, val interface{}) string {

	args := []interface{}{}
	tokens := splitIntoTokens(s)

	_val := reflect.ValueOf(val)

	for i, w := range tokens {
		if !(strings.HasPrefix(w, "<") && strings.HasSuffix(w, ">")) {
			continue
		}

		varName, fmtVerb := extractVarFmtInfo(w)

		_el := _val.FieldByName(varName)
		if !_el.IsValid() {
			continue
		}

		tokens[i] = fmtVerb
		args = append(args, _el.Interface())
	}

	fmtStr := strings.Join(tokens, "")
	d.DebugLog("fmtStr", fmtStr)

	return fmt.Sprintf(fmtStr, args...)
}

func splitIntoTokens(_s string) []string {
	s := _s[:]

	tokens := []string{}

	pattern := regexp.MustCompile(_TOKEN_REGEXP)
	indexPairs := pattern.FindAllStringIndex(s, -1)

	rOld := 0
	for _, pair := range indexPairs {
		l, r := pair[0], pair[1]
		tokens = append(tokens, s[rOld:l])
		tokens = append(tokens, s[l:r])
		rOld = r
	}
	tokens = append(tokens, s[rOld:])

	d.DebugLog("Tokens", `["`+strings.Join(tokens, `", "`)+`"]`)

	return tokens
}

func extractVarFmtInfo(w string) (string, string) {
	w = w[1 : len(w)-1]

	idx := strings.Index(w, " %")
	if idx == -1 {
		return w, "%v"
	}

	return w[:idx], strings.TrimLeft(w[idx:], " ")
}
