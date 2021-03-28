package debug

import (
	"log"
	"os"
)

const (
	_GolcondaLogFile = "/tmp/golconda.log"
	_LogFlags        = (log.Ltime | log.Lmicroseconds | log.Lmsgprefix)
)

type defLogData struct {
	_logger *log.Logger
}

// DefLogInit inits logs
func defLogInit() defLogData {
	d := defLogData{nil}

	// Nothing
	file, err := os.OpenFile(_GolcondaLogFile,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0644)

	if err != nil {
		Bug("Error Unable to open temp file", err)
	}

	d._logger = log.New(file, "", _LogFlags)

	return d
}

func (d defLogData) Println(v ...interface{}) {

	d._logger.Println(v...)
}
