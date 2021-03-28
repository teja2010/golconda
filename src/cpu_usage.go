package golconda

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"io/ioutil"
	ui "github.com/teja2010/golconda/src/ui"
	d "github.com/teja2010/golconda/src/debug"
)

const (
	_PROC_STAT = "/proc/stat"

	// regexps
	_HEADER_CPU_USAGE = "CPU USAGE:"
)

type CpuUsageConfig struct {
	UpdateInterval string
}

// read values from /proc/stat
func CPU_Usage(c chan<- ui.PrintData) {

	old_data := read_proc_stat()

	for {
		conf := GetConfig()

		update_interval := confCpuUpdateInterval(conf)
		duration, err := time.ParseDuration(update_interval)
		if err != nil {
			d.Error("Invalid Duration:", update_interval)
			duration = 1*time.Second
			//TODO read this value from the default config
		}
		time.Sleep(duration)

		old_data = __cpu_usage(c, old_data)
	}

}

func confCpuUpdateInterval(conf *GolcondaConfig) string {
	update_interval := conf.CpuUsage.UpdateInterval
	if update_interval == "" {
		update_interval = conf.Global.UpdateInterval
	}

	return update_interval
}

func __cpu_usage(c chan<- ui.PrintData, old_data []cpu_stat_data) []cpu_stat_data {

	new_data := read_proc_stat()

	pdata := ui.PrintData{
				ui.Tuple{0, 0},
				ui.Tuple{1, 100},
				[]string{_HEADER_CPU_USAGE},
			}

	for i := 0; i < len(old_data); i++ {
		fmt_usage_data := fmt_cpu_usage(new_data[i], old_data[i])
		pdata.Content = append(pdata.Content, fmt_usage_data)
	}

	// finally push into the channel
	c <- pdata

	d.DebugLog(pdata.Content[:3])
	return new_data
}

func fmt_cpu_usage(new_data, old_data cpu_stat_data) string {

	diff := cpu_stat_data{
		new_data.title,
		new_data.user - old_data.user,
		new_data.kern - old_data.kern,
		new_data.idle - old_data.idle,
		new_data.irq - old_data.irq,
		new_data.guest - old_data.guest,
	}

	sum := diff.user + diff.kern + diff.idle + diff.irq + diff.guest

	float_div := func(a, b int64) float32 {
		return 100.0*float32(a)/float32(b)
	}

	fmt_string := fmt.Sprintf("%s  User %6.02f%% | Kern %6.02f%% | " +
				  "Idle %6.02f%% | Irq %6.02f%% | " +
				  "Guest %6.02f%%",
				strings.ToUpper(diff.title),
				float_div(diff.user, sum),
				float_div(diff.kern, sum),
				float_div(diff.idle, sum),
				float_div(diff.irq, sum),
				float_div(diff.guest, sum),
			)

	return fmt_string
}

type cpu_stat_data struct {
	title string
	user int64
	kern int64
	idle int64
	irq int64
	guest int64
}

func read_proc_stat() []cpu_stat_data {
	_contents, err := ioutil.ReadFile(_PROC_STAT)
	if err != nil {
		d.Error("Unable to read", _PROC_STAT)
		return []cpu_stat_data{}
	}

	contents := string(_contents)
	lines := strings.Split(contents, _NEWLINE)

	cpu_lines := TakeWhile(lines, Regex2Func("^cpu"))

	usage_data := FmapSCpu_stat(cpu_lines, read_usage_line)

	return usage_data
}

func read_usage_line(line string) cpu_stat_data {
	_words := strings.Split(line, " ")

	not_empty := func (l string) bool { return l != "" }
	words := Filter(_words, not_empty)

	if d.DebugCheck(len(words) != 11) {
		d.Bug("We must have 11 words")
	}

	title := words[0]
	if title == "cpu" {
		title += " " // add padding so they align
	}

	convert2int := func(w string) int64 {
		i64, err := strconv.ParseInt(w, 10, 64)
		if err != nil { d.Bug("ParseInt failed") }
		return i64
	}
	usage_ints := FmapSI64(words[1:], convert2int)
	user_usage := usage_ints[0] + usage_ints[1]
	kern_usage := usage_ints[2]
	idle_usage := usage_ints[3]
	// 4 is IO wait
	irq_usage := usage_ints[5] + usage_ints[6] //hardirq + softirq
	// 7 is steal
	guest_usage := usage_ints[8] + usage_ints[9]

	return cpu_stat_data {
		title:title,
		user:user_usage,
		kern:kern_usage,
		idle:idle_usage,
		irq:irq_usage,
		guest:guest_usage,
	}
}

