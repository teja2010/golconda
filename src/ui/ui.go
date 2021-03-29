package ui

var _ui UI

// Tuple of ints
type Tuple struct {
	X, Y int
}

// PrintData to print data send data to a UI
type PrintData struct {
	Position Tuple
	Size     Tuple
	Content  []string
}

type UIConfig struct {
	UI         string
	SimpleTerm SimpleTermConfig
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
	//_ui = StdoutUI{}.New()
	_ui = SimpleTermUI{}.New()
}

// Update - update the UI
func Update(pdata PrintData) {
	_ui.Update(pdata)
}
