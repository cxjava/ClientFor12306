package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	CDN      []string `toml:"cdn"`
	LogLevel int      `toml:"log_level"`
	Proxy    bool
	ProxyUrl string `toml:"proxy_url"`
}

var Conf Config

func init() {
	ReadConfig()
}

//读取配置文件
func ReadConfig() error {
	if _, err := toml.DecodeFile("config.ini", &Conf); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
