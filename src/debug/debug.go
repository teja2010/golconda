package debug

// checks which will be run only when the debug flag is set

const (
	_DEBUG = false
)

// These checks must mostly be false
func DebugCheck (result bool) bool {
	if _DEBUG {
		if result {
			Bug("Debug Check failed")
		}
		return result
	} else {
		return false
	}
}
