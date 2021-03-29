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

// /proc/meminfo
const (
	_PROC_MEMINFO = "/proc/meminfo"

	_HEADER_MEMINFO = "Memory Info:"
)

// MemInfoConfig to read mem info
type MemInfoConfig struct {
	UpdateInterval string
	UIPosition     ui.Tuple
	UISize         ui.Tuple
}

// MemInfo reg. func to read mem info
func MemInfo(c chan<- ui.PrintData) {

	for {
		conf := GetConfig()

		updateInterval := confMemUpdateInterval(conf)
		duration, err := time.ParseDuration(updateInterval)
		if err != nil {
			d.Bug("Invalid Duration:", updateInterval)
			duration = 1 * time.Second
			// TODO read this value from the default config
		}
		time.Sleep(duration)

		_memInfo(c)
	}
}

func confMemUpdateInterval(conf *GolcondaConfig) string {
	updateInterval := conf.MemInfo.UpdateInterval
	if updateInterval == "" {
		updateInterval = conf.Global.UpdateInterval
	}

	return updateInterval
}

func _memInfo(c chan<- ui.PrintData) {
	_contents, err := ioutil.ReadFile(_PROC_MEMINFO)
	if err != nil {
		d.Error("Unable to read", _PROC_MEMINFO)
		return
	}

	contents := string(_contents)
	lines := strings.Split(contents, _NEWLINE)

	getMVal := func(prefix string) int64 {
		s3 := FindLine(lines, Regex2Func("^"+prefix))
		s2 := strings.TrimPrefix(s3, prefix)
		s1 := strings.TrimSuffix(s2, " kB")
		s := strings.TrimSpace(s1)

		i64, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			d.Bug("ParseInt failed")
		}
		return i64
	}

	totalMem := getMVal("MemTotal:")
	freeMem := getMVal("MemFree:")
	availMem := getMVal("MemAvailable:")
	usedMem := totalMem - availMem
	cacheMem := getMVal("Cached:")
	sharedMem := getMVal("Shmem:")

	activeMem := getMVal("Active:")
	//inactiveMem := getMVal("Inactive:")

	totalSwap := getMVal("SwapTotal:")
	freeSwap := getMVal("SwapFree:")
	usedSwap := totalSwap - freeSwap

	fmtMemstr := fmt.Sprintf(
		"Memory Total %s | Free %s | Available %s | Cache %s | "+
			"Shared %s | Used %s (Active %6.2f%%)",
		humanizeMem(totalMem),
		humanizeMem(freeMem),
		humanizeMem(availMem),
		humanizeMem(cacheMem),
		humanizeMem(sharedMem),
		humanizeMem(usedMem),
		(100.0 * float32(activeMem) / float32(usedMem)),
	)

	fmtSwpstr := fmt.Sprintf(
		"Swap   Total %s | Used %s | Free %s",
		humanizeMem(totalSwap),
		humanizeMem(usedSwap),
		humanizeMem(freeSwap),
	)

	pdata := ui.PrintData{
		Position: ui.Tuple{X: 0, Y: 0},
		Size:     ui.Tuple{X: 2, Y: 100},
		Content:  []string{_HEADER_MEMINFO, fmtMemstr, fmtSwpstr},
	}

	c <- pdata
}

func humanizeMem(memKb int64) string {
	const KILO int64 = 1000
	floatDiv := func(a, b int64) string {
		f := float32(a) / float32(b)
		return fmt.Sprintf("%6.2f", f)
	}

	if memKb < KILO {
		return floatDiv(memKb, 1) + " kB"
	} else if memKb < KILO*KILO {
		return floatDiv(memKb, KILO) + " MB"
	} else if memKb < KILO*KILO*KILO {
		return floatDiv(memKb, KILO*KILO) + " GB"
	} else if memKb < KILO*KILO*KILO*KILO {
		return floatDiv(memKb, KILO*KILO*KILO) + " TB"
	} else {
		return floatDiv(memKb, KILO^4) + " PB"
	}
}
