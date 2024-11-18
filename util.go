package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/astaxie/beego/logs"
)

func VersionGet() string {
	return "v0.1.0"
}

func SaveToFile(name string, body []byte) error {
	return os.WriteFile(name, body, 0664)
}

func CapSignal(proc func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		proc()
		logs.Error("recv signcal %s, ready to exit", sig.String())
		os.Exit(-1)
	}()
}

func StringList(list []string) string {
	var body string
	for idx, v := range list {
		if idx == len(list)-1 {
			body += fmt.Sprintf("%s", v)
		} else {
			body += fmt.Sprintf("%s;", v)
		}
	}
	return body
}

func StringDiff(oldlist []string, newlist []string) ([]string, []string) {
	del := make([]string, 0)
	add := make([]string, 0)
	for _, v1 := range oldlist {
		flag := false
		for _, v2 := range newlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if !flag {
			del = append(del, v1)
		}
	}
	for _, v1 := range newlist {
		flag := false
		for _, v2 := range oldlist {
			if v1 == v2 {
				flag = true
				break
			}
		}
		if !flag {
			add = append(add, v1)
		}
	}
	return del, add
}

func StringClone(list []string) []string {
	output := make([]string, len(list))
	copy(output, list)
	return output
}

func ReadFileList(dir string) ([]string, error) {
	output := make([]string, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".png" {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			output = append(output, absPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
