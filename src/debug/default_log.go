package debug

import (
	"bufio"
	"log"
	"os"
)

const (
	_GOLCONDA_LOG_FILE = "/tmp/golconda.log"
	_LOG_FLAGS         = (log.Ltime | log.Lmicroseconds |
			      log.Lmsgprefix)
)

type def_log struct {
	_writer *bufio.Writer
	_logger *log.Logger
}

func DefLogInit() def_log {
	d := def_log{nil, nil}

	// Nothing
	file, err := os.OpenFile(_GOLCONDA_LOG_FILE,
				 os.O_APPEND|os.O_WRONLY|os.O_CREATE,
				 0644)

	if err != nil {
		Bug("Error Unable to open temp file", err)
	}

	//d._writer = bufio.NewWriter(file)

	d._logger = log.New(file, "", _LOG_FLAGS)

	return d
}

func (d def_log) Println(v ...interface{}) {

	d._logger.Println(v...)
}
