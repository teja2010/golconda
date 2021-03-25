package config

import (
	"sync"
	"time"
	. "github.com/teja2010/golconda/src/any"
)

// A flat map. i.e. does not have maps within maps
// All Get calls must succeed. i.e. a default value must always be present.

type config_data struct {
	data TreeMap
	mt sync.RWMutex
}

var latest_config config_data

func GetStr(key string) string {
	latest_config.mt.RLock()
	defer latest_config.mt.RUnlock()

	return latest_config.data.Get(key).Str()
}

func GetInt(key string) int {
	latest_config.mt.RLock()
	defer latest_config.mt.RUnlock()

	return latest_config.data.Get(key).Int()
}

func GetMap(key string) TreeMap {
	latest_config.mt.RLock()
	defer latest_config.mt.RUnlock()

	return latest_config.data.Get(key).Map()
}

func ConfigInit() {
	// first fill the latest_config with the default values.
	// use https://golang.org/pkg/embed/ to embed the default configs

	latest_config.data = ReadMap("embed://default.json")
	latest_config.data.Set("update_interval", StrValue("2s"))
	latest_config.mt = sync.RWMutex{}
	go updater_thread()
}

// The only place where it will be written
func updater_thread() {
	// wait for the config file to change, if it changes, update the Map
	for {
		time.Sleep(time.Second)
	}
}
