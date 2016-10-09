package config

import (
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"os"
	"strings"
	"fmt"
	"path/filepath"
)

var PvLink = make(map[string]string)

func init() {
	//全部递归一次, 得到volume
	recursiveVolumes(Config["base_path"])
	//监视BasePath目录变更情况
	NewWatcher(initWatch(Config["base_path"]))
}

func NewWatcher(paths []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				fmt.Println("change:", ev.Name)
				go recursiveVolumes(Config["base_path"])
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()
	for _, path := range paths {
		fmt.Println("addWatch:", path)
		err = watcher.Watch(path)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func initWatch(path string) []string{
	dirs, _ := ioutil.ReadDir(path)
	var paths []string
	paths = append(paths, path)
	for _, fi := range dirs {
		if !fi.IsDir() {
			break
		} else {
			var temp = path+ "/" +fi.Name() + "/"+ Config["volume_path"]
			exists, _ :=PathExists(temp)
			if(exists){
				paths = append(paths, temp)
			}
		}
	}
	return paths
}

//递归volume目录
func recursiveVolumes(path string) {
	PvLink = make(map[string]string)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if ( f == nil ) {
			return err
		}
		if f.IsDir() {
			ok := strings.HasSuffix(path, Config["mount_path"])
			if(ok){
				dirs, _ := ioutil.ReadDir(path)
				for _, fi := range dirs {
					if !fi.IsDir() {
						break
					} else {
						PvLink[fi.Name()] = path + "/" + fi.Name()
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	printPv()
}

func printPv(){
	for pv,path := range PvLink{
		fmt.Println(pv, path)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}