package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "embed"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenOn         string   `yaml:"ListenOn"`
	Port             uint16   `yaml:"Port"`
	RelativePath     string   `yaml:"RelativePath"`
	LocalAssetsPaths []string `yaml:"LocalAssetsPaths"`
	NodeHosts        []string `yaml:"NodeHosts"`
}

//go:embed config.yaml
var configData []byte

var (
	confInstance *Config
	confOnce     sync.Once
)

func initConfig() {
	confOnce.Do(func() {
		configPath := filepath.Join("etc", "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// 创建目录并写入默认配置文件
			if err := os.MkdirAll("etc", 0755); err != nil {
				log.Default().Panicln("Failed to create etc directory:", err)
			}
			if err := os.WriteFile(configPath, configData, 0644); err != nil {
				log.Default().Panicln("Failed to write config file:", err)
			}
			log.Default().Println("Please edit etc/config.yaml.")
			// os.Exit(0)
			select {}
		} else if err != nil {
			log.Default().Panicln("Failed to check config file:", err)
		} else {

			data, err := os.ReadFile(configPath)
			if err != nil {
				log.Default().Panicln("Failed to read config file:", err)
			}
			confInstance = &Config{}
			if err := yaml.Unmarshal(data, confInstance); err != nil {
				log.Default().Panicln("Failed to parse config file:", err)
			}
		}
	})
}

func GetConfig() *Config {
	initConfig()
	return confInstance
}
