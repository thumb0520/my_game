package llm

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 游戏配置
type Config struct {
	LLM  LLMConfig  `yaml:"llm"`
	Game GameConfig `yaml:"game"`
}

// LLMConfig LLM 配置
type LLMConfig struct {
	BaseURL     string  `yaml:"base_url"`
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	Enabled     bool    `yaml:"enabled"`
}

// GameConfig 游戏配置
type GameConfig struct {
	Difficulty string `yaml:"difficulty"`
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg := &Config{
		LLM: LLMConfig{
			BaseURL:     "https://api.openai.com/v1",
			Model:       "gpt-4o-mini",
			Temperature: 0.8,
			MaxTokens:   512,
		},
		Game: GameConfig{
			Difficulty: "normal",
		},
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 支持环境变量覆盖
	if envKey := os.Getenv("DUNGEONLOG_API_KEY"); envKey != "" {
		cfg.LLM.APIKey = envKey
	}
	if envURL := os.Getenv("DUNGEONLOG_BASE_URL"); envURL != "" {
		cfg.LLM.BaseURL = envURL
	}

	return cfg, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		LLM: LLMConfig{
			BaseURL:     "https://api.openai.com/v1",
			Model:       "gpt-4o-mini",
			Temperature: 0.8,
			MaxTokens:   512,
			Enabled:     false,
		},
		Game: GameConfig{
			Difficulty: "normal",
		},
	}
}
