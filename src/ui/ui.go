package ui

var _ui UI

type Tuple struct {
	Fst, Snd int
}

type PrintData struct {
	Position Tuple
	Size Tuple
	Content []string
}

type UI interface {
	New() UI
	Update(PrintData)
}

// initialize a UI
func Init() {
	// TODO read config and decide which ui to choose
	_ui = StdoutUI{}.New()
}

// update the UI
func Update(pdata PrintData) {
	_ui.Update(pdata)
}
