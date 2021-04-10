package golconda

import (
	_ "embed" // embed the default json
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	d "github.com/teja2010/golconda/src/debug"
	"github.com/teja2010/golconda/src/jsonc"
	"github.com/teja2010/golconda/src/meta"
	"github.com/teja2010/golconda/src/ui"
)

type GlobalConfig struct {
	UpdateInterval string
	DebugLevel     int
	UI             string
}
type GolcondaConfig struct {
	Global      GlobalConfig
	UI          ui.UIConfig
	CpuUsage    CPUUsageConfig
	MemInfo     MemInfoConfig
	ProcMemInfo ProcessMemConfig
}

// RegisteredFunction is the type which registered functions must support
type RegisteredFunction func(chan<- ui.PrintData)

// RegisteredFunctions returns functions which are registered
func RegisteredFunctions() []RegisteredFunction {
	return []RegisteredFunction{
		CPUUsage,
		MemInfo,
		ProcessMemoryInfo,
	}
}

// add specific functions above.

var _defaultConfig GolcondaConfig
var _config *GolcondaConfig
var rwlock sync.RWMutex
var _searchConfigFiles = []string{"~/.golcondarc"}

// GetConfig returns a pointer to the config
func GetConfig() *GolcondaConfig {
	rwlock.RLock()
	defer rwlock.RUnlock()
	return _config
}

// ConfigInit initializes config
func ConfigInit() {
	rwlock = sync.RWMutex{}
	readArgs()
	readDefaultConfig()

	configIsReady := make(chan bool)

	go updateConfigThread(configIsReady)

	<-configIsReady // wait for it to be read atleast once
}

func updateConfigThread(configIsReady chan bool) {
	// if file has changed update the config pointer
	var lastRead time.Time
	var once sync.Once
	for {
		data, file := findConfData()
		readConfig(data)
		once.Do(func() { configIsReady <- true })

		if file == "" { // the file in args is still not readable...
			time.Sleep(time.Second)
			continue
		}

		lastRead = time.Now()
		waitForChange(file, lastRead)
	}
}

//go:embed embeded_files/default_config.json
var _defConfigData []byte

// read Default config only once.
func readDefaultConfig() {
	err := jsonc.Unmarshal(_defConfigData, &_defaultConfig)
	if err != nil {
		// should never fail
		d.Bug("Default config unmarshall failed", err)
	}

	d.DebugLog("=======================================================")
	d.DebugLog("Read default Config", d.ToString(_defaultConfig))
	d.DebugLog("=======================================================")
}

// the only writer
func readConfig(data []byte) {
	rwlock.Lock()
	defer rwlock.Unlock()
	_config = new(GolcondaConfig)

	err := jsonc.Unmarshal(data, _config)
	if err != nil {
		d.Error("Unmarshall config data failed", err)
		return
	}
	d.DebugLog("Read Config from File", d.ToString(_config))

	err = meta.LeftMerge(&_defaultConfig, _config, &GolcondaConfig{})
	if err != nil {
		d.Error("LeftMerge failed", err)
		//return
	}

	d.DebugLog("=======================================================")
	d.DebugLog("Merged Config", d.ToString(_config))
	d.DebugLog("=======================================================")
}

func findConfData() ([]byte, string) {
	for _, cfile := range _searchConfigFiles {
		data, err := ioutil.ReadFile(cfile)
		if err == nil {
			d.Log("Read data from", cfile)
			return data, cfile
		}
		d.Error("findConfData", cfile, "failed", err)
	}

	// return an empty conf data
	return []byte("{}"), ""
}

func waitForChange(cfile string, lastRead time.Time) {
	for {
		if cfile == "" {
			time.Sleep(60 * time.Second)
			continue
		}

		time.Sleep(1 * time.Second)

		info, err := os.Stat(cfile)
		if err != nil {
			d.Error("Stat failed", err)
			return
		}

		if info.ModTime().After(lastRead) {
			break
		}
	}
}

// only these arguments supported.
// 1. -h or --help
// 2. -f or --follow-config-file
// 3. -s or --sample-config
func readArgs() {
	args := os.Args[:]

	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			printHelp()
		} else if args[i] == "-f" {
			if i+1 > len(args) {
				fmt.Fprintf(os.Stderr, "File name missing\n"+
					"Run golconda -h for help")
				os.Exit(-1)
			}
			_searchConfigFiles = []string{args[i+1]}
			d.DebugLog("Set config file", _searchConfigFiles)
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
	fmt.Println(string(_defConfigData))
}
