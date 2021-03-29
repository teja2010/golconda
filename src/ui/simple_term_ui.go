package ui

import "time"

// SimpleTerm to print data ti the terminal
type SimpleTermUI struct {
}

type simpleTermUIData struct {
	lastUpdate time.Time
	size       Tuple
	content    []string
}

// New create the SimpleTermUI
func (s SimpleTermUI) New() UI {
	return s
}

// Update sends data to SimpleTermUI
func (s SimpleTermUI) Update(pdata PrintData) {
}
