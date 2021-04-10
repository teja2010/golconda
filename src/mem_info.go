package golconda

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	d "github.com/teja2010/golconda/src/debug"
	"github.com/teja2010/golconda/src/meta"
	"github.com/teja2010/golconda/src/ui"
)

// /proc/meminfo
const (
	_PROC_MEMINFO = "/proc/meminfo"
)

// MemInfoConfig to read mem info
type MemInfoConfig struct {
	UpdateInterval string
	UIPosition     ui.Tuple
	UISize         ui.Tuple
	FmtString      []string
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

		_memInfo(c, conf.MemInfo)
	}
}

func confMemUpdateInterval(conf *GolcondaConfig) string {
	updateInterval := conf.MemInfo.UpdateInterval
	if updateInterval == "" {
		updateInterval = conf.Global.UpdateInterval
	}

	return updateInterval
}

type memStat struct {
	TotalMemKB, FreeMemKB, AvailMemKB, UsedMemKB int64
	CachedMemKB, ActiveMemKB, SharedMemKB        int64
	TotalSwapKB, FreeSwapKB, UsedSwapKB          int64

	// Human readable alternatives
	TotalMemH, FreeMemH, AvailMemH, UsedMemH string
	CachedMemH, ActiveMemH, SharedMemH       string
	TotalSwapH, FreeSwapH, UsedSwapH         string
}

func _memInfo(c chan<- ui.PrintData, conf MemInfoConfig) {
	_contents, err := ioutil.ReadFile(_PROC_MEMINFO)
	if err != nil {
		d.Error("Unable to read", _PROC_MEMINFO)
		return
	}

	pdata := ui.PrintData{
		Position: conf.UIPosition,
		Size:     conf.UISize,
		Content:  []string{},
	}

	contents := string(_contents)
	lines := strings.Split(contents, _NEWLINE)

	mStat := parseMemStat(lines)

	for _, fmtStr := range conf.FmtString {
		contentStr := meta.Format(fmtStr, mStat)
		pdata.Content = append(pdata.Content, contentStr)
	}

	// Use registers for this
	//TODO Active (100.0 * float32(activeMem) / float32(usedMem)),

	c <- pdata
}

func humanizeMem(memKb int64) string {
	const KILO int64 = 1000
	floatDiv := func(a, b int64) string {
		f := float32(a) / float32(b)
		return fmt.Sprintf("%.2f", f)
	}

	if memKb < KILO {
		return floatDiv(memKb, 1) + " KB"
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

func parseMemStat(lines []string) memStat {
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

	mStat := memStat{}

	mStat.TotalMemKB = getMVal("MemTotal:")
	mStat.FreeMemKB = getMVal("MemFree:")
	mStat.AvailMemKB = getMVal("MemAvailable:")
	mStat.UsedMemKB = mStat.TotalMemKB - mStat.AvailMemKB
	mStat.CachedMemKB = getMVal("Cached:")
	mStat.SharedMemKB = getMVal("Shmem:")
	mStat.ActiveMemKB = getMVal("Active:")
	//inactiveMemKB  = getMVal("Inactive:")
	mStat.TotalSwapKB = getMVal("SwapTotal:")
	mStat.FreeSwapKB = getMVal("SwapFree:")
	mStat.UsedSwapKB = mStat.TotalSwapKB - mStat.FreeSwapKB

	// humanize

	mStat.TotalMemH = humanizeMem(mStat.TotalMemKB)
	mStat.FreeMemH = humanizeMem(mStat.FreeMemKB)
	mStat.AvailMemH = humanizeMem(mStat.AvailMemKB)
	mStat.UsedMemH = humanizeMem(mStat.UsedMemKB)
	mStat.CachedMemH = humanizeMem(mStat.CachedMemKB)
	mStat.SharedMemH = humanizeMem(mStat.SharedMemKB)
	mStat.ActiveMemH = humanizeMem(mStat.ActiveMemKB)
	//inacTiveMemH  = humanizeMem(inactivEMemKB)
	mStat.TotalSwapH = humanizeMem(mStat.TotalSwapKB)
	mStat.FreeSwapH = humanizeMem(mStat.FreeSwapKB)
	mStat.UsedSwapH = humanizeMem(mStat.FreeSwapKB)

	return mStat
}
