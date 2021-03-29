package golconda

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	d "github.com/teja2010/golconda/src/debug"
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

		oldData = _CPUUsage(c, oldData)
	}

}

func confCPUUpdateInterval(conf *GolcondaConfig) string {
	updateInterval := conf.CpuUsage.UpdateInterval
	if updateInterval == "" {
		updateInterval = conf.Global.UpdateInterval
	}

	return updateInterval
}

func _CPUUsage(c chan<- ui.PrintData, oldData []cpuStatData) []cpuStatData {

	newData := readProcStat()

	pdata := ui.PrintData{
		Position: ui.Tuple{X: 5, Y: 0},
		Size:     ui.Tuple{X: 10, Y: 100},
		Content:  []string{_HEADER_CPU_USAGE},
	}

	for i := 0; i < len(oldData); i++ {
		fmtUsageData := fmtCPUUsage(newData[i], oldData[i])
		pdata.Content = append(pdata.Content, fmtUsageData)
	}

	// finally push into the channel
	c <- pdata

	d.DebugLog(pdata.Content[:3])
	return newData
}

func fmtCPUUsage(newData, oldData cpuStatData) string {

	diff := cpuStatData{
		newData.title,
		newData.user - oldData.user,
		newData.kern - oldData.kern,
		newData.idle - oldData.idle,
		newData.irq - oldData.irq,
		newData.guest - oldData.guest,
	}

	sum := diff.user + diff.kern + diff.idle + diff.irq + diff.guest

	floatDiv := func(a, b int64) float32 {
		return 100.0 * float32(a) / float32(b)
	}

	fmtString := fmt.Sprintf("%s  User %6.02f%% | Kern %6.02f%% | "+
		"Idle %6.02f%% | Irq %6.02f%% | "+
		"Guest %6.02f%%",
		strings.ToUpper(diff.title),
		floatDiv(diff.user, sum),
		floatDiv(diff.kern, sum),
		floatDiv(diff.idle, sum),
		floatDiv(diff.irq, sum),
		floatDiv(diff.guest, sum),
	)

	return fmtString
}

type cpuStatData struct {
	title string
	user  int64
	kern  int64
	idle  int64
	irq   int64
	guest int64
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
		title: title,
		user:  userUsage,
		kern:  kernUsage,
		idle:  idleUsage,
		irq:   irqUsage,
		guest: guestUsage,
	}
}
