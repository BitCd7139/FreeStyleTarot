package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`

	AI struct {
		ApiKey string `yaml:"api_key"`
		Model  string `yaml:"model"`
	} `yaml:"ai"`
}

// GlobalConfig 全局配置变量，供其他包调用
var GlobalConfig *Config

// InitConfig 初始化配置
func InitConfig() {
	// 读取配置文件
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析 YAML
	GlobalConfig = &Config{}
	err = yaml.Unmarshal(yamlFile, GlobalConfig)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}
