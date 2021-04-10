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
)

// CPUUsageConfig config for cpu usage
type CPUUsageConfig struct {
	UpdateInterval string
	UIPosition     ui.Tuple
	UISize         ui.Tuple
	FmtString      []string
	PerCPUStatFmt  string
	CPUs           string
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

func _CPUUsage(c chan<- ui.PrintData, oldData cpuStatData,
	conf CPUUsageConfig) cpuStatData {

	newData := readProcStat()
	diffData := newData.makeDiff(oldData, conf)

	pdata := ui.PrintData{
		Position: ui.Tuple{X: 5, Y: 0},
		Size:     ui.Tuple{X: 10, Y: 100},
		Content:  []string{},
	}

	for _, fmtStr := range conf.FmtString {
		if fmtStr == "[PerCPUStatFmt]" {
			pdata.Content = append(pdata.Content,
				diffData.PerCPUStatFmt...)
		} else {
			contentStr := meta.Format(fmtStr, diffData)
			pdata.Content = append(pdata.Content, contentStr)
		}
	}

	// finally push into the channel
	c <- pdata
	return newData
}

func (newData cpuStatData) makeDiff(oldData cpuStatData,
	conf CPUUsageConfig) cpuStatData {

	diffData := cpuStatData{}
	NumOfCPUs := len(newData.stats)

	floatDiv := func(a, b int64) float32 {
		return 100.0 * float32(a) / float32(b)
	}

	for i := 0; i < NumOfCPUs; i++ {
		d.DebugLog("TTT", i, len(newData.stats))
		nd := newData.stats[i]
		od := oldData.stats[i]

		diff := perCpuStatData{
			Title: nd.Title,
			User:  nd.User - od.User,
			Kern:  nd.Kern - od.Kern,
			Idle:  nd.Idle - od.Idle,
			Irq:   nd.Irq - od.Irq,
			Guest: nd.Guest - od.Guest,
		}

		diff.Sum = diff.User + diff.Kern + diff.Idle + diff.Irq + diff.Guest
		diff.UserPercent = floatDiv(diff.User, diff.Sum)
		diff.KernPercent = floatDiv(diff.Kern, diff.Sum)
		diff.IdlePercent = floatDiv(diff.Idle, diff.Sum)
		diff.IrqPercent = floatDiv(diff.Irq, diff.Sum)
		diff.GuestPercent = floatDiv(diff.Guest, diff.Sum)

		diffData.stats = append(diffData.stats, diff)
	}

	rangeInt := func(start, end int) []int {
		ret := []int{}
		for i := start; i < end; i++ {
			ret = append(ret, i)
		}
		return ret
	}

	parseCPUs := func(cpus string, NumOfCPUs int) []int {
		// "overall,1,2,5,2" -> []int{0, 2, 3, 6, 3}
		if cpus == "all" {
			return rangeInt(0, NumOfCPUs)
		}

		ret := []int{}
		for _, n := range strings.Split(cpus, ",") {
			n = strings.TrimSpace(n)
			if n == "overall" {
				ret = append(ret, 0)
			} else {
				num, err := strconv.Atoi(n)
				if err == nil && num+1 < NumOfCPUs {
					ret = append(ret, num+1)
				}
			}
		}
		return ret
	}

	for _, statsIdx := range parseCPUs(conf.CPUs, NumOfCPUs) {
		diffData.PerCPUStatFmt = append(
			diffData.PerCPUStatFmt,
			meta.Format(conf.PerCPUStatFmt,
				diffData.stats[statsIdx]))
	}

	return diffData
}

type cpuStatData struct {
	stats         []perCpuStatData
	PerCPUStatFmt []string
}

type perCpuStatData struct {
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

func readProcStat() cpuStatData {
	_contents, err := ioutil.ReadFile(_PROC_STAT)
	if err != nil {
		d.Error("Unable to read", _PROC_STAT)
		return cpuStatData{}
	}

	contents := string(_contents)
	lines := strings.Split(contents, _NEWLINE)

	cpuLines := TakeWhile(lines, Regex2Func("^cpu"))

	usageData := cpuStatData{
		stats: FmapSCpuStat(cpuLines, readUsageLine),
	}

	return usageData
}

func readUsageLine(line string) perCpuStatData {
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

	return perCpuStatData{
		Title: strings.ToUpper(title),
		User:  userUsage,
		Kern:  kernUsage,
		Idle:  idleUsage,
		Irq:   irqUsage,
		Guest: guestUsage,
	}
}
