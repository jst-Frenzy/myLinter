package analyzer

import (
	"encoding/json"
	"os"
)

type Config struct {
	SensitiveWords []string `json:"sensitive_words"`

	Checks struct {
		SpecialChars  bool `json:"special_chars"`
		CapitalLetter bool `json:"capital_letter"`
		EnglishOnly   bool `json:"english_only"`
		SensitiveData bool `json:"sensitive_data"`
	} `json:"checks"`
}

func DefaultConfig() *Config {
	cfg := &Config{
		SensitiveWords: []string{
			"password:",
			"api_key=",
			"token:",
		},
	}
	cfg.Checks.CapitalLetter = true
	cfg.Checks.SpecialChars = true
	cfg.Checks.EnglishOnly = true
	cfg.Checks.SensitiveData = true

	return cfg
}

func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
