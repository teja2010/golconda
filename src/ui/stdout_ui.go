package ui

import (
	"fmt"
	"strings"
	"time"

	d "github.com/teja2010/golconda/src/debug"
)

// Simply print everything out to stdout

//StdoutUI is an empty struct to satisfy the UI interface
type StdoutUI struct {
}

// stdoutUIData data to print out to stdout
type stdoutUIData struct {
	lastUpdate time.Time
	size       Tuple
	content    []string
}

func (s stdoutUIData) String() string {
	return fmt.Sprintf("Time: %s\n%s",
		time.Now().Format(time.Stamp),
		s.formatContent(),
	)

}

// New create a new UI
func (s StdoutUI) New() UI {
	// nothing to Init
	d.DebugLog("Inited StdoutUI")

	return s
}

// Update sends data to UI
func (s StdoutUI) Update(pdata PrintData) {
	sd := stdoutUIData{
		lastUpdate: time.Now(),
		size:       pdata.Size,
		content:    pdata.Content,
	}

	fmt.Println(sd)
}

func (s stdoutUIData) formatContent() string {
	return strings.Join(s.content, "\n")
}
