package config

import (
	configtool "github.com/larspensjo/config"
	"flag"
	"fmt"
)

var (
	configFile = flag.String("configfile", "config.ini", "General configuration file")
)
//topic list
var Config = make(map[string]string)

func init()  {
	flag.Parse()
	//set config file std
	cfg, err := configtool.ReadDefault(*configFile)
	if err != nil {
		fmt.Println("Fail to find", *configFile, err)
	}
	//Initialized topic from the configuration
	for _, section := range cfg.Sections() {
		options, err := cfg.SectionOptions(section)
		if err == nil {
			for _, v := range options {
				options, err := cfg.String(section, v)
				if err == nil {
					Config[v] = options
					fmt.Println("load config -> "+v+":", options)
				}
			}
		}
	}
}
