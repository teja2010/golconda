package debug

// checks which will be run only when the debug flag is set

const (
	_DEBUG = true
)

// Unlikely are checks that are mostly false. They should be removed when
// _DEBUG is false
func Unlikely(result bool) bool {
	if _DEBUG {
		if result {
			Bug("Debug Check failed")
		}
		return result
	}

	return false
}
