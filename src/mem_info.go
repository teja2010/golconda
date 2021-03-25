package golconda

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"io/ioutil"
	ui "github.com/teja2010/golconda/src/ui"
	d "github.com/teja2010/golconda/src/debug"
	conf "github.com/teja2010/golconda/src/config"
)

// /proc/meminfo
const (
	_PROC_MEMINFO = "/proc/meminfo"

	_HEADER_MEMINFO = "Memory Info:"
)

func Meminfo(c chan<- ui.PrintData) {

	for {
		update_interval := conf.GetStr("update_interval")
		duration, err := time.ParseDuration(update_interval)
		if err != nil {
			d.Bug("Invalid Duration:", update_interval)
		}
		time.Sleep(duration)

		__meminfo(c)
	}
}

func __meminfo(c chan<- ui.PrintData) {
	_contents, err := ioutil.ReadFile(_PROC_MEMINFO)
	if err != nil {
		d.Error("Unable to read", _PROC_MEMINFO)
		return
	}

	contents := string(_contents)
	lines := strings.Split(contents, _NEWLINE)

	getMVal := func(prefix string) int64 {
		s3:= FindLine(lines, Regex2Func("^" + prefix))
		s2 := strings.TrimPrefix(s3, prefix)
		s1 := strings.TrimSuffix(s2, " kB")
		s := strings.TrimSpace(s1)

		i64, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			d.Bug("ParseInt failed")
		}
		return i64
	}

	total_mem  := getMVal("MemTotal:")
	free_mem   := getMVal("MemFree:")
	avail_mem  := getMVal("MemAvailable:")
	used_mem   := total_mem - avail_mem
	cache_mem  := getMVal("Cached:")
	shared_mem := getMVal("Shmem:")

	active_mem   := getMVal("Active:")
	//inactive_mem := getMVal("Inactive:")

	total_swap := getMVal("SwapTotal:")
	free_swap  := getMVal("SwapFree:")
	used_swap  := total_swap - free_swap


	fmt_memstr := fmt.Sprintf(
		"Memory Total %s | Free %s | Available %s | Cache %s | " +
		"Shared %s | Used %s (Active %6.2f%%)" ,
		humanizeMem(total_mem),
		humanizeMem(free_mem),
		humanizeMem(avail_mem),
		humanizeMem(cache_mem),
		humanizeMem(shared_mem),
		humanizeMem(used_mem),
		(100.0*float32(active_mem)/float32(used_mem)),
	)

	fmt_swpstr := fmt.Sprintf(
		"Swap   Total %s | Used %s | Free %s",
		humanizeMem(total_swap),
		humanizeMem(used_swap),
		humanizeMem(free_swap),
	)

	pdata := ui.PrintData{
			ui.Tuple{0, 0},
			ui.Tuple{1, 100},
			[]string{_HEADER_MEMINFO, fmt_memstr, fmt_swpstr},
		}
	
	c <- pdata
}


func humanizeMem(mem_kb int64) string {
	const KILO int64 = 1000
	float_div := func(a, b int64) string {
		f := float32(a)/float32(b)
		return fmt.Sprintf("%6.2f", f)
	}

	if mem_kb < KILO {
		return float_div(mem_kb, 1) + " kB"
	} else if  mem_kb < KILO * KILO {
		return float_div(mem_kb, KILO) + " MB"
	} else if mem_kb < KILO * KILO * KILO {
		return float_div(mem_kb, KILO * KILO) + " GB"
	} else if mem_kb < KILO * KILO * KILO * KILO {
		return float_div(mem_kb, KILO * KILO * KILO) + " TB"
	} else {
		return float_div(mem_kb, KILO^4) + " PB"
	}
}

