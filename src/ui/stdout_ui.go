package ui

import (
	"fmt"
	"time"
	"strings"
	d "github.com/teja2010/golconda/src/debug"
)

// print everything out to stdout

type StdoutUI struct {
	last_update time.Time
	size Tuple
	content []string
}

func (s StdoutUI) String() string {
	return fmt.Sprintf("Time: %s\n%s",
			   time.Now().Format(time.Stamp),
			   format_content(s.size.Fst, s.size.Snd, s.content),
			)

}

func (s StdoutUI) New() UI {
	// nothing to Init
	d.DebugLog("Inited StdoutUI")

	return StdoutUI{time.Now(),
			Tuple{1, 100},
			[]string{"Waiting for data..."} }
}

func (s StdoutUI) Update(pdata PrintData) {
	// add pdata to s

	s.last_update = time.Now()
	s.size = pdata.Size
	s.content = pdata.Content

	fmt.Println(s)
}

func format_content(sizeX int, sizeY int, content []string) string {
	// format  the content so it fits into these dimensions
	return strings.Join(content, "\n")
}
