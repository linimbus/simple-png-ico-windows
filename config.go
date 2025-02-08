package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Config struct {
	InputDir  string `json:"inputDir"`
	OutputDir string `json:"outputDir"`
	Mode      string `json:"mode"`
	Pixel     int    `json:"pixel"`
	TimeStamp bool   `json:"timestamp"`
}

var configCache = Config{
	InputDir:  "",
	OutputDir: "",
	Mode:      "",
	Pixel:     -1,
}

var configFilePath string
var configLock sync.Mutex

func configSyncToFile() error {
	configLock.Lock()
	defer configLock.Unlock()

	value, err := json.MarshalIndent(configCache, "\t", " ")
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
		return err
	}
	return os.WriteFile(configFilePath, value, 0664)
}

func ConfigGet() *Config {
	return &configCache
}

func InputDirSave(path string) error {
	configCache.InputDir = path
	return configSyncToFile()
}

func OutputDirSave(path string) error {
	configCache.OutputDir = path
	return configSyncToFile()
}

func TimeStampSave(flag bool) error {
	configCache.TimeStamp = flag
	return configSyncToFile()
}

func ModeSave(mode string) error {
	configCache.Mode = mode
	return configSyncToFile()
}

func PixelSave(pixel int) error {
	configCache.Pixel = pixel
	return configSyncToFile()
}

func ConfigInit() {
	configFilePath = fmt.Sprintf("%s%c%s", ConfigDirGet(), os.PathSeparator, "config.json")

	_, err := os.Stat(configFilePath)
	if err != nil {
		err = configSyncToFile()
		if err != nil {
			logs.Error("config sync to file fail, %s", err.Error())
			return
		}
	}

	value, err := os.ReadFile(configFilePath)
	if err != nil {
		logs.Error("read config file from app data dir fail, %s", err.Error())
		configSyncToFile()
		return
	}

	err = json.Unmarshal(value, &configCache)
	if err != nil {
		logs.Error("json unmarshal config fail, %s", err.Error())
		configSyncToFile()
		return
	}
}
