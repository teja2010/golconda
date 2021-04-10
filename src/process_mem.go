package golconda

import (
	d "github.com/teja2010/golconda/src/debug"
	"github.com/teja2010/golconda/src/meta"
	"github.com/teja2010/golconda/src/ui"
	"io/fs"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	_PMEM_PROC   = "/proc/"
	_PMEM_STATUS = "/status"
)

type ProcessMemConfig struct {
	UpdateInterval      string
	UIPosition          ui.Tuple
	UISize              ui.Tuple
	FmtString           []string
	PerProcessFmtString string
	SortBy              string
	Top, Bottom         int
}

// ProcessMemoryInfo - per process info
func ProcessMemoryInfo(c chan<- ui.PrintData) {

	for {
		conf := GetConfig()

		updateInterval := confProcUpdateInterval(conf)
		duration, err := time.ParseDuration(updateInterval)
		if err != nil {
			d.Bug("Invalid Duration:", updateInterval)
			duration = 1 * time.Second
			// TODO read this value from the default config
		}
		time.Sleep(duration)

		_processInfo(c, conf.ProcMemInfo)
	}
}

func confProcUpdateInterval(conf *GolcondaConfig) string {
	updateInterval := conf.ProcMemInfo.UpdateInterval
	if updateInterval == "" {
		updateInterval = conf.Global.UpdateInterval
	}

	return updateInterval
}

type procMemStat struct {
	data []perProcessStat
}

type perProcessStat struct {
	Name   string
	Pid    string
	VmRSS  int64
	VmHWM  int64
	VmSwap int64
}

func _processInfo(c chan<- ui.PrintData, conf ProcessMemConfig) {
	processArr := _readProcessData()

	sortedProcessArr := procStatSort(processArr,
		conf.SortBy, conf.Top, conf.Bottom)

	d.DebugLog("sorted", d.ToString(sortedProcessArr))

	content := []string{}
	for _, fmtStr := range conf.FmtString {
		if fmtStr == "[PerProcessFmtString]" {
			for _, pstat := range sortedProcessArr {
				str := meta.Format(conf.PerProcessFmtString,
					pstat)
				content = append(content, str)
			}
		} else {
			// TODO: not configurable...
			content = append(content, fmtStr)
		}
	}

	pdata := ui.PrintData{
		Position: conf.UIPosition,
		Size:     conf.UISize,
		Content:  content,
	}

	c <- pdata
	d.DebugLog("sent", d.ToString(pdata))
}

type perProcSort []perProcessStat

var ProcArrSortBy string
var ProcArrMut = sync.Mutex{}

func (a perProcSort) Len() int      { return len(a) }
func (a perProcSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a perProcSort) Less(i, j int) bool {
	return meta.ALessThanB(a[i], a[j], ProcArrSortBy)
}

func procStatSort(processArr []perProcessStat,
	SortBy string, Top int, Bottom int) []perProcessStat {

	if meta.Contains(perProcessStat{}, SortBy) == false {
		return []perProcessStat{}
	}

	if Top > len(processArr) {
		Top = len(processArr)
	}

	if Bottom > len(processArr) {
		Bottom = len(processArr)
	}

	ProcArrMut.Lock()
	ProcArrSortBy = SortBy
	sort.Sort(sort.Reverse(perProcSort(processArr)))
	ProcArrMut.Unlock()

	return append(processArr[:Top], processArr[len(processArr)-Bottom:]...)
}

func _readProcessData() []perProcessStat {

	files, err := ioutil.ReadDir(_PMEM_PROC)
	if err != nil {
		return []perProcessStat{}
	}

	_processFiles := FilterFileInfo(files, filterPid)
	processFiles := FmapFilePerProcessData(_processFiles, readProcessFiles)

	notEmpty := func(s string) bool { return s != "" }

	nonEmptyProcessFiles := Filter(processFiles, notEmpty)

	return FmapPerProcessData(nonEmptyProcessFiles, extractperProcessStat)
}

func extractperProcessStat(finfo string) perProcessStat {

	lines := strings.Split(finfo, _NEWLINE)

	getStr := func(prefix string) string {
		s3 := TryFindLine(lines, Regex2Func("^"+prefix), "0 kB")
		s2 := strings.TrimPrefix(s3, prefix)
		s := strings.TrimSpace(s2)
		return s
	}

	getMVal := func(prefix string) int64 {
		s2 := getStr(prefix)
		s1 := strings.TrimSuffix(s2, " kB")
		s := strings.TrimSpace(s1)
		i64, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			d.Bug("ParseInt failed")
		}
		return i64
	}

	ret := perProcessStat{}

	ret.Name = getStr("Name:")
	ret.Pid = getStr("Pid:")
	ret.VmRSS = getMVal("VmRSS:")
	ret.VmHWM = getMVal("VmHWM:")
	ret.VmSwap = getMVal("VmSwap:")

	return ret
}

func filterPid(finfo fs.FileInfo) (string, bool) {
	name := finfo.Name()
	onlyHasDigits := Regex2Func(`^[0-9]+$`) // only digits
	return name, onlyHasDigits(name)
}

func readProcessFiles(file string) string {
	contents, err := ioutil.ReadFile(_PMEM_PROC + file + _PMEM_STATUS)
	if err != nil {
		return ""
	}

	return string(contents)
}
