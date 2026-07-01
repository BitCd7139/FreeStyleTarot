package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AgentSettings struct {
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

type AgentsConfig struct {
	IntentClarifier       AgentSettings `yaml:"intent_clarifier"`
	SpreadAnalyst         AgentSettings `yaml:"spread_analyst"`
	Advisor               AgentSettings `yaml:"advisor"`
	Persona               AgentSettings `yaml:"persona"`
	BriefMaterialMaxRunes int           `yaml:"brief_material_max_runes"`
}

type AuthYamlSettings struct {
	ForceLogin bool `yaml:"force_login"` // true 时必须先登录；false 时直接进入占卜界面
	SkipVerify bool `yaml:"skip_verify"` // true 时架空邮件验证码
}

type AnnouncementSettings struct {
	Enabled bool   `yaml:"enabled"`
	Title   string `yaml:"title"`
	Content string `yaml:"content"`
}

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`

	Auth AuthYamlSettings `yaml:"auth"`

	Announcement AnnouncementSettings `yaml:"announcement"`

	AI struct {
		ApiKey string `yaml:"api_key"`
		Model  string `yaml:"model"`
	} `yaml:"ai"`

	Agents AgentsConfig `yaml:"agents"`
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

	loadAuthConfig()
}
