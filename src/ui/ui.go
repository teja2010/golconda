package ui

var _ui UI

// Tuple of ints
type Tuple struct {
	Fst, Snd int
}

// PrintData to print data send data to a UI
type PrintData struct {
	Position Tuple
	Size     Tuple
	Content  []string
}

// A UI interface
type UI interface {
	New() UI
	Update(PrintData)
	// Notify(NotifyData)
}

// Init - initialize a UI
func Init() {
	// TODO read config and decide which ui to choose
	_ui = StdoutUI{}.New()
}

// Update - update the UI
func Update(pdata PrintData) {
	_ui.Update(pdata)
}
