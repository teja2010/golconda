package ui

import (
	"fmt"
	d "github.com/teja2010/golconda/src/debug"
	"strings"
)

const (
	CLR_SCREEN = "\033[2J"
)

type SimpleTermConfig struct {
	Size Tuple
}

// SimpleTermUI to print data ti the terminal
type SimpleTermUI struct {
	sizeX, sizeY int
	uiContent    []string
}

func (s SimpleTermUI) String() string {
	return strings.Join(s.uiContent, "\n")
}

// New create the SimpleTermUI
func (s SimpleTermUI) New() UI {
	s.sizeX = 30
	s.sizeY = 120

	s.uiContent = make([]string, s.sizeX)
	for row := 0; row < s.sizeX; row++ {
		s.uiContent[row] = padTo("", s.sizeY)
	}

	return s
}

func gotoPosition(x, y int) string {
	return fmt.Sprintf("\033[%d;%df", x, y)
}

// Update sends data to SimpleTermUI
func (s SimpleTermUI) Update(pdata PrintData) {
	startRow := (pdata.Position.X)
	endRow := (startRow + pdata.Size.X)

	startCol := (pdata.Position.Y)
	endCol := (startCol + pdata.Size.Y)

	paddedContent := simpleTermUIpad(pdata)
	d.DebugLog(paddedContent)

	for row := startRow; row < endRow; row++ {
		s.uiContent[row] = (s.uiContent[row][:startCol] +
			paddedContent[row-startRow] +
			s.uiContent[row][endCol:])
	}

	d.DebugLog(s.uiContent)

	fmt.Printf(CLR_SCREEN + gotoPosition(0, 0))
	fmt.Print(s.String())
}

func simpleTermUIpad(pdata PrintData) []string {
	paddedData := []string{}

	for i := 0; i < pdata.Size.X; i++ {
		if i < len(pdata.Content) {
			paddedData = append(paddedData,
				padTo(pdata.Content[i], pdata.Size.Y))
		} else {
			paddedData = append(paddedData,
				padTo("", pdata.Size.Y))
		}

	}

	return paddedData
}

func padTo(s string, l int) string {
	if l < len(s) {
		return s[:l]
	}

	widthFmt := "%-" + fmt.Sprintf("%d", l) + "s"

	return fmt.Sprintf(widthFmt, s)
}
