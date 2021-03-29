package golconda

import (
	//"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	d "github.com/teja2010/golconda/src/debug"
	"github.com/teja2010/golconda/src/meta"
	ui "github.com/teja2010/golconda/src/ui"
)

const (
	_PROC_STAT = "/proc/stat"

	_HEADER_CPU_USAGE = "CPU USAGE:"
)

// CPUUsageConfig config for cpu usage
type CPUUsageConfig struct {
	UpdateInterval string
	UIPosition     ui.Tuple
	UISize         ui.Tuple
	FmtString      string
}

// CPUUsage reads values from /proc/stat to display cpu stats
func CPUUsage(c chan<- ui.PrintData) {

	oldData := readProcStat()

	for {
		conf := GetConfig()

		updateInterval := confCPUUpdateInterval(conf)
		duration, err := time.ParseDuration(updateInterval)
		if err != nil {
			d.Error("Invalid Duration:", updateInterval)
			duration = 1 * time.Second
			//TODO read this value from the default config
		}
		time.Sleep(duration)

		oldData = _CPUUsage(c, oldData, conf.CpuUsage)
	}

}

func confCPUUpdateInterval(conf *GolcondaConfig) string {
	updateInterval := conf.CpuUsage.UpdateInterval
	if updateInterval == "" {
		updateInterval = conf.Global.UpdateInterval
	}

	return updateInterval
}

func _CPUUsage(c chan<- ui.PrintData, oldData []cpuStatData,
	conf CPUUsageConfig) []cpuStatData {

	newData := readProcStat()

	pdata := ui.PrintData{
		Position: ui.Tuple{X: 5, Y: 0},
		Size:     ui.Tuple{X: 10, Y: 100},
		Content:  []string{_HEADER_CPU_USAGE},
	}

	for i := 0; i < len(oldData); i++ {
		fmtUsageData := fmtCPUUsage(newData[i], oldData[i],
			conf.FmtString)
		pdata.Content = append(pdata.Content, fmtUsageData)
	}

	// finally push into the channel
	c <- pdata

	d.DebugLog(pdata.Content[:3])
	return newData
}

func fmtCPUUsage(newData, oldData cpuStatData, confFmt string) string {

	diff := cpuStatData{
		Title: newData.Title,
		User:  newData.User - oldData.User,
		Kern:  newData.Kern - oldData.Kern,
		Idle:  newData.Idle - oldData.Idle,
		Irq:   newData.Irq - oldData.Irq,
		Guest: newData.Guest - oldData.Guest,
	}

	diff.Sum = diff.User + diff.Kern + diff.Idle + diff.Irq + diff.Guest

	floatDiv := func(a, b int64) float32 {
		return 100.0 * float32(a) / float32(b)
	}

	diff.UserPercent = floatDiv(diff.User, diff.Sum)
	diff.KernPercent = floatDiv(diff.Kern, diff.Sum)
	diff.IdlePercent = floatDiv(diff.Idle, diff.Sum)
	diff.IrqPercent = floatDiv(diff.Irq, diff.Sum)
	diff.GuestPercent = floatDiv(diff.Guest, diff.Sum)

	fmtString := meta.Format(confFmt, diff)

	return fmtString
}

type cpuStatData struct {
	Title        string
	User         int64
	Kern         int64
	Idle         int64
	Irq          int64
	Guest        int64
	Sum          int64
	UserPercent  float32
	KernPercent  float32
	IdlePercent  float32
	IrqPercent   float32
	GuestPercent float32
}

func readProcStat() []cpuStatData {
	_contents, err := ioutil.ReadFile(_PROC_STAT)
	if err != nil {
		d.Error("Unable to read", _PROC_STAT)
		return []cpuStatData{}
	}

	contents := string(_contents)
	lines := strings.Split(contents, _NEWLINE)

	cpuLines := TakeWhile(lines, Regex2Func("^cpu"))

	usageData := FmapSCpuStat(cpuLines, readUsageLine)

	return usageData
}

func readUsageLine(line string) cpuStatData {
	_words := strings.Split(line, " ")

	notEmpty := func(l string) bool { return l != "" }
	words := Filter(_words, notEmpty)

	if d.Unlikely(len(words) != 11) {
		d.Bug("We must have 11 words")
	}

	title := words[0]
	if title == "cpu" {
		title += " " // add padding so they align
	}

	convert2int := func(w string) int64 {
		i64, err := strconv.ParseInt(w, 10, 64)
		if err != nil {
			d.Bug("ParseInt failed")
		}
		return i64
	}
	usageInts := FmapSI64(words[1:], convert2int)
	userUsage := usageInts[0] + usageInts[1]
	kernUsage := usageInts[2]
	idleUsage := usageInts[3]
	// 4 is IO wait
	irqUsage := usageInts[5] + usageInts[6] //hardirq + softirq
	// 7 is steal
	guestUsage := usageInts[8] + usageInts[9]

	return cpuStatData{
		Title: strings.ToUpper(title),
		User:  userUsage,
		Kern:  kernUsage,
		Idle:  idleUsage,
		Irq:   irqUsage,
		Guest: guestUsage,
	}
}
