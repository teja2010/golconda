package golconda

import (
	"os"
	"fmt"
	"time"
	"sync"
	_ "embed"
	d "github.com/teja2010/golconda/src/debug"
	"github.com/teja2010/golconda/src/ui"
	"github.com/teja2010/golconda/src/jsonc"
)

type GlobalConfig struct {
	UpdateInterval string
}
type GolcondaConfig struct {
	Global GlobalConfig
	CpuUsage CpuUsageConfig
	MemInfo MemInfoConfig
}

type RegisteredFunction func(chan<- ui.PrintData)
func RegisteredFunctions() []RegisteredFunction {
	return []RegisteredFunction{
			CPU_Usage,
			Meminfo,
		}
}

// add specific functions above.

var _default_config *GolcondaConfig
var _config *GolcondaConfig
var rwlock sync.RWMutex
var config_file string = ""
var _search_config_files = []string{"~/.golcondarc"}

// return config
func GetConfig() *GolcondaConfig {
	return _config
}

func ConfigInit() {
	rwlock = sync.RWMutex{}
	readArgs()
	readDefaultConfig()
	readConfig()

	go updateConfigThread()
}

func updateConfigThread() {
	// if file has changed update the config pointer
	for {
		time.Sleep(1*time.Second)
		readConfig()
	}
}

//go:embed embeded_files/default_config.json
var _def_config_data []byte
// read Default config only once.
func readDefaultConfig() {
	_default_config := new(GolcondaConfig)
	err := jsonc.Unmarshal(_def_config_data, _default_config)
	if err != nil {
		// should never fail
		d.Bug("Default config unmarshall failed", err)
	}

	d.DebugLog(fmt.Sprintf("%+v", _default_config))
}

// the only writer
func readConfig() {
	conf := new(GolcondaConfig)

	conf.CpuUsage.UpdateInterval = "3s"
	conf.MemInfo.UpdateInterval = "1s"

	rwlock.Lock()
	_config = conf
	rwlock.Unlock()
}

// only these arguments supported.
// 1. -h or --help
// 2. -f or --follow-config-file
// 3. -s or --sample-config
func readArgs() {
	args := os.Args[:]

	for i := 0 ; i < len(args) ; i++ {
		if args[i] == "-h" || args[i] == "--help" {
			printHelp()
		} else if args[i] == "-f" {
			if i+1 > len(args) {
				fmt.Fprintf(os.Stderr, "File name missing\n" +
						"Run golconda -h for help")
				os.Exit(-1)
			}
			config_file = args[i+1]
			d.DebugLog("Set config file", config_file)
			i++
		} else if args[i] == "-d" {
			printDefaultConfig()
		}
	}
}

func printHelp() {
	defer os.Exit(0)

	fmt.Println("" +
`Golconda 0.0.1

USAGE:
  -f, --follow-config-file FILE_PATH
      Follow config file to configure Golconda. When values are updated in
      the file, they will reflect automatically.
  -d, --default-config
      Print default config file to STDOUT. Use it to build custom config files.
  -h
      Print this help :)
`)
}

func printDefaultConfig() {
	defer os.Exit(0)
	fmt.Println(string(_def_config_data))
}
